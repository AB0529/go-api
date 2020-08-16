package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Config the configuration structure
type Config struct {
	Host        string `json:"host"`
	Port        string `json:"port"`
	APIKey      string `json:"api_key"`
	MongoURI    string `json:"mongo_uri"`
	RedirectURL string `json:"redirect_url"`
}

// Response the response structure the API will return
type Response struct {
	Status int         `json:"status"`
	State  string      `json:"state"`
	Result interface{} `json:"result"`
}

var config Config
var ctx context.Context
var cancel context.CancelFunc
var db *mongo.Collection

func main() {
	// Load config
	file, _ := ioutil.ReadFile("./config.json")
	json.Unmarshal(file, &config)

	// Mongo Setup
	client, err := mongo.NewClient(options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		panic(err)
	}
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	// The Screenshots collection
	db = client.Database("ab-db").Collection("screenshots")
	LogI.Println("Mongo connection... OK")

	// API setup
	router := mux.NewRouter()

	// Register Routes
	router.HandleFunc("/", Home)
	router.HandleFunc("/img/{name}", ImageHandler)
	router.HandleFunc("/screenshots/{id}", GetScreenshot).Methods("GET")
	router.HandleFunc("/screenshots/{key}/{id}", DeleteScreenshot).Methods("DELETE")
	router.HandleFunc("/screenshots", CreaeteScreenshot).Methods("POST")

	LogI.Printf("Server running on %s%s\n", config.Host, config.Port)
	log.Fatal(http.ListenAndServe(config.Port, router))
}
