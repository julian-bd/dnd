package main

import (
	"fmt"

	"github.com/julian-bd/dnd/data"
	"github.com/julian-bd/dnd/ripper"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	err := data.InitDB()
	if err != nil {
		fmt.Println(err)
	}

	if !data.HasBeenSeeded() {
		ripper.Seed()
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.File("/css/styles.css", "server/views/css/styles.css")
	e.Renderer = newTemplate()

	e.GET("/", index)
	e.GET("/PlayableRaces/:name", playableRacesName)
	e.GET("/PlayableRaces", playableRaces)
	e.POST("/PlayableRaces", postPlayableRace)

	e.Logger.Fatal(e.Start(":8080"))
}
