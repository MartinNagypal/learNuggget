package functions

import (
	"golang.org/x/crypto/bcrypt"
)

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

	sessionToken, sessionTokenErr := GenerateSessionToken()
	if sessionTokenErr != nil {
		LogError(GetFile(), sessionTokenErr)
		return "", sessionTokenErr
	}

	tokenHash := GenerateHash(sessionToken)

	_, saveSessionErr := db.Exec(
		`INSERT INTO SessionT (userID, sessionToken, expiresAt, createdAt) VALUES (?, ?, unixepoch() + 43200, unixepoch())`,
		userID,
		tokenHash,
	)

	if saveSessionErr != nil {
		LogError(GetFile(), saveSessionErr)
		return "", saveSessionErr
	}

	return sessionToken, nil
}
