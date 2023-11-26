package data

import (
	firebase "firebase.google.com/go"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"log"
	"os"
)

func main() {

	opt := option.WithCredentialsFile("./credential-token.json")
	app, err := initializeApp(opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	ctx := context.Background()
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	defer client.Close()
}

func initializeApp(opt option.ClientOption) (*firebase.App, error) {
	projectID, _ := os.LookupEnv("ProjectID")
	config := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(context.Background(), config, opt)
	return app, err
}
