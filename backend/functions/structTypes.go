package functions

type CollectionElement struct {
	Name string
	Id   int
	Gid  int
}

type Card struct {
	Question     string `json:"question"`
	Answer       string `json:"answer"`
	CardID       int    `json:"cardId"`
	CollectionID int    `json:"collectionID"`
}

type UserElement struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLogin struct {
	Username string
	Password string
	id       int64
}

type CollectionStats struct {
	Len         int     `json:"len"`
	AvgScore    float64 `json:"avgScore"`
	LastScore   float64 `json:"lastScore"`
	TotalPlayed int     `json:"totalPlayed"`
}

type UpdatedCollectionStatsElement struct {
	AvgScore     float64 `json:"avgScore"`
	LastScore    float64 `json:"lastScore"`
	CollectionID int     `json:"collectionID"`
}

type GroupElement struct {
	Name    string `json:"name"`
	GroupID int    `json:"groupID"`
}

type RankElement struct {
	Username string  `json:"username"`
	AvgScore float64 `json:"avgScore"`
}
