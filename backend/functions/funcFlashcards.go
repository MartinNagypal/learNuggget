package functions

import (
	"errors"
	"strings"
)

// ErrCardNotFound is returned when a card does not exist or does not belong to the requesting user.
var ErrCardNotFound = errors.New("Card not found")

// CardInput is a question/answer pair to be inserted.
type CardInput struct {
	Question string
	Answer   string
}

// Placeholders for escaped separators while splitting. \x00 can't be typed
// into the textarea, so it never collides with real card text.
const (
	escapedCardSep = "\x00CARDSEP\x00"
	escapedQASep   = "\x00QASEP\x00"
)

// ParseBulkCards splits the bulk text "question@@answer;;question@@answer;;..."
// into cards. To use the separators inside a card, escape them with a
// backslash: \;; becomes ";;" and \@@ becomes "@@". Question and answer are
// trimmed; parts without "@@" are skipped. The caller validates each card.
func ParseBulkCards(bulk string) []CardInput {
	text := strings.ReplaceAll(bulk, `\;;`, escapedCardSep)
	text = strings.ReplaceAll(text, `\@@`, escapedQASep)
	restore := strings.NewReplacer(escapedCardSep, ";;", escapedQASep, "@@")

	cards := []CardInput{}
	for _, element := range strings.Split(text, ";;") {
		sepElement := strings.SplitN(element, "@@", 2)
		if len(sepElement) != 2 {
			continue
		}
		cards = append(cards, CardInput{
			Question: strings.TrimSpace(restore.Replace(sepElement[0])),
			Answer:   strings.TrimSpace(restore.Replace(sepElement[1])),
		})
	}
	return cards
}

func AddNewCard(question, answer string, collectionID int) error {
	return AddCards(collectionID, []CardInput{{Question: question, Answer: answer}})
}

// AddCards inserts the cards using one connection and one transaction: either
// all cards are added or none (no half-finished imports). The caller checks
// each card with ValidateCard first. Returns ErrTooManyCards if the collection
// would exceed MaxCardsPerCollection.
func AddCards(collectionID int, cards []CardInput) error {
	db, dbErr := openDb()
	if dbErr != nil {
		LogError(GetFile(), dbErr)
		return dbErr
	}

	tx, txErr := db.Begin()
	if txErr != nil {
		LogError(GetFile(), txErr)
		return txErr
	}
	defer tx.Rollback()

	stmt, prepErr := tx.Prepare("INSERT INTO Card (question, answer, collectionID) VALUES (?,?,?)")
	if prepErr != nil {
		LogError(GetFile(), prepErr)
		return prepErr
	}
	defer stmt.Close()

	for _, card := range cards {
		if _, execErr := stmt.Exec(card.Question, card.Answer, collectionID); execErr != nil {
			LogError(GetFile(), execErr)
			return execErr
		}
	}

	// Counted after the inserts (inside the transaction) so the write lock is
	// already held and concurrent imports cannot slip past the limit together.
	var total int
	if countErr := tx.QueryRow("SELECT COUNT(*) FROM Card WHERE collectionID = ?", collectionID).Scan(&total); countErr != nil {
		LogError(GetFile(), countErr)
		return countErr
	}
	if total > MaxCardsPerCollection {
		return ErrTooManyCards
	}

	if commitErr := tx.Commit(); commitErr != nil {
		LogError(GetFile(), commitErr)
		return commitErr
	}
	return nil
}

func FetchCards(collectionID int) ([]Card, error) {
	db, dbErr := openDb()
	if dbErr != nil || db == nil {
		LogError(GetFile(), dbErr)
		return nil, dbErr
	}

	rows, queryErr := db.Query("SELECT question, answer, cardID, collectionID, marked FROM Card WHERE collectionID = ?", collectionID)
	if queryErr != nil {
		LogError(GetFile(), queryErr)
		return nil, queryErr
	}
	defer rows.Close()

	cards := []Card{}
	for rows.Next() {
		var cardElement Card
		scanErr := rows.Scan(&cardElement.Question, &cardElement.Answer, &cardElement.CardID, &cardElement.CollectionID, &cardElement.Marked)
		if scanErr != nil {
			LogError(GetFile(), scanErr)
			return nil, scanErr
		}
		cards = append(cards, cardElement)
	}
	return cards, nil
}

func DeleteCard(cardID, userID int) error {
	db, dbErr := openDb()
	if dbErr != nil {
		LogError(GetFile(), dbErr)
		return dbErr
	}

	result, execErr := db.Exec(`
		DELETE FROM Card
		WHERE cardID = ?
		  AND collectionID IN (SELECT collectionID FROM Collection WHERE userID = ?)`,
		cardID, userID)
	if execErr != nil {
		LogError(GetFile(), execErr)
		return execErr
	}

	affected, affectedErr := result.RowsAffected()
	if affectedErr != nil {
		LogError(GetFile(), affectedErr)
		return affectedErr
	}
	if affected == 0 {
		return ErrCardNotFound
	}
	return nil
}

func MarkCard(cardID int, marked bool, userID int) error {
	db, dbErr := openDb()
	if dbErr != nil {
		LogError(GetFile(), dbErr)
		return dbErr
	}

	result, execErr := db.Exec(`
		UPDATE Card SET marked = ?
		WHERE cardID = ?
		  AND collectionID IN (SELECT collectionID FROM Collection WHERE userID = ?)`,
		marked, cardID, userID)
	if execErr != nil {
		LogError(GetFile(), execErr)
		return execErr
	}

	affected, affectedErr := result.RowsAffected()
	if affectedErr != nil {
		LogError(GetFile(), affectedErr)
		return affectedErr
	}
	if affected == 0 {
		return ErrCardNotFound
	}
	return nil
}

func FetchMarkedCards(userID int) ([]Card, error) {
	db, dbErr := openDb()
	if dbErr != nil {
		LogError(GetFile(), dbErr)
		return nil, dbErr
	}

	rows, queryErr := db.Query(`
		SELECT c.question, c.answer, c.cardID, c.collectionID, c.marked
		FROM Card c
		JOIN Collection col ON c.collectionID = col.collectionID
		WHERE col.userID = ? AND c.marked = 1
		ORDER BY col.name, c.cardID`, userID)
	if queryErr != nil {
		LogError(GetFile(), queryErr)
		return nil, queryErr
	}
	defer rows.Close()

	cards := []Card{}
	for rows.Next() {
		var cardElement Card
		scanErr := rows.Scan(&cardElement.Question, &cardElement.Answer, &cardElement.CardID, &cardElement.CollectionID, &cardElement.Marked)
		if scanErr != nil {
			LogError(GetFile(), scanErr)
			return nil, scanErr
		}
		cards = append(cards, cardElement)
	}
	return cards, nil
}
