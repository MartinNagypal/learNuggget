package functions

import (
	"database/sql"
	"errors"
	"fmt"
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
	defer db.Close()

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

	var userCount int
	countErr := db.QueryRow("SELECT COUNT(*) FROM UserT").Scan(&userCount)
	if countErr != nil {
		LogError(GetFile(), countErr)
		return "", countErr
	}
	root := 0
	if userCount == 0 {
		root = 1
	}

	result, execErr := db.Exec("INSERT INTO UserT (username, password, root) VALUES (?, ?, ?) ", username, string(encryptedPassword), root)
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

const maxLoginAttempts = 5
const loginLockSeconds = 300

func Login(username, password string) (string, int, error) {
	lockedFor, lockErr := getLoginLock(username)
	if lockErr != nil {
		return "", http.StatusInternalServerError, lockErr
	}
	if lockedFor > 0 {
		return "", http.StatusTooManyRequests, lockedError(lockedFor)
	}

	db, openErr := openDb()
	if openErr != nil {
		LogError(GetFile(), openErr)
		return "", http.StatusInternalServerError, openErr
	}
	defer db.Close()
	var user UserLogin
	queryErr := db.QueryRow("SELECT userID, username, password FROM UserT WHERE username = (?)", username).Scan(&user.id, &user.Username, &user.Password)
	if queryErr != nil {
		if queryErr == sql.ErrNoRows {
			return failedLogin(username)
		}
		LogError(GetFile(), queryErr)
		return "", http.StatusInternalServerError, queryErr
	}

	encryptErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if encryptErr != nil {
		LogError(GetFile(), encryptErr)
		return failedLogin(username)
	}

	resetErr := resetLoginAttempts(username)
	if resetErr != nil {
		return "", http.StatusInternalServerError, resetErr
	}

	sessionToken, sessionErr := generateSession(user.id)
	if sessionErr != nil {
		return "", http.StatusInternalServerError, sessionErr
	}

	return sessionToken, http.StatusOK, nil
}

func failedLogin(username string) (string, int, error) {
	attemptsLeft, attemptErr := registerFailedLogin(username)
	if attemptErr != nil {
		return "", http.StatusInternalServerError, attemptErr
	}
	if attemptsLeft == 0 {
		return "", http.StatusTooManyRequests, lockedError(loginLockSeconds)
	}
	if attemptsLeft == 1 {
		return "", http.StatusForbidden, errors.New("Wrong username or password. 1 attempt left")
	}
	return "", http.StatusForbidden, fmt.Errorf("Wrong username or password. %d attempts left", attemptsLeft)
}

func lockedError(seconds int64) error {
	minutes := (seconds + 59) / 60
	if minutes == 1 {
		return errors.New("Too many failed attempts. Try again in 1 minute")
	}
	return fmt.Errorf("Too many failed attempts. Try again in %d minutes", minutes)
}

func getLoginLock(username string) (int64, error) {
	db, openErr := openDb()
	if openErr != nil {
		LogError(GetFile(), openErr)
		return 0, openErr
	}
	defer db.Close()

	var lockedUntil int64
	queryErr := db.QueryRow("SELECT lockedUntil FROM LoginAttemptT WHERE username = ?", username).Scan(&lockedUntil)
	if queryErr != nil {
		if queryErr == sql.ErrNoRows {
			return 0, nil
		}
		LogError(GetFile(), queryErr)
		return 0, queryErr
	}

	remaining := lockedUntil - time.Now().Unix()
	if remaining < 0 {
		return 0, nil
	}
	return remaining, nil
}

func registerFailedLogin(username string) (int, error) {
	db, openErr := openDb()
	if openErr != nil {
		LogError(GetFile(), openErr)
		return 0, openErr
	}
	defer db.Close()

	_, execErr := db.Exec(`
		INSERT INTO LoginAttemptT (username, attempts) VALUES (?, 1)
		ON CONFLICT(username) DO UPDATE SET attempts = attempts + 1`, username)
	if execErr != nil {
		LogError(GetFile(), execErr)
		return 0, execErr
	}

	var attempts int
	queryErr := db.QueryRow("SELECT attempts FROM LoginAttemptT WHERE username = ?", username).Scan(&attempts)
	if queryErr != nil {
		LogError(GetFile(), queryErr)
		return 0, queryErr
	}

	if attempts >= maxLoginAttempts {
		_, lockErr := db.Exec("UPDATE LoginAttemptT SET attempts = 0, lockedUntil = unixepoch() + ? WHERE username = ?", loginLockSeconds, username)
		if lockErr != nil {
			LogError(GetFile(), lockErr)
			return 0, lockErr
		}
		return 0, nil
	}
	return maxLoginAttempts - attempts, nil
}

func resetLoginAttempts(username string) error {
	db, openErr := openDb()
	if openErr != nil {
		LogError(GetFile(), openErr)
		return openErr
	}
	defer db.Close()

	_, execErr := db.Exec("DELETE FROM LoginAttemptT WHERE username = ?", username)
	if execErr != nil {
		LogError(GetFile(), execErr)
		return execErr
	}
	return nil
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

func FuncLogout(sessionToken string) error {
	db, openErr := openDb()
	if openErr != nil {
		LogError(GetFile(), openErr)
		return openErr
	}
	defer db.Close()

	// In der DB steht nur der Hash, nicht der Token selbst
	tokenHash := GenerateHash(sessionToken)
	_, execErr := db.Exec("DELETE FROM SessionT WHERE sessionToken = ?", tokenHash)
	if execErr != nil {
		LogError(GetFile(), execErr)
		return execErr
	}
	return nil
}

func GetUsername(userID int) (string, error) {
	db, openErr := openDb()
	if openErr != nil {
		LogError(GetFile(), openErr)
		return "", openErr
	}
	defer db.Close()

	var username string
	scanErr := db.QueryRow("SELECT username FROM UserT WHERE userID = ?", userID).Scan(&username)
	if scanErr != nil {
		LogError(GetFile(), scanErr)
		return "", scanErr
	}
	return username, nil
}

func DeleteAccount(userID int) error {
	root, rootErr := IsRoot(userID)
	if rootErr != nil {
		return rootErr
	}
	if root {
		return errors.New("Root user can't be deleted")
	}

	db, openErr := openDb()
	if openErr != nil {
		LogError(GetFile(), openErr)
		return openErr
	}
	defer db.Close()

	tx, txErr := db.Begin()
	if txErr != nil {
		LogError(GetFile(), txErr)
		return txErr
	}
	// Wenn Commit schon gelaufen ist, macht Rollback nichts mehr
	defer tx.Rollback()

	// Reihenfolge: zuerst die abhängigen Daten, ganz zum Schluss der User selbst
	queries := []string{
		`DELETE FROM CollectionProgress WHERE collectionID IN (SELECT collectionID FROM Collection WHERE userID = ?)`,
		`DELETE FROM Card WHERE collectionID IN (SELECT collectionID FROM Collection WHERE userID = ?)`,
		`DELETE FROM Collection WHERE userID = ?`,
		`DELETE FROM CollectionGroup WHERE userID = ?`,
		`DELETE FROM SessionT WHERE userID = ?`,
		`DELETE FROM UserT WHERE userID = ?`,
	}

	for _, query := range queries {
		_, execErr := tx.Exec(query, userID)
		if execErr != nil {
			LogError(GetFile(), execErr)
			return execErr
		}
	}

	commitErr := tx.Commit()
	if commitErr != nil {
		LogError(GetFile(), commitErr)
		return commitErr
	}
	return nil
}

func LogoutAllSessions(userID int) error {
	db, openErr := openDb()
	if openErr != nil {
		LogError(GetFile(), openErr)
		return openErr
	}
	defer db.Close()

	_, execErr := db.Exec("DELETE FROM SessionT WHERE userID = ?", userID)
	if execErr != nil {
		LogError(GetFile(), execErr)
		return execErr
	}
	return nil
}
