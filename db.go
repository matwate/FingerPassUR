package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var schema = `
create table  Usuarios (
  Hash varchar(255) primary key,
  Correo varchar(255),
  Nombre varchar(255),
  Programa varchar(255)
);`

func InsertUser(hash, correo, nombre, programa string) {
	db, err := sqlx.Connect(
		"postgres",
		"user=postgres password=postgres dbname=FingerPassUR sslmode=disable",
	)
	if err != nil {
		log.Fatal(err)
	}

	tx := db.MustBegin()
	tx.MustExec(
		"INSERT INTO Usuarios (Hash, Correo, Nombre, Programa) VALUES ($1, $2, $3, $4)",
		hash,
		correo,
		nombre,
		programa,
	)
	tx.Commit()
}

func GetUser(hash string) []User {
	db, err := sqlx.Connect(
		"postgres",
		"user=postgres password=postgres dbname=FingerPassUR sslmode=disable",
	)
	if err != nil {
		log.Fatal(err)
	}

	people := []User{}

	db.Select(&people, "SELECT * FROM Usuarios WHERE Hash = $1", hash)
	return people
}
