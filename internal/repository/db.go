package repository

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"

  "github.com/matwate/corner/internal/model"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const connectionString string = "user=zryogi password=pPassZr dbname=FingerPassUR sslmode=disable" 

/*

create table Usuarios (
  id serial primary key,
  Correo varchar(255),
  Nombre varchar(255),
  Programa varchar(255)
);


create table Images (
  id serial primary key,
  path varchar(512),
  user_id integer references Usuarios(id)
)

*/

func LoadSchema() string {
	// Load the file query.sql
	file, err := os.ReadFile("query.sql")
	if err != nil {
		log.Fatal(err)
	}
	return string(file)
}

func InsertUser(template, correo, nombre, programa string) {
	db, err := sqlx.Connect(
		"postgres",
    connectionString,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Convert the template base64 to a file and get its path
	dec, err := base64.StdEncoding.DecodeString(template)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("./templates/" + correo + ".bin")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.Write(dec) // Write the decoded template to the file
	if err != nil {
		log.Fatal(err)
	}
	template_path := "./templates/" + correo + ".bin"

	tx := db.MustBegin()

	// Create an user
	tx.MustExec(
		"INSERT INTO Usuarios (Correo, Nombre, Programa) VALUES ($1, $2, $3)",
		correo,
		nombre,
		programa,
	)

	// Get the user id
	var id int
	tx.Get(&id, "SELECT id FROM Usuarios WHERE Correo = $1", correo)

	// Insert the template path
	tx.MustExec(
		"INSERT INTO Images (path, user_id) VALUES ($1, $2)",
		template_path,
		id,
	)

	tx.Commit()
}

func InsertImage(user_id int, path string) {
	db, err := sqlx.Connect(
		"postgres",
		connectionString,
	)
	if err != nil {
		log.Fatal(err)
	}

	tx := db.MustBegin()
	tx.MustExec(
		"INSERT INTO Images (path, user_id) VALUES ($1, $2)",
		path,
		user_id,
	)
	tx.Commit()
}

func ListUsers(from, to int) []model.User {
	db, err := sqlx.Connect(
		"postgres",
		connectionString,
	)
	if err != nil {
		log.Fatal(err)
	}

	var users []model.User
	db.Select(&users, "SELECT * FROM Usuarios where id between $1 and $2", from, to)
	return users
}

func ListALLImages() []model.Image {
	db, err := sqlx.Connect(
		"postgres",
		connectionString,
	)
	if err != nil {
		log.Fatal(err)
	}

	var images []model.Image
	db.Select(&images, "SELECT * FROM Images")
	fmt.Println(images)
	return images
}
