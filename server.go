package main

import (
	"fmt"
	"net/http"
)

// Home the home route of the API
func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Home")
}