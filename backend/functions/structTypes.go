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

type GroupElement struct {
	Name    string `json:"name"`
	GroupID int    `json:"groupID"`
}
