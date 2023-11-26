package main

import (
	"SS-Database/lib/handlers"
	"SS-Database/lib/user"
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
	app := DataBaseAPI{HandlerIns: handlers.HandlerInstance{}}
	err := app.createClient()
	if err != nil {
		log.Fatalf("Failed adding aturing: %v", err)
	}
	app.HandlerIns.UserController = user.User{}

	//data, _ := os.ReadFile("./example-data.json")
	//
	//err = app.HandlerIns.UserController.AddUser(context.Background(), app.HandlerIns.Client, data)
	//if err != nil {
	//	log.Fatalf("Failed adding aturing: %v", err)
	//}
	//err := userIns.AddFavoritePlace(context.Background(), client, data)
	//if err != nil {
	//	log.Fatalf("Failed adding aturing: %v", err)
	//}

	//err := userIns.RemoveFavoritePlace(context.Background(), client, data)
	//if err != nil {
	//	log.Fatalf("Failed adding aturing: %v", err)
	//}
	//err := userIns.AddReview(context.Background(), client, data)
	//if err != nil {
	//	log.Fatalf("Failed adding aturing: %v", err)
	//}
	log.Print("starting server...")
	http.HandleFunc("/AddUser", app.HandlerIns.UserHandler)

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
