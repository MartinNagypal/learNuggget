package functions

func AddNewCollection(collectionName string, userID, groupID int) error {
	db, openErr := openDb()
	if openErr != nil {
		LogError(GetFile(), openErr)
	}
	defer db.Close()

	_, execErr := db.Exec(`INSERT INTO Collection (name, userID, groupID) VALUES (?, ?, ?)`, collectionName, userID, groupID)
	if execErr != nil {
		LogError(GetFile(), execErr)
		return execErr
	}
	return openErr
}

func RemoveCollection(collectionName string, userID int) error {
	db, openErr := openDb()
	if openErr != nil {
		LogError(GetFile(), openErr)
		return openErr
	}
	defer db.Close()
	_, cardErr := db.Exec("DELETE FROM Card WHERE collectionID = (SELECT collectionID FROM Collection WHERE name = ? AND userID = ?)", collectionName, userID)
	if cardErr != nil {
		LogError(GetFile(), cardErr)
		return cardErr
	}

	_, collErr := db.Exec("DELETE FROM Collection WHERE name = ? AND userID = ?", collectionName, userID)
	if collErr != nil {
		LogError(GetFile(), collErr)
		return collErr
	}
	return nil
}

func FetchCollections(userID, groupID int) ([]CollectionElement, error) {
	//order: 0=createdAt DESC 1=createdAt, 2=name 3= name DESC
	db, openErr := openDb()
	if openErr != nil {
		LogError(GetFile(), openErr)
		return nil, openErr
	}
	defer db.Close()

	rows, queryErr := db.Query("SELECT collectionID, name, groupID FROM Collection WHERE userID = (?)", userID)
	if queryErr != nil {
		LogError(GetFile(), queryErr)
		return nil, queryErr

	}
	var collections []CollectionElement
	for rows.Next() {
		var collection CollectionElement
		if err := rows.Scan(&collection.Id, &collection.Name, &collection.Gid); err != nil {
			LogError(GetFile(), err)
			return nil, err
		}
		collections = append(collections, collection)
	}
	if err := rows.Err(); err != nil {
		LogError(GetFile(), err)
		return nil, err
	}
	return collections, nil
}

func GetCollectionLen(collectionID int) (int, error) {
	db, dbErr := openDb()

	if dbErr != nil {
		LogError(GetFile(), dbErr)
		return 0, dbErr
	}
	defer db.Close()

	var len int
	queryErr := db.QueryRow("SELECT COUNT(*) FROM Card WHERE collectionID = ?", collectionID).Scan(&len)

	if queryErr != nil {
		LogError(GetFile(), queryErr)
		return 0, queryErr
	}

	return len, nil
}
