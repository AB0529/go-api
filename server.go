package main

import (
	"encoding/json"
	"fmt"

	"net/http"
)

// Home the home route of the API
func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Home")
}

// SendJSON will send a JSON response to the API according the Response type
func SendJSON(w http.ResponseWriter, resp Response) (bool, error) {
	// Decode JSON
	json, err := json.Marshal(resp)

	if err != nil {
		return false, err
	}

	// Write response
	w.Header().Set("Content-Type", "applicaiton/json")
	w.Write(json)

	return true, nil
}
