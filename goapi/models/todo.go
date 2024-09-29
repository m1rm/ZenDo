package models

type Todo struct {
	Id          *int64 `json:"id"`
	Description string `json:"text"`
	Status      int    `json:"status"`
}
