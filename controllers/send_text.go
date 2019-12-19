package controllers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	models "github.com/renato-macedo/whatsapi/models"
	utils "github.com/renato-macedo/whatsapi/utils"
	waconnection "github.com/renato-macedo/whatsapi/waconnection"
)

// SendText handles the request to send new text messages
func SendText(c echo.Context) error {
	id := c.Param("id")
	connectionIsActive, err := utils.ConnectionIsActive(waconnection.Connections, id)
	if err != nil {
		response := &models.Response{Success: false, Message: "Erro no servidor"}
		return c.JSON(http.StatusInternalServerError, response)
	}

	// check if the id exists
	if !connectionIsActive {
		response := &models.Response{Success: false, Message: "Conexao não está ativa"}
		return c.JSON(http.StatusInternalServerError, response)
	}

	message := &models.TextMessageDTO{}

	err = c.Bind(message)

	if err != nil {
		fmt.Printf("%v", err)
		return err
	}
	wac := utils.FindConnectionByID(waconnection.Connections, id)
	if wac == nil {
		return c.JSON(http.StatusOK, &models.Response{Success: false, Message: "Could not find connection"})
	}
	err = waconnection.SendTextMessage(wac, message.Number, message.Text)

	if err != nil {
		return c.JSON(http.StatusOK, &models.Response{Success: false, Message: "Something is wrong with the whatsapp"})
	}

	return c.JSON(http.StatusOK, &models.Response{Success: true, Message: "Message sent"})
}
