package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/julian-bd/dnd/data"
	"github.com/labstack/echo/v4"
)

func headers() headerContent {
	playableRaceNames, err := data.PlayableRaceNames()
	if err != nil {
		fmt.Println(err)
	}

	return headerContent{
		PlayableRaceNames: playableRaceNames,
	}
}

func index(c echo.Context) error {
		content := indexContent{
			HeaderContent: headers(),
		}
		return c.Render(200, "index", content)
}

func playableRacesName(c echo.Context) error {
    name := c.Param("name")
    data, err := data.PlayableRaceByName(name)
    if err != nil {
        fmt.Println(err)
        return c.Redirect(404, "not_found")
    }
    content := playableRaceContent{
        HeaderContent:  headers(),
        AbilityBonuses: data.AbilityBonuses,
        Languages:      data.StartingLanguages,
        Name:           data.Name,
        Proficiencies:  data.StartingProficiencies,
        SubRaces:       data.SubRaces,
        Traits:         data.Traits,
        Speed:          data.Speed,
    }

    return c.Render(200, "playable_race", content)
}

func playableRaces(c echo.Context) error {
    l, err := data.LanguageNames()
    if err != nil {
        fmt.Println(err)
        return c.String(http.StatusBadRequest, "getting language names failed")
    }
    a, err := data.AbilityNames()
    if err != nil {
        fmt.Println(err)
        return c.String(http.StatusBadRequest, "getting ability names failed")
    }
    t, err := data.TraitNames()
    if err != nil {
        fmt.Println(err)
        return c.String(http.StatusBadRequest, "getting trait names failed")
    }
    p, err := data.ProficiencyNames()
    if err != nil {
        fmt.Println(err)
        return c.String(http.StatusBadRequest, "getting proficiency names failed")
    }
    s, err := data.PlayableRaceNames()
    if err != nil {
        fmt.Println(err)
        return c.String(http.StatusBadRequest, "getting playable race names failed")
    }
    d := createPlayableRaceContent{
        HeaderContent: headers(),
        Abilities:     a,
        Languages:     l,
        Proficiencies: p,
        SubRaces:      s,
        Traits:        t,
    }
    return c.Render(200, "add_race", d)
}

func postPlayableRace(c echo.Context) error {
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
