package controllers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	models "github.com/renato-macedo/whatsapi/models"
	utils "github.com/renato-macedo/whatsapi/utils"
	waconnection "github.com/renato-macedo/whatsapi/waconnection"
)

// CreateSession handles the request for a new session
func CreateSession(c echo.Context) error {
	// criando uma "instancia" do objeto que virá no request
	session := &models.Session{}
	err := c.Bind(session)
	if err != nil {
		fmt.Printf("%v", err)
		return err
	}
	sessionExists, err := utils.SessionExists(session.Name)
	if err != nil {
		response := &models.Response{Success: false, Message: "Erro no servidor"}
		return c.JSON(http.StatusInternalServerError, response)
	}
	if sessionExists {
		// se existir entao cria-se um objeto para se enviar na resposta
		response := &models.Response{Success: false, Message: "Esta sessão já existe"}
		return c.JSON(http.StatusOK, response)
	}

	done := make(chan models.Result)
	go waconnection.NewConnection(session.Name, done)

	result := <-done
	if result.Success == true {
		response := &models.Response{Success: true, Message: "Sessão criada!"}
		return c.JSON(http.StatusCreated, response)
	}
	response := &models.Response{Success: false, Message: result.Message}
	return c.JSON(http.StatusOK, response)
}

// SendText handles the request to send new text messages
func SendText(c echo.Context) error {
	id := c.Param("id")
	sessionExists, err := utils.SessionExists(id)
	if err != nil {
		response := &models.Response{Success: false, Message: "Erro no servidor"}
		return c.JSON(http.StatusInternalServerError, response)
	}

	// check if the id exists
	if !sessionExists {
		response := &models.Response{Success: false, Message: "Sessão não existe"}
		return c.JSON(http.StatusInternalServerError, response)
	}

	message := &models.Message{}

	err = c.Bind(message)

	if err != nil {
		fmt.Printf("%v", err)
		return err
	}
	waconnection.
		waconnection.SendTextMessage(wac, number, text)

	return c.JSON(http.StatusOK, message)
}

// SendImage handles the request to send a image
//func SendImage(c echo.Context) error {}
