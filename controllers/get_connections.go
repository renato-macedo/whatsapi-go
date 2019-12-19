package controllers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	waconnection "github.com/renato-macedo/whatsapi/waconnection"
)

// GetConnections return all conections must be deleted later
func GetConnections(c echo.Context) error {
	var connectionNames []string
	for _, conection := range waconnection.Connections {
		connectionNames = append(connectionNames, conection.Info.Wid)
	}
	response := fmt.Sprintf("conexoes %v", connectionNames)
	return c.String(http.StatusOK, response)
}
