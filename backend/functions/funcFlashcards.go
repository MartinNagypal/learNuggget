package functions

func AddNewCard(question, answer string, collectionID int) error {
	db, dbErr := openDb()
	if dbErr != nil || db == nil {
		LogError(GetFile(), dbErr)
		return dbErr
	}

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

	rows, queryErr := db.Query("SELECT question, answer, id, collectionID FROM Card WHERE collectionID = ?", collectionID)
	if queryErr != nil {
		LogError(GetFile(), queryErr)
		return nil, queryErr
	}
	defer rows.Close()

	cards := []Card{}
	for rows.Next() {
		var cardElement Card
		scanErr := rows.Scan(&cardElement.Question, &cardElement.Answer, &cardElement.CardID, &cardElement.CollectionID)
		if scanErr != nil {
			LogError(GetFile(), scanErr)
			return nil, scanErr
		}
		cards = append(cards, cardElement)
	}
	return cards, nil
}

func DeleteCard(cardId int) error {
	db, dbErr := openDb()
	defer db.Close()

	if dbErr != nil {
		LogError(GetFile(), dbErr)
	}

	_, execErr := db.Exec("DELETE FROM card WHERE id = ?", cardId)

	if execErr != nil {
		LogError(GetFile(), execErr)
		return execErr
	}

	return nil
}
