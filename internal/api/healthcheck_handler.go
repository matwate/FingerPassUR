package api

import (
	"net/http"

  "github.com/matwate/corner/internal/dto"
)


func HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	response := dto.Response{Message: "Server is running"}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response.Message))
}
