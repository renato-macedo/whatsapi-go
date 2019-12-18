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
	SessionDTO := &models.Session{}
	err := c.Bind(SessionDTO)
	if err != nil {
		fmt.Printf("%v", err)
		return err
	}
	sessionExists, err := utils.SessionExists(SessionDTO.Name)
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
	go waconnection.NewConnection(SessionDTO.Name, done)

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

// GetConnections return all conections must be deleted later
func GetConnections(c echo.Context) error {
	var connectionNames []string
	for _, conection := range waconnection.Connections {
		connectionNames = append(connectionNames, conection.Info.Wid)
	}
	response := fmt.Sprintf("conexoes %v", connectionNames)
	return c.String(http.StatusOK, response)
}

// SendImage handles the request to send a image
func SendImage(c echo.Context) error {
	message := &models.ImageMessageDTO{}

	err := c.Bind(message)

	if err != nil {
		fmt.Printf("%v", err)
		return err
	}

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

	wac := utils.FindConnectionByID(waconnection.Connections, id)
	if wac == nil {
		return c.JSON(http.StatusOK, &models.Response{Success: false, Message: "Could not find connection"})
	}
	img, err := utils.DownloadImage(message.URL)
	if err != nil {
		return err
	}
	//img, err := os.Open("img/example1.jpeg")
	err = waconnection.SendImageMessage(wac, message.Number, img, "🥱")

	if err != nil {
		return c.JSON(http.StatusOK, &models.Response{Success: false, Message: "Something is wrong with the whatsapp"})
	}

	return c.JSON(http.StatusOK, &models.Response{Success: true, Message: "Message sent"})
}
