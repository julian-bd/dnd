package main

import (
	"fmt"
	"net/http"
	"strings"

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

	playableRaceNames, err := data.PlayableRaceNames()
	if err != nil {
		fmt.Println(err)
	}

	headerContent := headerContent{
		PlayableRaceNames: playableRaceNames,
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.File("/css/styles.css", "server/views/css/styles.css")
	e.Renderer = newTemplate()

	e.GET("/", func(c echo.Context) error {
		content := indexContent{
			HeaderContent: headerContent,
		}
		return c.Render(200, "index", content)
	})

	e.GET("/PlayableRaces/:name", func(c echo.Context) error {
		name := c.Param("name")
		data, err := data.PlayableRaceByName(name)
		if err != nil {
			e.Logger.Error(err)
			return c.Redirect(404, "not_found")
		}
		content := playableRaceContent{
			HeaderContent:  headerContent,
			AbilityBonuses: data.AbilityBonuses,
			Languages:      data.StartingLanguages,
			Name:           data.Name,
			Proficiencies:  data.StartingProficiencies,
			SubRaces:       data.SubRaces,
			Traits:         data.Traits,
			Speed:          data.Speed,
		}

		return c.Render(200, "playable_race", content)
	})

	e.GET("/PlayableRaces", func(c echo.Context) error {
		l, err := data.LanguageNames()
		a, err := data.AbilityNames()
		t, err := data.TraitNames()
		p, err := data.ProficiencyNames()
		s, err := data.PlayableRaceNames()
		if err != nil {
			// TODO: Something
		}
		d := createPlayableRaceContent{
			HeaderContent: headerContent,
			Abilities:     a,
			Languages:     l,
			Proficiencies: p,
			SubRaces:      s,
			Traits:        t,
		}
		return c.Render(200, "add_race", d)
	})

	e.POST("/PlayableRaces", func(c echo.Context) error {
		var request createPlayableRace
		err := c.Bind(&request)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, "binding failed")
		}
		fmt.Println(request)

		pr, err := playableRace(request)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, "mapping failed")
		}
		fmt.Println(pr)
		_, err = data.InsertPlayableRace(pr)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusInternalServerError, "data saving failed")
		}

		n, err := data.PlayableRaceByName(pr.Name)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusInternalServerError, "data retrieval failed")
		}
		return c.Render(200, "/PlayableRaces/"+pr.Name, n)
	})

	e.Logger.Fatal(e.Start(":8080"))
}

func playableRace(request createPlayableRace) (data.PlayableRace, error) {
	var pr data.PlayableRace
	pr.Name = request.Name
	pr.Speed = request.Speed
	pr.StartingLanguages = strings.Split(request.Languages, ",")
	pr.StartingProficiencies = strings.Split(request.Proficiencies, ",")
	pr.Traits = strings.Split(request.Traits, ",")
	pr.SubRaces = strings.Split(request.SubRaces, ",")
	// TODO: ability bonuses
	return pr, nil
}
