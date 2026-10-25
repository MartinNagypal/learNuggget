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
		    id INTEGER PRIMARY KEY,
		    name TEXT NOT NULL UNIQUE,
		    right INTEGER DEFAULT 0,
		    wrong INTEGER DEFAULT 0
			)
		`)

		if errCollection != nil {
			log.Println("Error creating collection table in dbHandler.go:")
			log.Fatal(errCollection)
			return errCollection
		}

		_, errCard := db.Exec(`
		CREATE TABLE IF NOT EXISTS card (
		    id INTEGER PRIMARY KEY,
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
	if db != nil && openErr == nil {
		defer db.Close()
		_, execErr := db.Exec(`DELETE FROM collection WHERE name = ?`, collectionName)
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

func FetchCollections() ([]string, error) {
	db, openErr := openDb()
	if db != nil && openErr == nil {
		rows, queryErr := db.Query(`SELECT name FROM collection ORDER BY id`)
		if queryErr != nil {
			log.Fatal(queryErr)
			return nil, queryErr
		} else {
			var collections []string
			for rows.Next() {
				var name string
				if err := rows.Scan(&name); err != nil {
					return nil, err
				}
				collections = append(collections, name)
			}
			if err := rows.Err(); err != nil {
				return nil, err
			}
			return collections, nil
		}
	}
	return nil, openErr
}
