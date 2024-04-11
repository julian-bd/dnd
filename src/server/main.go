package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/julian-bd/dnd/data"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/gin-gonic/gin"
)

func getPlayableRaces(c *gin.Context) {
	result, err := data.PlayableRaceNames()
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "playable race names not found"})
		c.AbortWithError(http.StatusConflict, err)
	}
	c.IndentedJSON(http.StatusOK, result)
}

func getPlayableRace(c *gin.Context) {
	name := c.Param("name")
	result, err := data.PlayableRaceByName(name)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "playable race not found"})
		c.AbortWithError(http.StatusConflict, err)
	}
	c.IndentedJSON(http.StatusOK, result)
}

func postPlayableRaces(c *gin.Context) {
	var newPlayableRace data.PlayableRace
	if err := c.BindJSON(&newPlayableRace); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	_, err := data.InsertPlayableRace(newPlayableRace)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "saving playable race failed"})
		c.AbortWithError(http.StatusConflict, err)
	}
	c.IndentedJSON(http.StatusCreated, newPlayableRace)
}

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("server/views/*.html")),
	}
}

type indexContent struct {
	PlayableRaceNames []string
}

func main() {
	err := data.InitDB()
	if err != nil {
		fmt.Println(err)
	}
	//router := gin.Default()
	//router.GET("/api/playable_races", getPlayableRaces)
	//router.GET("/api/playable_races/:name", getPlayableRace)
	//router.POST("/api/playable_races", postPlayableRaces)
	//router.Run("localhost:8080")

	playableRaceNames, err := data.PlayableRaceNames()
	if err != nil {
		fmt.Println(err)
	}
	content := indexContent{
		PlayableRaceNames: playableRaceNames,
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Renderer = newTemplate()

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", content)
	})

	e.GET("/PlayableRaces/:name", func(c echo.Context) error {
		name := c.Param("name")
		data, err := data.PlayableRaceByName(name)
		if err != nil {
			e.Logger.Error(err)
			return c.Redirect(404, "not_found")
		}
		return c.Render(200, "playable_race", data)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
