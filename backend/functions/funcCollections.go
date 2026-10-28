package functions

import (
	"fmt"
)

func AddNewCollection(collectionName string) error {
	db, openErr := openDb()
	if openErr != nil {
		LogError(GetFile(), openErr)
	}
	defer db.Close()

	_, execErr := db.Exec(`INSERT INTO Collection (name) VALUES (?)`, collectionName)
	if execErr != nil {
		LogError(GetFile(), execErr)
		return execErr
	}
	return openErr
}

func RemoveCollection(collectionName string) error {
	db, openErr := openDb()
	defer db.Close()
	if openErr != nil {
		LogError(GetFile(), openErr)
		return openErr
	}
	_, cardErr := db.Exec("DELETE FROM Card WHERE collectionID = (SELECT collectionID FROM Collection WHERE name = ?)", collectionName)
	if cardErr != nil {
		LogError(GetFile(), cardErr)
		return cardErr
	}

	_, collErr := db.Exec("DELETE FROM Collection WHERE name = ?", collectionName)
	if collErr != nil {
		LogError(GetFile(), collErr)
		return collErr
	}
	return nil
}

func FetchCollections(orderState int) ([]CollectionElement, error) {
	//order: 0=createdAt DESC 1=createdAt, 2=name 3= name DESC
	db, openErr := openDb()
	if openErr != nil {
		LogError(GetFile(), openErr)
		return nil, openErr
	}

	var orderBy, order string
	switch orderState {
	case 0:
		orderBy = "createdAt"
		order = "DESC"
	case 1:
		orderBy = "createdAt"
		order = "ASC"
	case 2:
		orderBy = "name"
		order = "ASC"
	case 3:
		orderBy = "name"
		order = "DESC"
	default:
		orderBy = "name"
		order = "DESC"
	}

	query := fmt.Sprintf("SELECT collectionID, name FROM Collection ORDER BY %s COLLATE NOCASE %s", orderBy, order)

	rows, queryErr := db.Query(query) //Id auch fetchen, error handling verbessern
	if queryErr != nil {
		LogError(GetFile(), queryErr)
		return nil, queryErr

	}
	var collections []CollectionElement
	for rows.Next() {
		var collection CollectionElement
		if err := rows.Scan(&collection.Id, &collection.Name); err != nil {
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

	var len int
	queryErr := db.QueryRow("SELECT COUNT(*) FROM Card WHERE collectionID = ?", collectionID).Scan(&len)

	if queryErr != nil {
		LogError(GetFile(), queryErr)
		return 0, queryErr
	}

	return len, nil
}
