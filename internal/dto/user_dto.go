package dto

type UserRequest struct {
	Template string `json:"template"` // Base64 encoded template
	Correo   string `json:"correo"`
	Nombre   string `json:"nombre"`
	Programa string `json:"programa"`
}
