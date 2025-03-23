package main

import (
	"log"
	"net/http"
  
  "github.com/matwate/corner/internal/api"
  "github.com/matwate/corner/internal/service"

	"github.com/miqdadyyy/go-sourceafis/config"
	"github.com/miqdadyyy/go-sourceafis/templates"
)

func main() {
	config.LoadDefaultConfig()

	var templates []*templates.SearchTemplate
	service.LoadTemplates(&templates)
	log.Print("Templates loaded")
	log.Print(templates)

	router := http.NewServeMux()
	router.HandleFunc("/health", api.HandleHealthCheck)
	router.HandleFunc("POST /user", api.HandleCreateNewUser)
	router.HandleFunc("POST /image", api.HandleCreateImage)
	router.HandleFunc("GET /user/{hash}", api.HandleGetUserFromFpPrint)
	router.HandleFunc("GET /user/{from}/{to}", api.HandleGetUsers)
	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Print("Server is running on port 8080")
	log.Fatal(server.ListenAndServe())
}
