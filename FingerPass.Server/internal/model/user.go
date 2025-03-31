package model

type User struct {
	Id       int    `db:"id" json:"id"`
	Correo   string `db:"correo" json:"correo"`
	Nombre   string `db:"nombre" json:"nombre"`
	Programa string `db:"programa" json:"programa"`
}
