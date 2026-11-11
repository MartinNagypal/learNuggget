package functions

import "database/sql"

func IsRoot(userID int) (bool, error) {
	db, openErr := openDb()
	if openErr != nil {
		LogError(GetFile(), openErr)
		return false, openErr
	}

	var root bool
	queryErr := db.QueryRow("SELECT root FROM UserT WHERE userID = ?", userID).Scan(&root)
	if queryErr != nil {
		if queryErr == sql.ErrNoRows {
			return false, nil
		}
		LogError(GetFile(), queryErr)
		return false, queryErr
	}
	return root, nil
}

func GetAdminStats() (AdminStats, error) {
	db, openErr := openDb()
	if openErr != nil {
		LogError(GetFile(), openErr)
		return AdminStats{}, openErr
	}

	var stats AdminStats
	queryErr := db.QueryRow(`
		SELECT
			(SELECT COUNT(*) FROM UserT),
			(SELECT COUNT(*) FROM Collection),
			(SELECT COUNT(*) FROM CollectionGroup),
			(SELECT COUNT(*) FROM Card)`).Scan(&stats.Users, &stats.Collections, &stats.Groups, &stats.Cards)
	if queryErr != nil {
		LogError(GetFile(), queryErr)
		return AdminStats{}, queryErr
	}
	return stats, nil
}

func FetchAllUsers() ([]AdminUserElement, error) {
	users := []AdminUserElement{}
	db, openErr := openDb()
	if openErr != nil {
		LogError(GetFile(), openErr)
		return users, openErr
	}

	rows, queryErr := db.Query(`
		SELECT u.userID, u.username, u.root,
			(SELECT COUNT(*) FROM Collection col WHERE col.userID = u.userID),
			(SELECT COUNT(*) FROM Card c JOIN Collection col ON c.collectionID = col.collectionID WHERE col.userID = u.userID),
			(SELECT MAX(s.createdAt) FROM SessionT s WHERE s.userID = u.userID)
		FROM UserT u
		ORDER BY u.userID`)
	if queryErr != nil {
		LogError(GetFile(), queryErr)
		return users, queryErr
	}
	defer rows.Close()

	for rows.Next() {
		var user AdminUserElement
		var lastActive sql.NullInt64
		if err := rows.Scan(&user.UserID, &user.Username, &user.Root, &user.Collections, &user.Cards, &lastActive); err != nil {
			LogError(GetFile(), err)
			return users, err
		}
		user.LastActive = lastActive.Int64
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		LogError(GetFile(), err)
		return users, err
	}
	return users, nil
}
