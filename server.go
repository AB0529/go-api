package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"net/http"

	"github.com/gorilla/mux"
	s "github.com/gosimple/slug"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// // Home the home route of the API
// func Home(w http.ResponseWriter, r *http.Request) {
// 	http.ServeFile(w, r, "./api-docs/index.html")
// }

// ImageHandler shows the image from the database
func ImageHandler(w http.ResponseWriter, r *http.Request) {
	filename := mux.Vars(r)["name"]
	ext := filepath.Ext(filename)

	// Only allow jpg, png, and gif
	if ext != ".png" && ext != ".jpg" && ext != ".gif" {
		fmt.Fprint(w, "404 not found")
		return
	}

	// Search DB for filename
	result, _, err := FindScreenshot(strings.TrimSuffix(filename, ext))

	// Handle not found
	if err != nil {
		LogE.Println("(FindScreenshot) error: ", err)
		fmt.Fprint(w, "404 not found")
		return
	}
	var screenshot struct {
		Image []byte
	}
	bsonBytes, err := bson.Marshal(result)
	bson.Unmarshal(bsonBytes, &screenshot)

	if err != nil {
		LogE.Println("(bsonMarshal) error: ", err)
		return
	}

	// Set headers
	w.Header().Set("Content-Type", result["mime"].(string))

	// If found, serve file
	http.ServeContent(w, r, filename, time.Now(), bytes.NewReader(screenshot.Image))
}

// GetAllScreenshots will return every screenshot in the database
func GetAllScreenshots(w http.ResponseWriter, r *http.Request) {
	// Key auth
	key := s.Make(mux.Vars(r)["key"])

	if key != config.APIKey {
		SendJSON(w, Response{
			Status: 401,
			State:  "fail",
			Result: "error: unauthorized key",
		})
		return
	}

	// Get all screenshots
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// Find filter in db
	cur, _ := db.Find(ctx, bson.M{})

	// Item could not be found
	if cur.RemainingBatchLength() <= 0 {
		LogW.Println("No items found in database")
		return
	}

	var items []bson.M

	for cur.Next(ctx) {
		var res bson.M
		cur.Decode(&res)

		if res != nil {
			items = append(items, res)
		}
	}

	if err := cur.Err(); err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	SendJSON(w, Response{
		Status: 200,
		State:  "ok",
		Result: items,
	})
}

// GetScreenshot will return info on the screenshot
func GetScreenshot(w http.ResponseWriter, r *http.Request) {
	id := s.Make(mux.Vars(r)["id"])
	screenshot, _, err := FindScreenshot(id)

	// If not found, return 404
	if err != nil {
		SendJSON(w, Response{
			Status: 404,
			State:  "fail",
			Result: "error: could not find screenshot",
		})
		return
	}

	// If found, send the screenshot data
	SendJSON(w, Response{
		Status: 200,
		State:  "ok",
		Result: screenshot,
	})

}

// CreaeteScreenshot will create a screenshot and save it to the database
func CreaeteScreenshot(w http.ResponseWriter, r *http.Request) {
	// Validate key
	r.ParseMultipartForm(0)
	key := r.FormValue("key")

	if key != config.APIKey {
		SendJSON(w, Response{
			Status: 401,
			State:  "fail",
			Result: "error: unauthorized key",
		})
		return
	}

	// Get filename and the screenshot itself
	name := r.MultipartForm.Value["name"][0]
	file, header, err := r.FormFile("screenshot")
	defer file.Close()
	screenshot, _ := ioutil.ReadAll(file)
	if err != nil {
		SendJSON(w, Response{
			Status: 400,
			State:  "fail",
			Result: "error: no image sent",
		})
		return
	}

	contentType := http.DetectContentType(screenshot)
	ext := ""

	switch contentType {
	case "image/png":
		ext = ".png"
		break
	case "image/jpg":
		ext = ".jpg"
		break
	case "image/gif":
		ext = ".gif"
		break
	}

	LogI.Println(fmt.Sprintf("`%s` : %.2fKB", name+ext, float32(header.Size/1000)))

	// Save file data to Mongo
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err = db.InsertOne(ctx, bson.M{
		"image":     screenshot,
		"name":      name,
		"mime":      contentType,
		"timestamp": time.Now(),
	})

	if err != nil {
		LogE.Println("error: create screeshot: ", err)
	}

	fmt.Fprint(w, config.RedirectURL+name+ext)
}

// DeleteScreenshot will delete a screenshot from the database
func DeleteScreenshot(w http.ResponseWriter, r *http.Request) {
	// Validate key
	key := s.Make(mux.Vars(r)["key"])

	if key != config.APIKey {
		SendJSON(w, Response{
			Status: 401,
			State:  "fail",
			Result: "error: unauthorized key",
		})
		return
	}

	id := s.Make(mux.Vars(r)["id"])
	res, _, _ := FindScreenshot(id)

	// No items found
	if res == nil {
		SendJSON(w, Response{
			Status: 404,
			State:  "fail",
			Result: "error: could not find screenshot",
		})
		return
	}

	// Delete screenshot
	if res["name"] == id {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if _, err := db.DeleteOne(ctx, bson.M{"name": id}); err != nil {
			LogE.Println(err)
		}

		// Send deleted data
		SendJSON(w, Response{
			Status: 200,
			State:  "ok",
			Result: "Deleted screenshot " + id,
		})
	}
}

// SendJSON will send a JSON response to the API according the Response type
func SendJSON(w http.ResponseWriter, resp Response) (bool, error) {
	// Decode JSON
	json, err := json.Marshal(resp)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false, nil
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)

	return true, nil
}

// FindScreenshot will find a screenshot in the database given a filename
func FindScreenshot(filename string) (primitive.M, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// Find filter in db
	cur, _ := db.Find(ctx, bson.M{"name": filename})

	// Item could not be found
	if cur.RemainingBatchLength() <= 0 {
		return nil, cancel, errors.New("no items found")
	}

	var item bson.M

	for cur.Next(ctx) {
		var res bson.M
		cur.Decode(&res)

		if res != nil {
			item = res
			break
		}
	}

	if err := cur.Err(); err != nil {
		panic(err)
	}

	return item, cancel, nil
}
