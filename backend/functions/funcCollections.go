package functions

import "errors"

// ErrCollectionNotFound is returned when a collection does not exist or does
// not belong to the requesting user. Both cases are reported the same way so
// that other users' collection IDs cannot be probed.
var ErrCollectionNotFound = errors.New("Collection not found")

// UserOwnsCollection reports whether the collection with the given ID belongs to the user.
func UserOwnsCollection(collectionID, userID int) (bool, error) {
	db, dbErr := openDb()
	if dbErr != nil {
		LogError(GetFile(), dbErr)
		return false, dbErr
	}
	defer db.Close()

	var exists bool
	queryErr := db.QueryRow("SELECT EXISTS(SELECT 1 FROM Collection WHERE collectionID = ? AND userID = ?)", collectionID, userID).Scan(&exists)
	if queryErr != nil {
		LogError(GetFile(), queryErr)
		return false, queryErr
	}
	return exists, nil
}

func AddNewCollection(collectionName string, userID, groupID int) error {
	db, openErr := openDb()
	if openErr != nil {
		LogError(GetFile(), openErr)
		return openErr
	}
	defer db.Close()

	result, execErr := db.Exec(`
		INSERT INTO Collection (name, userID, groupID)
		SELECT ?, ?, ?
		WHERE ? = 0
		   OR EXISTS (SELECT 1 FROM CollectionGroup WHERE groupID = ? AND userID = ?)`,
		collectionName, userID, groupID,
		groupID,
		groupID, userID,
	)
	if execErr != nil {
		LogError(GetFile(), execErr)
		return execErr
	}

	rows, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		LogError(GetFile(), rowsErr)
		return rowsErr
	}
	if rows == 0 {
		return errors.New("Group not found")
	}

	collectionID, idErr := result.LastInsertId()
	if idErr != nil {
		LogError(GetFile(), idErr)
		return idErr
	}

	_, progressErr := db.Exec("INSERT INTO CollectionProgress (CollectionID) VALUES (?)", collectionID)
	if progressErr != nil {
		LogError(GetFile(), progressErr)
		return progressErr
	}

	return nil
}

func RemoveCollection(collectionID, userID int) error {
	db, openErr := openDb()
	if openErr != nil {
		LogError(GetFile(), openErr)
		return openErr
	}
	defer db.Close()

	owns, ownErr := UserOwnsCollection(collectionID, userID)
	if ownErr != nil {
		return ownErr
	}
	if !owns {
		return ErrCollectionNotFound
	}

	_, cardErr := db.Exec("DELETE FROM Card WHERE collectionID = (?)", collectionID)
	if cardErr != nil {
		LogError(GetFile(), cardErr)
		return cardErr
	}

	_, collErr := db.Exec("DELETE FROM Collection WHERE collectionID = (?) AND userID = (?)", collectionID, userID)
	if collErr != nil {
		LogError(GetFile(), collErr)
		return collErr
	}

	_, collProgErr := db.Exec("DELETE FROM CollectionProgress WHERE CollectionID = (?)", collectionID)
	if collProgErr != nil {
		LogError(GetFile(), collProgErr)
		return collProgErr
	}
	return nil
}

func FetchCollections(userID, groupID int) ([]CollectionElement, error) {
	collections := []CollectionElement{}
	db, openErr := openDb()
	if openErr != nil {
		LogError(GetFile(), openErr)
		return collections, openErr
	}
	defer db.Close()

	rows, queryErr := db.Query("SELECT collectionID, name, groupID FROM Collection WHERE userID = (?) AND groupID = (?) ORDER BY collectionID DESC", userID, groupID)
	if queryErr != nil {
		LogError(GetFile(), queryErr)
		return collections, queryErr

	}
	defer rows.Close()
	for rows.Next() {
		var collection CollectionElement
		if err := rows.Scan(&collection.Id, &collection.Name, &collection.Gid); err != nil {
			LogError(GetFile(), err)
			return collections, err
		}
		collections = append(collections, collection)
	}
	if err := rows.Err(); err != nil {
		LogError(GetFile(), err)
		return collections, err
	}
	return collections, nil
}

func GetCollectionLen(collectionID int) (int, error) {
	db, dbErr := openDb()
	if dbErr != nil {
		LogError(GetFile(), dbErr)
		return 0, dbErr
	}
	defer db.Close()

	var len int
	queryErr := db.QueryRow("SELECT COUNT(*) FROM Card WHERE collectionID = ?", collectionID).Scan(&len)

	if queryErr != nil {
		LogError(GetFile(), queryErr)
		return 0, queryErr
	}

	return len, nil
}

func GetCollectionStats(collectionID int) (CollectionStats, error) {
	db, dbErr := openDb()
	if dbErr != nil {
		LogError(GetFile(), dbErr)
		return CollectionStats{}, dbErr
	}
	defer db.Close()
	var stats CollectionStats
	len, lenErr := GetCollectionLen(collectionID)
	if lenErr != nil {
		LogError(GetFile(), lenErr)
		return CollectionStats{}, lenErr
	}
	stats.Len = len
	markedErr := db.QueryRow("SELECT COUNT(*) FROM Card WHERE collectionID = ? AND marked = 1", collectionID).Scan(&stats.Marked)
	if markedErr != nil {
		LogError(GetFile(), markedErr)
		return CollectionStats{}, markedErr
	}
	execErr := db.QueryRow("SELECT avgScore, lastScore, totalPlayed FROM CollectionProgress WHERE collectionID = (?)", collectionID).Scan(&stats.AvgScore, &stats.LastScore, &stats.TotalPlayed)
	if execErr != nil {
		LogError(GetFile(), execErr)
		return CollectionStats{}, execErr
	}
	return stats, nil
}

func UpdateCollectionStats(stats UpdatedCollectionStatsElement) error {
	db, dbErr := openDb()
	if dbErr != nil {
		LogError(GetFile(), dbErr)
		return dbErr
	}
	defer db.Close()

	_, execErr := db.Exec("UPDATE CollectionProgress SET avgScore = ?, lastScore = ?, totalPlayed = totalPlayed + 1 WHERE collectionID = (?)", stats.AvgScore, stats.LastScore, stats.CollectionID)
	if execErr != nil {
		LogError(GetFile(), execErr)
		return execErr
	}
	return nil
}

func FetchMarkedCollections(userID int) ([]MarkedCollectionElement, error) {
	collections := []MarkedCollectionElement{}
	db, openErr := openDb()
	if openErr != nil {
		LogError(GetFile(), openErr)
		return collections, openErr
	}
	defer db.Close()

	rows, queryErr := db.Query(`
		SELECT col.collectionID, col.name, COUNT(c.cardID)
		FROM Collection col
		JOIN Card c ON c.collectionID = col.collectionID
		WHERE col.userID = ? AND c.marked = 1
		GROUP BY col.collectionID, col.name
		ORDER BY col.name`, userID)
	if queryErr != nil {
		LogError(GetFile(), queryErr)
		return collections, queryErr
	}
	defer rows.Close()

	for rows.Next() {
		var collection MarkedCollectionElement
		if err := rows.Scan(&collection.Id, &collection.Name, &collection.Marked); err != nil {
			LogError(GetFile(), err)
			return collections, err
		}
		collections = append(collections, collection)
	}
	if err := rows.Err(); err != nil {
		LogError(GetFile(), err)
		return collections, err
	}
	return collections, nil
}
