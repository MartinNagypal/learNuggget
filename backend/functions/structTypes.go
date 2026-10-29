package functions

type CollectionElement struct {
	Name string
	Id   int
}

type Card struct {
	Question     string `json:"question"`
	Answer       string `json:"answer"`
	CardID       int    `json:"cardId"`
	CollectionID int    `json:"collectionID"`
}

type RegisterUserElement struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
