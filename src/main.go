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
		data := new([]map[string]interface{})
		jsonFile, _ := os.ReadFile("data.json")
		json.Unmarshal(jsonFile, data)

		placeController.AddPlaceBatch(context.Background(), app.HandlerIns.Client, *data)
		os.Exit(0)
	}
	app.HandlerIns.UserController = users.UserController{}
	app.HandlerIns.PlaceController = places.PlaceController{}

	log.Print("starting server...")
	http.HandleFunc("/AddUser", app.HandlerIns.HandleRequest)
	http.HandleFunc("/UpdateUser", app.HandlerIns.HandleRequest)
	http.HandleFunc("/RemoveUser", app.HandlerIns.HandleRequest)
	http.HandleFunc("/GetUserInfo", app.HandlerIns.HandleRequest)
	http.HandleFunc("/ReturnPassword", app.HandlerIns.HandleRequest)
	http.HandleFunc("/AddPlace", app.HandlerIns.HandleRequest)
	http.HandleFunc("/GetPlaceInfo", app.HandlerIns.HandleRequest)
	http.HandleFunc("/GetPlaceName", app.HandlerIns.HandleRequest)
	http.HandleFunc("/GetAllPlaces", app.HandlerIns.HandleRequest)
	http.HandleFunc("/GetAllUsers", app.HandlerIns.HandleRequest)
	http.HandleFunc("/AddFeedback", app.HandlerIns.HandleRequest)
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
