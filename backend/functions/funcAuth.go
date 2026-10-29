package functions

import (
	"golang.org/x/crypto/bcrypt"
)

func generateSession(userId int64) (string, error) {
	sessionToken, sessionTokenErr := GenerateSessionToken()
	if sessionTokenErr != nil {
		LogError(GetFile(), sessionTokenErr)
		return "", sessionTokenErr
	}

	tokenHash := GenerateHash(sessionToken)

	db, openErr := openDb()
	if openErr != nil {
		LogError(GetFile(), openErr)
		return "", openErr
	}

	_, saveSessionErr := db.Exec(
		`INSERT INTO SessionT (userID, sessionToken, expiresAt, createdAt) VALUES (?, ?, unixepoch() + 43200, unixepoch())`,
		userId,
		tokenHash,
	)

	if saveSessionErr != nil {
		LogError(GetFile(), saveSessionErr)
		return "", saveSessionErr
	}
	return sessionToken, nil
}

func CreateUser(username, password string) (string, error) {
	db, openDbErr := openDb()
	if openDbErr != nil {
		LogError(GetFile(), openDbErr)
		return "", openDbErr
	}
	defer db.Close()
	encryptedPassword, encryptionErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if encryptionErr != nil {
		LogError(GetFile(), encryptionErr)
		return "", encryptionErr
	}

	result, execErr := db.Exec("INSERT INTO UserT (username, password) VALUES (?, ?) ", username, string(encryptedPassword))
	if execErr != nil {
		LogError(GetFile(), execErr)
		return "", execErr
	}
	userID, userIdErr := result.LastInsertId()
	if userIdErr != nil {
		LogError(GetFile(), userIdErr)
		return "", userIdErr
	}

	sessionToken, sessionErr := generateSession(userID)
	if sessionErr != nil {
		LogError(GetFile(), sessionErr)
		return "", sessionErr
	}

	return sessionToken, nil
}

func CheckIfUserExists(username string) (bool, error) {
	db, openErr := openDb()
	if openErr != nil {
		LogError(GetFile(), openErr)
		return false, openErr
	}
	defer db.Close()
	var exists bool
	queryErr := db.QueryRow("SELECT EXISTS(SELECT 1 FROM UserT WHERE username = ?)", username).Scan(&exists)
	if queryErr != nil {
		LogError(GetFile(), queryErr)
		return false, queryErr
	}
	return exists, nil
}
