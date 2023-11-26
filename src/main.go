package main

import (
	"SS-Database/lib/user"
	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"log"
	"os"
)

func main() {
	userIns := user.User{}
	data, _ := os.ReadFile("./example-data.json")
	client, _ := initializeApp()

	err := userIns.AddUser(context.Background(), client, data)
	if err != nil {
		log.Fatalf("Failed adding aturing: %v", err)
	}

	defer client.Close()
}

func initializeApp() (*firestore.Client, error) {
	opt := option.WithCredentialsFile("./credential-token.json")
	projectID, _ := os.LookupEnv("ProjectID")
	config := &firebase.Config{ProjectID: projectID}
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, config, opt)
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	return client, err
}
