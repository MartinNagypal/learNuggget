package functions

import "errors"

// ErrCardNotFound is returned when a card does not exist or does not belong to the requesting user.
var ErrCardNotFound = errors.New("Card not found")

func AddNewCard(question, answer string, collectionID int) error {
	db, dbErr := openDb()
	if dbErr != nil || db == nil {
		LogError(GetFile(), dbErr)
		return dbErr
	}
	defer db.Close()
	_, execErr := db.Exec("INSERT INTO Card (question, answer, collectionID) VALUES (?,?,?)", question, answer, collectionID)
	if execErr != nil {
		LogError(GetFile(), execErr)
		return execErr
	} else {
		return nil
	}
}

func FetchCards(collectionID int) ([]Card, error) {
	db, dbErr := openDb()
	if dbErr != nil || db == nil {
		LogError(GetFile(), dbErr)
		return nil, dbErr
	}
	defer db.Close()

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
	defer db.Close()

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
	defer db.Close()

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
		return errors.New("Card not found")
	}
	return nil
}

func FetchMarkedCards(userID int) ([]Card, error) {
	db, dbErr := openDb()
	if dbErr != nil {
		LogError(GetFile(), dbErr)
		return nil, dbErr
	}
	defer db.Close()

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
