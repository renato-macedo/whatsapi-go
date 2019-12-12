package models

// Response é o objeto para ser enviado nas respostas das requisições
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
