package functions

func GetUserAvgScore(userID int) (float64, error) {
	db, dbErr := openDb()
	if dbErr != nil {
		LogError(GetFile(), dbErr)
		return 0, dbErr
	}

	var avg float64
	queryErr := db.QueryRow(`
		SELECT COALESCE(
			SUM(p.avgScore * p.totalPlayed) / NULLIF(SUM(p.totalPlayed), 0),
			0)
		FROM CollectionProgress p
		JOIN Collection c ON c.collectionID = p.collectionID
		WHERE c.userID = (?)
		  AND p.totalPlayed >= 3
		  AND (SELECT COUNT(*) FROM Card k WHERE k.collectionID = p.collectionID) >= 5`,
		userID).Scan(&avg)
	if queryErr != nil {
		LogError(GetFile(), queryErr)
		return 0, queryErr
	}
	return avg, nil
}

func GetTopUsers() ([]RankElement, error) {
	ranking := []RankElement{}
	db, dbErr := openDb()
	if dbErr != nil {
		LogError(GetFile(), dbErr)
		return ranking, dbErr
	}

	rows, queryErr := db.Query(`
		SELECT u.username,
		       SUM(p.avgScore * p.totalPlayed) / SUM(p.totalPlayed) AS score
		FROM CollectionProgress p
		JOIN Collection c ON c.collectionID = p.collectionID
		JOIN UserT u ON u.userID = c.userID
		WHERE p.totalPlayed >= 3
		  AND (SELECT COUNT(*) FROM Card k WHERE k.collectionID = p.collectionID) >= 5
		GROUP BY u.userID
		ORDER BY score DESC
		LIMIT 5`)
	if queryErr != nil {
		LogError(GetFile(), queryErr)
		return ranking, queryErr
	}
	defer rows.Close()

	for rows.Next() {
		var entry RankElement
		if err := rows.Scan(&entry.Username, &entry.AvgScore); err != nil {
			LogError(GetFile(), err)
			return ranking, err
		}
		ranking = append(ranking, entry)
	}
	if err := rows.Err(); err != nil {
		LogError(GetFile(), err)
		return ranking, err
	}
	return ranking, nil
}
