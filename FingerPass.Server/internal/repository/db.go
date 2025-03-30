package repository

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/google/uuid"

	"github.com/matwate/corner/internal/model"
)

var user = os.Getenv("POSTGRES_USER")
var password = os.Getenv("POSTGRES_PASSWORD")
var dbname = os.Getenv("POSTGRES_DB")
var db_host = os.Getenv("POSTGRES_HOST")
var connectionString string = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", db_host, user, password, dbname)

func InsertUser(template, correo, nombre, programa string) {
	template_uuid := uuid.New()
	template_path := fmt.Sprintf("./templates/%s.bin", template_uuid.String())

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

	f, err := os.Create(template_path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.Write(dec) // Write the decoded template to the file
	if err != nil {
		log.Fatal(err)
	}

	tx := db.MustBegin()

	// Create an user
	tx.MustExec(
		"INSERT INTO Usuarios (correo, nombre, programa) VALUES ($1, $2, $3)",
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

func GetUser(id int) (model.User, error) {
	db, err := sqlx.Connect(
		"postgres",
		connectionString,
	)
	if err != nil {
		return model.User{}, err
	}

	var user model.User
	db.Get(&user, "SELECT * FROM Usuarios WHERE id = $1", id)
	return user, nil
}
