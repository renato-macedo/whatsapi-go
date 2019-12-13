package main

import (
	"fmt"
	"github.com/labstack/echo"
	controllers "github.com/renato-macedo/whatsapi/controllers"
)

func main() {
	fmt.Println("Hello, world!")

	//waconnection.NewConnection("renato")
	e := echo.New()

	e.POST("/session", controllers.CreateSession)
	e.POST("/:id/sendText", controllers.SendText)
	e.Logger.Fatal(e.Start(":1323"))
}
