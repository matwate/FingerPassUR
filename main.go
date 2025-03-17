package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/miqdadyyy/go-sourceafis/templates"
)

type Response struct {
	Message string `json:"message"`
}

type User struct {
	Id       int    `json:"id"`
	Correo   string `json:"correo"`
	Nombre   string `json:"nombre"`
	Programa string `json:"programa"`
}

type Image struct {
	Id       int    `json:"id"       db:"id"`
	Template string `json:"template" db:"path"`
	User_id  int    `json:"user_id"  db:"user_id"`
}

type UserRequest struct {
	Template string `json: template` // Base64 encoded template
	Correo   string `json: correo`
	Nombre   string `json: nombre`
	Programa string `json: programa`
}

func HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	response := Response{Message: "Server is running"}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response.Message))
}

func StoreHash(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body",
			http.StatusInternalServerError)
		return
	}

	var req UserRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Error parsing request body",
			http.StatusBadRequest)
	}

	InsertUser(req.Template, req.Correo, req.Nombre, req.Programa)
	response := Response{
		Message: fmt.Sprintf(
			"Template: %s, Nombre: %s, Correo: %s, Programa: %s",
			req.Template, req.Nombre, req.Correo, req.Programa,
		),
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response.Message))
}

func HandleGetUserFromHash(w http.ResponseWriter, r *http.Request) {
	// Uninmplemented
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Not implemented"))
}

func HandleGetUsers(w http.ResponseWriter, r *http.Request) {
	from, err := strconv.Atoi(r.PathValue("from"))
	if err != nil {
		log.Fatal(err)
	}
	to, err := strconv.Atoi(r.PathValue("to"))
	if err != nil {
		log.Fatal(err)
	}

	UserList := ListUsers(from, to)
	response, err := json.Marshal(UserList)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func HandleInserImage(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body",
			http.StatusInternalServerError)
		return
	}

	var req Image
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Error parsing request body",
			http.StatusBadRequest)
	}

	InsertImage(req.User_id, req.Template)
	response := Response{
		Message: fmt.Sprintf(
			"Template: %s",
			req.Template,
		),
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response.Message))
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/health", HandleHealthCheck)
	router.HandleFunc("POST /store", StoreHash)
	router.HandleFunc("POST /image", HandleInserImage)
	router.HandleFunc("GET /user/{hash}", HandleGetUserFromHash)
	router.HandleFunc("GET /user/{from}/{to}", HandleGetUsers)
	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("Server is running on port 8080")
	var templates []*templates.SearchTemplate
	LoadTemplates(&templates)
	fmt.Println("Templates loaded")
	fmt.Println(templates)

	server.ListenAndServe()
}
