package controllers

import (
	"fmt"
	"net/http"

	"os"

	"github.com/labstack/echo"
	connections "github.com/renato-macedo/whatsapi/connections"
	models "github.com/renato-macedo/whatsapi/models"
	utils "github.com/renato-macedo/whatsapi/utils"
	waconnection "github.com/renato-macedo/whatsapi/waconnection"
)

// SendImage handles the request to send a image
func SendImage(c echo.Context) error {
	message := &models.ImageMessageDTO{}

	err := c.Bind(message)

	if err != nil {
		fmt.Printf("%v", err)
		return err
	}

	id := c.Param("id")
	connectionIsActive, err := utils.ConnectionIsActive(connections.Connections, id)
	if err != nil {
		response := &models.Response{Success: false, Message: "Erro no servidor"}
		return c.JSON(http.StatusInternalServerError, response)
	}

	// check if the id exists
	if !connectionIsActive {
		response := &models.Response{Success: false, Message: "Conexao não está ativa"}
		return c.JSON(http.StatusInternalServerError, response)
	}

	wac := utils.FindConnectionByID(connections.Connections, id)
	if wac == nil {
		return c.JSON(http.StatusOK, &models.Response{Success: false, Message: "Could not find connection"})
	}
	// criando channel para o tipo do arquivo
	ch := make(chan *os.File)

	for _, URL := range message.URLs {
		fmt.Println(URL)
		go utils.DownloadImage(URL, ch)
	}

	hasErrors := false
	for range message.URLs {
		err = waconnection.SendImageMessage(wac, message.Number, <-ch, "")
		if err != nil {
			hasErrors = true
		}
	}

	//img, err := utils.DownloadImage(message.)
	// if err != nil {
	// 	return err
	// }
	//img, err := os.Open("img/example1.jpeg")
	//err = waconnection.SendImageMessage(wac, message.Number, img, "🥱")

	if hasErrors {
		return c.JSON(http.StatusOK, &models.Response{Success: false, Message: "Some messages were not sent"})
	}

	return c.JSON(http.StatusOK, &models.Response{Success: true, Message: "All messages sent"})
}
