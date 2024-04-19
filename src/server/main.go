package main

import (
	"log"

	"github.com/gofor-little/env"
	"github.com/julian-bd/dnd/data"
	"github.com/julian-bd/dnd/ripper"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
    if err := env.Load("server/.env"); err != nil {
		log.Fatal(err)
    }
	err := data.InitDB()
	if err != nil {
		log.Fatal(err)
	}

    err = data.CreateTables()
    if err != nil {
        log.Fatal(err)
    }

    ripper.Seed()

	e := echo.New()
	e.Use(middleware.Logger())
	e.File("/css/styles.css", "server/views/css/styles.css")
	e.Renderer = newTemplate()

	e.GET("/", index)
	e.GET("/PlayableRaces/:name", playableRacesName)
	e.GET("/PlayableRaces", playableRaces)
	e.POST("/PlayableRaces", postPlayableRace)
	e.GET("/Traits", traits)

	e.Logger.Fatal(e.Start(":8080"))
}
