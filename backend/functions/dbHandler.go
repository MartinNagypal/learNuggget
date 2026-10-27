package functions

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type CollectionElement struct {
	Name string
	Id   int
}

type Card struct {
	Question     string `json:"question"`
	Answer       string `json:"answer"`
	CollectionID int    `json:"collectionID"`
}

func openDb() (*sql.DB, error) {
	db, dbErr := sql.Open("sqlite", "learNuggget.db")
	if dbErr != nil {
		log.Fatal(dbErr)
		return nil, dbErr
	} else {
		return db, nil
	}
}

func InitTables() error {
	db, openErr := openDb()
	if db != nil && openErr == nil {
		defer db.Close()
		_, errCollection := db.Exec(`
		CREATE TABLE IF NOT EXISTS collection (
		    Id INTEGER PRIMARY KEY,
		    name TEXT NOT NULL UNIQUE,
		    right INTEGER DEFAULT 0,
		    wrong INTEGER DEFAULT 0,
		    createdAt DEFAULT CURRENT_TIMESTAMP
			)
		`)

		if errCollection != nil {
			log.Println("Error creating collection table in dbHandler.go:")
			log.Fatal(errCollection)
			return errCollection
		}

		_, errCard := db.Exec(`
		CREATE TABLE IF NOT EXISTS Card (
		    Id INTEGER PRIMARY KEY,
		    question TEXT NOT NULL UNIQUE,
			answer TEXT NOT NULL UNIQUE,
			marked INTEGER DEFAULT 0,
			collectionID INTEGER NOT NULL
			)
		`)

		if errCard != nil {
			log.Println("Error creating Card table in dbHandler.go:")
			log.Fatal(errCard)
			return errCard
		}
		return nil

	} else {
		return openErr
	}
}

func AddNewCollection(collectionName string) error {
	db, openErr := openDb()
	if db != nil && openErr == nil {
		defer db.Close()
		_, execErr := db.Exec(`INSERT INTO collection (name) VALUES (?)`, collectionName)
		if execErr != nil {
			log.Println("Error executing db command in dbHandler.go:")
			log.Println(execErr)
			return execErr
		}
	} else {
		log.Fatal(openErr)
		return openErr
	}
	return openErr
}

func RemoveCollection(collectionName string) error {
	db, openErr := openDb()
	defer db.Close()
	if openErr != nil {
		log.Println("Error in dbHandler.go")
		log.Println(openErr)
		return openErr
	}
	_, cardErr := db.Exec("DELETE FROM Card WHERE collectionID = (SELECT Id FROM collection WHERE name = ?)", collectionName)
	if cardErr != nil {
		log.Println("Error in dbHandler.go:")
		log.Println(cardErr)
		return cardErr
	}

	_, collErr := db.Exec("DELETE FROM collection WHERE name = ?", collectionName)
	if collErr != nil {
		log.Println("Error in dbHandler.go:")
		log.Println(collErr)
		return collErr
	}

	return nil
}

func FetchCollections(orderState int) ([]CollectionElement, error) {
	//order: 0=createdAt DESC 1=createdAt, 2=name 3= name DESC
	db, openErr := openDb()
	if openErr != nil {
		log.Println("Error in dbHandler.go")
		log.Println(openErr)
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

	query := fmt.Sprintf("SELECT Id, name FROM collection ORDER BY %s COLLATE NOCASE %s", orderBy, order)

	rows, queryErr := db.Query(query) //Id auch fetchen, error handling verbessern
	if queryErr != nil {
		log.Fatal(queryErr)
		return nil, queryErr

	}
	var collections []CollectionElement
	for rows.Next() {
		var collection CollectionElement
		if err := rows.Scan(&collection.Id, &collection.Name); err != nil {
			return nil, err
		}
		collections = append(collections, collection)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return collections, nil
}

func AddNewCard(question, answer string, collectionID int) error {
	db, dbErr := openDb()
	if dbErr != nil || db == nil {
		log.Println("Error in dbHandler.go:")
		log.Println(dbErr)
		return dbErr
	}

	_, execErr := db.Exec("INSERT INTO Card (question, answer, collectionID) VALUES (?,?,?)", question, answer, collectionID)
	if execErr != nil {
		log.Println("Error in dbHandler.go")
		log.Println(execErr)
		return execErr
	} else {
		return nil
	}
}

func FetchCards(collectionID int) ([]Card, error) {
	db, dbErr := openDb()
	defer db.Close()
	if dbErr != nil || db == nil {
		log.Println("Error in dbHandler.go:")
		log.Println(dbErr)
		return nil, dbErr
	}

	rows, queryErr := db.Query("SELECT question, answer, collectionID FROM Card WHERE collectionID = ?", collectionID)
	defer rows.Close()
	if queryErr != nil {
		log.Println("Error in dbHandler.go:")
		log.Println(queryErr)
		return nil, queryErr
	}

	cards := []Card{}
	for rows.Next() {
		var cardElement Card
		scanErr := rows.Scan(&cardElement.Question, &cardElement.Answer, &cardElement.CollectionID)
		if scanErr != nil {
			log.Println("Error in dbHandler.go:")
			log.Println(scanErr)
			return nil, scanErr
		}
		cards = append(cards, cardElement)
	}
	return cards, nil

}
