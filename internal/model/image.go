package model

type Image struct {
	Id       int    `json:"id"       db:"id"`
	Template string `json:"template" db:"path"`
	User_id  int    `json:"user_id"  db:"user_id"`
}
