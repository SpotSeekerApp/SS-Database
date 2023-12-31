package main

import (
	"SS-Database/lib/handlers"
	places "SS-Database/lib/places"
	"SS-Database/lib/users"
	"encoding/json"
	firebase "firebase.google.com/go"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"
)

type DataBaseAPI struct {
	HandlerIns handlers.HandlerInstance
}

func main() {
	app := DataBaseAPI{}
	err := app.createClient()
	if err != nil {
		log.Fatalf("Failed adding aturing: %v", err)
	}
	isDataPush := os.Getenv("DATAPUSH")
	if isDataPush == "true" {
		placeController := places.NewPlaceController()
		data := make([]map[string]interface{}, 0)
		dataFiles, _ := os.ReadDir("SS-Vision-and-NLP/data")
		for _, file := range dataFiles {
			fileData := make([]map[string]interface{}, 0)
			log.Printf("Reading %s", file.Name())
			jsonFile, _ := os.ReadFile("SS-Vision-and-NLP/data/" + file.Name())
			json.Unmarshal(jsonFile, &fileData)
			data = append(data, fileData...)
		}
		tagData := make([]map[string]interface{}, 0)
		tagJsonFile, _ := os.ReadFile("SS-Vision-and-NLP/outputs/output_tags.json")
		tags := make(map[string]map[string]float64, 0)
		json.Unmarshal(tagJsonFile, &tagData)
		for _, tagMap := range tagData {
			placeTags := make(map[string]float64)

			for _, tagVals := range tagMap["tags"].([]interface{}) {
				placeTags[tagVals.([]interface{})[0].(string)] = tagVals.([]interface{})[1].(float64)
			}
			tags[tagMap["place_id"].(string)] = placeTags
		}

		placeController.AddPlaceBatch(context.Background(), app.HandlerIns.Client, data, tags)
		os.Exit(0)
	}
	app.HandlerIns.UserController = users.UserController{}
	app.HandlerIns.PlaceController = places.NewPlaceController()

	log.Print("starting server...")
	http.HandleFunc("/FilterPlaces", app.HandlerIns.HandleRequest)
	http.HandleFunc("/AddUser", app.HandlerIns.HandleRequest)
	http.HandleFunc("/UpdateUser", app.HandlerIns.HandleRequest)
	http.HandleFunc("/RemoveUser", app.HandlerIns.HandleRequest)
	http.HandleFunc("/GetUserInfo", app.HandlerIns.HandleRequest)
	http.HandleFunc("/AddPlace", app.HandlerIns.HandleRequest)
	http.HandleFunc("/RemovePlace", app.HandlerIns.HandleRequest)
	http.HandleFunc("/UpdatePlace", app.HandlerIns.HandleRequest)
	http.HandleFunc("/GetPlaceInfo", app.HandlerIns.HandleRequest)
	http.HandleFunc("/GetPlaceName", app.HandlerIns.HandleRequest)
	http.HandleFunc("/GetAllTags", app.HandlerIns.HandleRequest)
	http.HandleFunc("/GetTagByPlace", app.HandlerIns.HandleRequest)
	http.HandleFunc("/GetAllPlaces", app.HandlerIns.HandleRequest)
	http.HandleFunc("/GetAllUsers", app.HandlerIns.HandleRequest)
	http.HandleFunc("/AddFeedback", app.HandlerIns.HandleRequest)
	http.HandleFunc("/GetFeedbacks", app.HandlerIns.HandleRequest)
	http.HandleFunc("/AddFavoritePlace", app.HandlerIns.HandleRequest)
	http.HandleFunc("/RemoveFavoritePlace", app.HandlerIns.HandleRequest)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

	defer app.HandlerIns.Client.Close()
}

func (d *DataBaseAPI) createClient() error {
	opt := option.WithCredentialsFile("./credential-token.json")
	projectID, _ := os.LookupEnv("ProjectID")
	config := &firebase.Config{ProjectID: projectID}
	ctx := context.Background()

	app, _ := firebase.NewApp(ctx, config, opt)
	client, err := app.Firestore(ctx)

	d.HandlerIns.Client = client
	return err
}
