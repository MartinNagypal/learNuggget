package functions

import (
	"errors"
	"strings"
)

func CreateGroup(name string, userID int) error {
	db, openErr := openDb()
	if openErr != nil {
		LogError(GetFile(), openErr)
		return openErr
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("name is empty")
	}

	_, queryErr := db.Exec("INSERT INTO CollectionGroup (name, userID) VALUES (?,?)", name, userID)
	if queryErr != nil {
		LogError(GetFile(), queryErr)
		return queryErr
	}
	return nil
}

func GetGroups(userID int) ([]GroupElement, error) {
	db, openErr := openDb()
	if openErr != nil {
		LogError(GetFile(), openErr)
		return nil, openErr
	}

	rows, queryErr := db.Query("SELECT GroupID, name FROM CollectionGroup WHERE userID = (?) ORDER BY GroupID DESC", userID)
	if queryErr != nil {
		LogError(GetFile(), queryErr)
		return nil, queryErr
	}
	defer rows.Close()
	// empty list instead of nil, so the JSON is [] and not null when the user has no groups
	groups := []GroupElement{}
	for rows.Next() {
		var group GroupElement
		if scanErr := rows.Scan(&group.GroupID, &group.Name); scanErr != nil {
			LogError(GetFile(), scanErr)
			return nil, scanErr
		}
		groups = append(groups, group)
	}
	return groups, nil
}

func DelGroup(userID, groupID int) error {
	db, openErr := openDb()
	if openErr != nil {
		LogError(GetFile(), openErr)
		return openErr
	}

	// Transaction: collections are only moved out if the group is really deleted
	tx, txErr := db.Begin()
	if txErr != nil {
		LogError(GetFile(), txErr)
		return txErr
	}
	defer tx.Rollback()

	result, err := tx.Exec("DELETE FROM CollectionGroup WHERE groupID = (?) AND userID = (?)", groupID, userID)
	if err != nil {
		LogError(GetFile(), err)
		return err
	}
	deleted, affectedErr := result.RowsAffected()
	if affectedErr != nil {
		LogError(GetFile(), affectedErr)
		return affectedErr
	}
	if deleted == 0 {
		return ErrGroupNotFound
	}

	_, migrationErr := tx.Exec("UPDATE Collection SET groupID = 0 WHERE groupID = (?) AND userID = (?)", groupID, userID)
	if migrationErr != nil {
		LogError(GetFile(), migrationErr)
		return migrationErr
	}

	commitErr := tx.Commit()
	if commitErr != nil {
		LogError(GetFile(), commitErr)
		return commitErr
	}
	return nil
}
