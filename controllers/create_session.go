package controllers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	models "github.com/renato-macedo/whatsapi/models"
	waconnection "github.com/renato-macedo/whatsapi/waconnection"
)

// CreateSession handles the request for a new session
func CreateSession(c echo.Context) error {
	// criando uma "instancia" do objeto que virá no request
	SessionDTO := &models.Session{}
	err := c.Bind(SessionDTO)
	if err != nil {
		fmt.Printf("%v", err)
		return err
	}
	// sessionExists, err := utils.SessionExists(SessionDTO.Name)
	// if err != nil {
	// 	response := &models.Response{Success: false, Message: "Erro no servidor"}
	// 	return c.JSON(http.StatusInternalServerError, response)
	// }
	// if sessionExists {
	// 	// se existir entao cria-se um objeto para se enviar na resposta
	// 	response := &models.Response{Success: false, Message: "Esta sessão já existe"}
	// 	return c.JSON(http.StatusOK, response)
	// }

	done := make(chan models.Result)
	go waconnection.NewConnection(SessionDTO.Name, done)

	result := <-done
	if result.Success == true {
		response := &models.Response{Success: true, Message: result.Message}
		return c.JSON(http.StatusCreated, response)
	}
	response := &models.Response{Success: false, Message: result.Message}
	return c.JSON(http.StatusOK, response)
}
