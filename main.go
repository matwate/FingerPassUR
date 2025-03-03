package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

type User struct {
	Hash     string `json:"hash"`
	Correo   string `json:"correo"`
	Nombre   string `json:"nombre"`
	Programa string `json:"programa"`
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

	var req User
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Error parsing request body",
			http.StatusBadRequest)
	}

	InsertUser(req.Hash, req.Correo, req.Nombre, req.Programa)
	response := Response{
		Message: fmt.Sprintf(
			"Hash: %s, Nombre: %s, Correo: %s, Programa: %s",
			req.Hash, req.Nombre, req.Correo, req.Programa,
		),
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response.Message))
}

func HandleGetUserFromHash(w http.ResponseWriter, r *http.Request) {
	hash := r.PathValue("hash")
	user := GetUser(hash)[0]
	json, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Error parsing request body",
			http.StatusInternalServerError)
	}
	response := Response{
		Message: string(json),
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response.Message))
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/health", HandleHealthCheck)
	router.HandleFunc("POST /store", StoreHash)
	router.HandleFunc("GET /user/{hash}", HandleGetUserFromHash)
	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("Server is running on port 8080")
	server.ListenAndServe()
}
