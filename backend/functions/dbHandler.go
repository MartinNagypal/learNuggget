package functions

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

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
		    question TEXT NOT NULL,
			answer TEXT NOT NULL,
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
