package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/joho/godotenv/autoload"
	"github.com/miqdadyyy/go-sourceafis/config"
	"github.com/miqdadyyy/go-sourceafis/templates"

	"github.com/matwate/corner/internal/api"
	"github.com/matwate/corner/internal/service"
	u "github.com/matwate/corner/internal/utils"
)

func main() {
	config.LoadDefaultConfig()
	var Templates []*templates.SearchTemplate
	service.LoadTemplates(&Templates)
	log.Print("Templates loaded")
	log.Print(Templates)
	service.Templates = Templates
	router := http.NewServeMux()
	router.HandleFunc("/health", api.HandleHealthCheck)
	router.HandleFunc("POST /user", api.HandleCreateNewUser)
	router.HandleFunc("POST /image", api.HandleCreateImage)
	router.HandleFunc("POST /user/fetch", api.HandleGetUserFromFpPrint)
	router.HandleFunc("GET /user/{from}/{to}", api.HandleGetUsers)
	server := http.Server{
		Addr:    fmt.Sprintf(":%s", u.GetEnv("PORT", "8080")),
		Handler: router,
	}

	log.Print("Server is running on port 8080")
	log.Fatal(server.ListenAndServe())
}
