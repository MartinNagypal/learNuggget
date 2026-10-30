package functions

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

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

func Login(username, password string) (string, int, error) {
	db, openErr := openDb()
	if openErr != nil {
		LogError(GetFile(), openErr)
		return "", http.StatusInternalServerError, openErr
	}
	var user UserLogin
	queryErr := db.QueryRow("SELECT userID, username, password FROM UserT WHERE username = (?)", username).Scan(&user.id, &user.Username, &user.Password)
	if queryErr != nil {
		if queryErr == sql.ErrNoRows {
			return "", http.StatusForbidden, errors.New("Wrong username or password")
		}
		LogError(GetFile(), queryErr)
		return "", http.StatusInternalServerError, queryErr
	}

	encryptErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if encryptErr != nil {
		LogError(GetFile(), encryptErr)
		return "", http.StatusForbidden, errors.New("Wrong username or password")
	}

	sessionToken, sessionErr := generateSession(user.id)
	if sessionErr != nil {
		return "", http.StatusInternalServerError, sessionErr
	}

	return sessionToken, http.StatusOK, nil
}

func VerifySession(sessionToken string) (int, error) {
	db, openErr := openDb()
	if openErr != nil {
		LogError(GetFile(), openErr)
		return 0, openErr
	}
	defer db.Close()

	tokenHash := GenerateHash(sessionToken)
	var expiresAt int
	var userId int
	queryErr := db.QueryRow("SELECT expiresAt, userID FROM SessionT WHERE sessionToken = (?)", tokenHash).Scan(&expiresAt, &userId)
	if queryErr != nil {
		if queryErr == sql.ErrNoRows {
			return 0, errors.New("Session not found")
		}
		LogError(GetFile(), queryErr)
		return 0, queryErr
	}

	if int64(expiresAt) <= time.Now().Unix() {
		_, deleteErr := db.Exec("DELETE FROM SessionT WHERE sessionToken = ?", tokenHash)
		if deleteErr != nil {
			LogError(GetFile(), deleteErr)
			return 0, deleteErr
		}
		return 0, errors.New("Session Expired")
	}
	return userId, nil
}

func FuncLogout(userID int) error {
	db, dbErr := openDb()
	if dbErr != nil {
		LogError(GetFile(), dbErr)
		return dbErr
	}

	_, execErr := db.Exec("DELETE FROM SessionT WHERE userID = ?", userID)
	if execErr != nil {
		LogError(GetFile(), execErr)
		return execErr
	}
	return nil
}
