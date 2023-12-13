package main

import (
	"SS-Database/lib/handlers"
	ssPlaces "SS-Database/lib/places"
	"SS-Database/lib/users"
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
	app.HandlerIns.UserController = users.UserController{}
	app.HandlerIns.PlaceController = ssPlaces.PlaceController{}

	log.Print("starting server...")
	http.HandleFunc("/AddUser", app.HandlerIns.HandleRequest)
	http.HandleFunc("/UpdateUser", app.HandlerIns.HandleRequest)
	http.HandleFunc("/AddPlace", app.HandlerIns.HandleRequest)
	http.HandleFunc("/GetPlaceInfo", app.HandlerIns.HandleRequest)
	//http.HandleFunc("/AddReview", app.HandlerIns.HandleRequest)
	http.HandleFunc("/GetAllUsers", app.HandlerIns.HandleRequest)
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
