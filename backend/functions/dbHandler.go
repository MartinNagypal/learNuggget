package functions

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

type collectionElement struct {
	Name string
	Id   int
}

type card struct {
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
		CREATE TABLE IF NOT EXISTS card (
		    Id INTEGER PRIMARY KEY,
		    question TEXT NOT NULL UNIQUE,
			answer TEXT NOT NULL UNIQUE,
			marked INTEGER DEFAULT 0,
			collectionID INTEGER NOT NULL
			)
		`)

		if errCard != nil {
			log.Println("Error creating card table in dbHandler.go:")
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
	_, cardErr := db.Exec("DELETE FROM card WHERE collectionID = (SELECT Id FROM collection WHERE name = ?)", collectionName)
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

func FetchCollections() ([]collectionElement, error) {
	db, openErr := openDb()
	if openErr != nil {
		log.Println("Error in dbHandler.go")
		log.Println(openErr)
		return nil, openErr
	}

	rows, queryErr := db.Query(`SELECT Id, name FROM collection ORDER BY createdAt DESC`) //Id auch fetchen, error handling verbessern
	if queryErr != nil {
		log.Fatal(queryErr)
		return nil, queryErr

	}
	var collections []collectionElement
	for rows.Next() {
		var collection collectionElement
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

	_, execErr := db.Exec("INSERT INTO card (question, answer, collectionID) VALUES (?,?,?)", question, answer, collectionID)
	if execErr != nil {
		log.Println("Error in dbHandler.go")
		log.Println(execErr)
		return execErr
	} else {
		return nil
	}
}

func FetchCards(collectionID int) ([]card, error) {
	db, dbErr := openDb()
	defer db.Close()
	if dbErr != nil || db == nil {
		log.Println("Error in dbHandler.go:")
		log.Println(dbErr)
		return nil, dbErr
	}

	rows, queryErr := db.Query("SELECT question, answer, collectionID FROM card WHERE collectionID = ?", collectionID)
	defer rows.Close()
	if queryErr != nil {
		log.Println("Error in dbHandler.go:")
		log.Println(queryErr)
		return nil, queryErr
	}

	cards := []card{}
	for rows.Next() {
		var cardElement card
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
