package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
  "github.com/matwate/corner/internal/repository"

  "github.com/matwate/corner/internal/dto"
  "github.com/matwate/corner/internal/model"
)

func HandleCreateNewUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body",
			http.StatusInternalServerError)
		return
	}

	var req dto.UserRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Error parsing request body",
			http.StatusBadRequest)
	}

	repository.InsertUser(req.Template, req.Correo, req.Nombre, req.Programa)
	response := dto.Response{
		Message: fmt.Sprintf(
			"Template: %s, Nombre: %s, Correo: %s, Programa: %s",
			req.Template, req.Nombre, req.Correo, req.Programa,
		),
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response.Message))
}

func HandleGetUserFromFpPrint(w http.ResponseWriter, r *http.Request) {
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

	UserList := repository.ListUsers(from, to)
	response, err := json.Marshal(UserList)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func HandleCreateImage(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body",
			http.StatusInternalServerError)
		return
	}

	var req model.Image
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Error parsing request body",
			http.StatusBadRequest)
	}

	repository.InsertImage(req.User_id, req.Template)
	response := dto.Response{
		Message: fmt.Sprintf(
			"Template: %s",
			req.Template,
		),
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response.Message))
}
