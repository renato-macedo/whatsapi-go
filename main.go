package main

import (
	"fmt"
	"github.com/labstack/echo"
	model "github.com/renato-macedo/whatsapi/model"
	waconnection "github.com/renato-macedo/whatsapi/waconnection"
	"net/http"
)

func main() {
	fmt.Println("Hello, world")

	//waconnection.NewConnection("renato")
	e := echo.New()

	e.GET("/", createSession)
	e.Logger.Fatal(e.Start(":1323"))
	fmt.Printf("Tudo certo")
}

func createSession(c echo.Context) error {
	sessionName := c.QueryParam("session")
	if sessionName == "" {
		return c.String(http.StatusOK, "O parametro `session` é necessário")
	}
	done := make(chan model.Result)
	go waconnection.NewConnection(sessionName, done)
	fmt.Println("hmmmmmm")
	result := <-done
	if result.Success == true {
		return c.String(http.StatusOK, "Session "+sessionName+" created")
	}
	return c.String(http.StatusOK, result.Message)
}
