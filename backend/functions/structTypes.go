package functions

type CollectionElement struct {
	Name string
	Id   int
	Gid  int
}

type MarkedCollectionElement struct {
	Name   string `json:"name"`
	Id     int    `json:"id"`
	Marked int    `json:"marked"`
}

type AdminStats struct {
	Users       int `json:"users"`
	Collections int `json:"collections"`
	Groups      int `json:"groups"`
	Cards       int `json:"cards"`
}

type AdminUserElement struct {
	UserID      int    `json:"userID"`
	Username    string `json:"username"`
	Root        bool   `json:"root"`
	Collections int    `json:"collections"`
	Cards       int    `json:"cards"`
	LastActive  int64  `json:"lastActive"`
}

type Card struct {
	Question     string `json:"question"`
	Answer       string `json:"answer"`
	CardID       int    `json:"cardId"`
	CollectionID int    `json:"collectionID"`
	Marked       bool   `json:"marked"`
}

type MarkCardRequest struct {
	CardID int  `json:"cardId"`
	Marked bool `json:"marked"`
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
	Marked      int     `json:"marked"`
	AvgScore    float64 `json:"avgScore"`
	LastScore   float64 `json:"lastScore"`
	TotalPlayed int     `json:"totalPlayed"`
}

// LearnResult is what the client reports after finishing a learning session.
// Only raw counts are accepted; the score is computed on the server.
type LearnResult struct {
	CollectionID int `json:"collectionID"`
	Mode         int `json:"mode"`
	RightCnt     int `json:"rightCnt"`
	WrongCnt     int `json:"wrongCnt"`
}

// LearnResultResponse is returned after a learning session was recorded.
type LearnResultResponse struct {
	Score             float64 `json:"score"`
	PreviousLastScore float64 `json:"previousLastScore"`
	PreviousPlayed    int     `json:"previousPlayed"`
}

type GroupElement struct {
	Name    string `json:"name"`
	GroupID int    `json:"groupID"`
}

type RankElement struct {
	Username string  `json:"username"`
	AvgScore float64 `json:"avgScore"`
}
