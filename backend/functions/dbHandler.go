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
	if openErr != nil {
		LogError(GetFile(), openErr)
		return openErr
	}
	defer db.Close()
	_, errCollection := db.Exec(`
		CREATE TABLE IF NOT EXISTS Collection (
		    collectionID INTEGER PRIMARY KEY,
		    name TEXT NOT NULL UNIQUE,
		    userID INTEGER NOT NULL,
		    createdAt DEFAULT CURRENT_TIMESTAMP
			)
		`)

	if errCollection != nil {
		LogError(GetFile(), errCollection)
		return errCollection
	}

	_, errCard := db.Exec(`
		CREATE TABLE IF NOT EXISTS Card (
		    cardID INTEGER PRIMARY KEY,
		    question TEXT NOT NULL,
			answer TEXT NOT NULL,
			collectionID INTEGER NOT NULL
			)
		`)

	if errCard != nil {
		LogError(GetFile(), errCard)
		return errCard
	}

	_, errUser := db.Exec(`
		CREATE TABLE IF NOT EXISTS UserT (
			userID INTEGER PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL
			)
		`)

	if errUser != nil {
		LogError(GetFile(), errUser)
		return errUser
	}

	_, errSession := db.Exec(`
		CREATE TABLE IF NOT EXISTS SessionT (
			sessionID INTEGER PRIMARY KEY,
			sessionToken TEXT UNIQUE NOT NULL,
			expiresAt INTEGER NOT NULL,
			createdAt INTEGER NOT NULL DEFAULT (unixepoch()),
		    userID INTEGER NOT NULL
			)
		`)

	if errSession != nil {
		LogError(GetFile(), errSession)
		return errSession
	}
	return nil
}
