package user

import (
	types "SS-Database/lib/types"
	"cloud.google.com/go/firestore"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"log"
)

type PlaceController struct {
	Id          string `firebase:"placeId"`
	Name        string `firebase:"placeName"`
	Location    string `firebase:"location"`
	Link2Photo  string `firebase:"link2Photo"`
	PhoneNumber string `firebase:"phoneNumber"`
	InitReview  string `firebase:"initReview"`
}

func (s PlaceController) AddPlace(ctx context.Context, client *firestore.Client, data []byte) codes.Code {
	placeInfo := new(types.PlaceRequest)
	err := json.Unmarshal(data, placeInfo)
	fmt.Println(placeInfo)

	_, err = client.Doc("Places/"+placeInfo.Id).Create(ctx, map[string]interface{}{
		"place_name":   placeInfo.Name,
		"location":     placeInfo.Location,
		"photo_link":   placeInfo.Link2Photo,
		"phone_number": placeInfo.PhoneNumber,
	})
	if err != nil {
		log.Fatalf("Failed adding users: %v", err)
		return codes.Aborted
	}

	return codes.OK
}

func (s PlaceController) AddReview(ctx context.Context, client *firestore.Client, data []byte) codes.Code {
	placeInfo := new(types.ReviewRequest)
	err := json.Unmarshal(data, placeInfo)
	fmt.Println(placeInfo)

	ref := client.Doc("Places/" + placeInfo.PlaceId)
	if ref == nil {
		return codes.NotFound
	}
	ref = client.Doc("Places/" + placeInfo.PlaceId + "/Reviews/" + placeInfo.ReviewId)
	err = client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		return tx.Create(ref, map[string]interface{}{
			"comment": placeInfo.Comment,
			"rating":  placeInfo.Rating,
		})
	})
	if err != nil {
		// Handle any errors appropriately in this section.
		log.Printf("An error has occurred: %s", err)
		return codes.Aborted
	}
	return codes.OK
}
