package functions

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

// busy_timeout: when another request holds the write lock, wait up to 5s
// instead of failing immediately with SQLITE_BUSY ("database is locked").
const dbDSN = "learNuggget.db?_pragma=busy_timeout(5000)"

func openDb() (*sql.DB, error) {
	db, dbErr := sql.Open("sqlite", dbDSN)
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
		    name TEXT NOT NULL,
		    userID INTEGER NOT NULL,
		    GroupID INTEGER DEFAULT 0,
		    createdAt DEFAULT CURRENT_TIMESTAMP
			)
		`)

	if errCollection != nil {
		LogError(GetFile(), errCollection)
		return errCollection
	}

	_, errGroups := db.Exec(`
		CREATE TABLE IF NOT EXISTS CollectionGroup (
		    GroupID INTEGER PRIMARY KEY,
		    name TEXT NOT NULL,
		    userID INTEGER NOT NULL
			)
		`)

	if errGroups != nil {
		LogError(GetFile(), errGroups)
		return errGroups
	}

	_, errCard := db.Exec(`
		CREATE TABLE IF NOT EXISTS Card (
		    cardID INTEGER PRIMARY KEY,
		    question TEXT NOT NULL,
			answer TEXT NOT NULL,
			marked BOOLEAN NOT NULL DEFAULT 0,
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
			password TEXT NOT NULL,
			root INTEGER NOT NULL DEFAULT 0
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

	_, errLoginAttempt := db.Exec(`
		CREATE TABLE IF NOT EXISTS LoginAttemptT (
			username TEXT PRIMARY KEY,
			attempts INTEGER NOT NULL DEFAULT 0,
			lockedUntil INTEGER NOT NULL DEFAULT 0
			)
		`)

	if errLoginAttempt != nil {
		LogError(GetFile(), errLoginAttempt)
		return errLoginAttempt
	}

	_, errCollectionProgress := db.Exec(`
		CREATE TABLE IF NOT EXISTS CollectionProgress (
			collectionProgressID INTEGER PRIMARY KEY,
			avgScore REAL DEFAULT 0,
			lastScore REAL DEFAULT 0,
			totalPlayed INTEGER DEFAULT 0,
		    collectionID INTEGER NOT NULL
			)
		`)

	if errCollectionProgress != nil {
		LogError(GetFile(), errCollectionProgress)
		return errCollectionProgress
	}
	return nil
}
