package model

type User struct {
	Id       int    `json:"id"`
	Correo   string `json:"correo"`
	Nombre   string `json:"nombre"`
	Programa string `json:"programa"`
}
