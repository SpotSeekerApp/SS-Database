package data

//
//import (
//	"golang.org/x/net/context"
//	"google.golang.org/api/option"
//	"log"
//)
//
//func main() {
//
//	opt := option.WithCredentialsFile("./credential-token.json")
//	app, err := initializeApp(opt)
//	if err != nil {
//		log.Fatalf("error initializing app: %v\n", err)
//	}
//	ctx := context.Background()
//	client, err := app.Firestore(ctx)
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	defer client.Close()
//}
