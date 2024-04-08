package main

import (
	"fmt"
	"net/http"

	"github.com/julian-bd/dnd/data"

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

func main() {
	err := data.InitDB()
	if err != nil {
		fmt.Println(err)
	}
	router := gin.Default()
	router.GET("/playable_races", getPlayableRaces)
	router.GET("/playable_races/:name", getPlayableRace)
	router.POST("/playable_races", postPlayableRaces)
	router.Run("localhost:8080")
}
