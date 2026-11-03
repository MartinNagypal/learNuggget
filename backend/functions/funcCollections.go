package functions

import "errors"

func AddNewCollection(collectionName string, userID, groupID int) error {
	db, openErr := openDb()
	if openErr != nil {
		LogError(GetFile(), openErr)
		return openErr
	}
	defer db.Close()

	result, execErr := db.Exec(`
		INSERT INTO Collection (name, userID, groupID)
		SELECT ?, ?, ?
		WHERE ? = 0
		   OR EXISTS (SELECT 1 FROM CollectionGroup WHERE groupID = ? AND userID = ?)`,
		collectionName, userID, groupID,
		groupID,
		groupID, userID,
	)
	if execErr != nil {
		LogError(GetFile(), execErr)
		return execErr
	}

	rows, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		LogError(GetFile(), rowsErr)
		return rowsErr
	}
	if rows == 0 {
		return errors.New("´Group not found")
	}
	return nil
}

func RemoveCollection(collectionID, userID int) error {
	db, openErr := openDb()
	if openErr != nil {
		LogError(GetFile(), openErr)
		return openErr
	}
	defer db.Close()
	_, cardErr := db.Exec("DELETE FROM Card WHERE collectionID = (?)", collectionID)
	if cardErr != nil {
		LogError(GetFile(), cardErr)
		return cardErr
	}

	_, collErr := db.Exec("DELETE FROM Collection WHERE collectionID = (?) AND userID = (?)", collectionID, userID)
	if collErr != nil {
		LogError(GetFile(), collErr)
		return collErr
	}
	return nil
}

func FetchCollections(userID, groupID int) ([]CollectionElement, error) {
	db, openErr := openDb()
	if openErr != nil {
		LogError(GetFile(), openErr)
		return nil, openErr
	}
	defer db.Close()

	rows, queryErr := db.Query("SELECT collectionID, name, groupID FROM Collection WHERE userID = (?) AND groupID = (?)", userID, groupID)
	if queryErr != nil {
		LogError(GetFile(), queryErr)
		return nil, queryErr

	}
	defer rows.Close()
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
