package main

import (
	"context"
	"encoding/json"
	"fmt"
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
	Host     string `json:"host"`
	Port     string `json:"port"`
	MongoURI string `json:"mongo_uri"`
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
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	// The Screenshots collection
	db = client.Database("ab-db").Collection("screenshots")

	// API setup
	router := mux.NewRouter()

	// Register Routes
	router.HandleFunc("/", Home)

	fmt.Printf("Server us running on %s%s\n", config.Host, config.Port)
	log.Fatal(http.ListenAndServe(config.Port, router))
}
