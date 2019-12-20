package main

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	controllers "github.com/renato-macedo/whatsapi/controllers"
)

func main() {
	fmt.Println("Hello, world!")

	//waconnection.NewConnection("renato")
	e := echo.New()
	e.Use(middleware.CORS())

	e.POST("/session", controllers.CreateSession)
	e.POST("/:id/text", controllers.SendText)
	e.POST("/:id/image", controllers.SendImage)
	e.POST("/:id/audio", controllers.SendAudio)
	e.GET("/connections", controllers.GetConnections)
	e.Logger.Fatal(e.Start(":1323"))
}
