package models

type Todo struct {
	Id          *int64 `json:"id"`
	Description string `json:"description"`
	Status      int    `json:"status"`
}
