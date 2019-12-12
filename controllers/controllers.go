package controllers

import (
	"fmt"
	"github.com/labstack/echo"
	models "github.com/renato-macedo/whatsapi/models"
	waconnection "github.com/renato-macedo/whatsapi/waconnection"
	"net/http"
	"os"
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
	// obtendo a lista de nomes do arquivos que existem na pasta sessions
	folder, err := os.Open("./sessions")
	files, err := folder.Readdirnames(-1)

	// verificando pelo nome se a sessao já existe
	for _, filename := range files {
		if session.Name+".gob" == filename {
			// se existir entao cria-se um objeto para se enviar na resposta
			response := &models.Response{Success: false, Message: "Esta sessão já existe"}
			return c.JSON(http.StatusOK, response)
		}
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
	message := &models.Message{}

	err := c.Bind(message)
	if err != nil {
		fmt.Printf("%v", err)
		return err
	}
	return c.JSON(http.StatusOK, message)
}

// SendImage handles the request to send a image
//func SendImage(c echo.Context) error {}
