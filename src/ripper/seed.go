package ripper

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julian-bd/dnd/data"
)

func seedTraits() error {
	results, err := getResults("/api/traits")
	if err != nil {
		return err
	}
	// TODO: This really should be a batch insert
	for _, t := range results.Results {
        resp, err := http.Get(baseUrl + t.Url)
        if err != nil {
            return err
        }
        defer resp.Body.Close()
        var trait traitResponse
        err = json.NewDecoder(resp.Body).Decode(&trait)
		if err != nil {
			return err
		}
        description := strings.Join(trait.Desc, "\n")
		err = data.InsertTrait(trait.Name, description)
		if err != nil {
			return err
		}
	}
	return nil
}

func seedProficiencies() error {
	results, err := getResults("/api/proficiencies")
	if err != nil {
		return err
	}
	// TODO: This really should be a batch insert
	for _, t := range results.Results {
		err := data.InsertProficiency(t.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

func seedLanguages() error {
	results, err := getResults("/api/languages")
	if err != nil {
		return err
	}
	// TODO: This really should be a batch insert
	for _, t := range results.Results {
		err := data.InsertLanguage(t.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

func seedRaces(endpoint string) error {
	es, err := getResults(endpoint)
	if err != nil {
		return err
	}
	var domainRaces []data.PlayableRace
	for _, e := range es.Results {
		r, err := getRace(e.Url)
		if err != nil {
			return err
		}
		domainRace := toDomainRace(r)
		domainRaces = append(domainRaces, domainRace)
	}
	for _, r := range domainRaces {
		_, err := data.InsertPlayableRace(r)
		if err != nil {
			return err
		}
	}
	return nil
}
