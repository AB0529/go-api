package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Config the configuration structure
type Config struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	MongoURI string `json:"mongo_uri"`
}

// Response the response structure the API will return
type Resposne struct {
	Status int         `json:"status"`
	State  string      `json:"state"`
	Result interface{} `json:"result"`
}

var config Config

func main() {
	// Load config
	file, _ := ioutil.ReadFile("./config.json")
	json.Unmarshal(file, &config)

	// API setup
	router := mux.NewRouter()

	// Register Routes
	router.HandleFunc("/", Home)

	fmt.Printf("Server us running on %s%s\n", config.Host, config.Port)
	log.Fatal(http.ListenAndServe(config.Port, router))
}
