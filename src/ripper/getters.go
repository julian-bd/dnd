package ripper

import (
	"encoding/json"
	"net/http"
)

func getRace(endpoint string) (raceResponse, error) {
	var r raceResponse
	resp, err := http.Get(baseUrl + endpoint)
	if err != nil {
		return r, nil
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return r, nil
	}
	return r, nil
}

func getResults(endpoint string) (resultsResponse, error) {
	var results resultsResponse
	resp, err := http.Get(baseUrl + endpoint)
	if err != nil {
		return results, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		return results, err
	}
	return results, nil
}
