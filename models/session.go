package models

// Session é objeto para ser enviado no corpo da requisição de criar sessão
type Session struct {
	Name string `json:"name"`
}
