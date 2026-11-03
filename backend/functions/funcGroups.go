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
	defer db.Close()

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
	defer db.Close()

	rows, queryErr := db.Query("SELECT GroupID, name FROM CollectionGroup WHERE userID = (?) ORDER BY GroupID DESC", userID)
	if queryErr != nil {
		LogError(GetFile(), queryErr)
		return nil, queryErr
	}

	var groups []GroupElement
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
	defer db.Close()

	_, migrationErr := db.Exec("UPDATE Collection SET groupID = 0 WHERE groupID = (?) AND userID = (?)", groupID, userID)
	if migrationErr != nil {
		LogError(GetFile(), migrationErr)
	}

	_, err := db.Exec("DELETE FROM CollectionGroup WHERE groupID = (?) AND userID = (?)", groupID, userID)
	if err != nil {
		LogError(GetFile(), err)
		return err
	}
	return nil
}
