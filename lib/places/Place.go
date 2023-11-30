package user

import (
	types "SS-Database/lib/types"
	"cloud.google.com/go/firestore"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"log"
)

type UserController struct {
	Id             string            `firebase:"userId"`
	Name           string            `firebase:"name"`
	Surname        string            `firebase:"surname"`
	Born           string            `firebase:"born"`
	Email          string            `firebase:"email"`
	FavoritePlaces map[string]string `firebase:"favorite_places"`
}

type Feedback struct {
	Name    string            `json:"name"`
	PlaceId string            `json:"place_id"`
	Fields  map[string]string `json:"fields"`
}

func (s UserController) AddUser(ctx context.Context, client *firestore.Client, data []byte) error {
	userInfo := new(types.UserRequest)
	err := json.Unmarshal(data, &userInfo)
	fmt.Println(*userInfo)

	_, err = client.Collection("Users").Doc(userInfo.Id).Create(ctx, map[string]interface{}{
		"name":    userInfo.Name,
		"surname": userInfo.Surname,
		"born":    userInfo.Born,
		"email":   userInfo.Email,
	})
	if err != nil {
		log.Fatalf("Failed adding user: %v", err)
	}

	return err
}

func (s UserController) AddFavoritePlace(ctx context.Context, client *firestore.Client, data []byte) error {
	favReqInfo := new(types.FavoritePlaceRequest)
	err := json.Unmarshal(data, &favReqInfo)
	fmt.Println(*favReqInfo)

	ref := client.Collection("Users").Doc(favReqInfo.UserId)
	err = client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {

		doc, err := tx.Get(ref)
		if err != nil {
			return err
		}
		var userdata UserController
		_ = doc.DataTo(&userdata)

		fmt.Println(userdata.FavoritePlaces)
		if userdata.FavoritePlaces == nil {
			userdata.FavoritePlaces = map[string]string{favReqInfo.PlaceId: favReqInfo.PlaceName}
		} else {
			userdata.FavoritePlaces[favReqInfo.PlaceId] = favReqInfo.PlaceName
			if err != nil {
				return err
			}
		}
		return tx.Set(ref, map[string]interface{}{
			"favorite_places": userdata.FavoritePlaces,
		}, firestore.MergeAll)

	})
	if err != nil {
		// Handle any errors appropriately in this section.
		log.Printf("An error has occurred: %s", err)
	}

	return err
}

func (s UserController) RemoveFavoritePlace(ctx context.Context, client *firestore.Client, data []byte) error {
	favReqInfo := new(types.FavoritePlaceRequest)
	err := json.Unmarshal(data, &favReqInfo)
	fmt.Println(*favReqInfo)

	ref := client.Collection("Users").Doc(favReqInfo.UserId)
	err = client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		doc, err := tx.Get(ref) // tx.Get, NOT ref.Get!
		if err != nil {
			return err
		}
		var userdata UserController
		_ = doc.DataTo(&userdata)
		fmt.Println(userdata.FavoritePlaces)

		delete(userdata.FavoritePlaces, favReqInfo.PlaceId)
		fmt.Println(userdata.FavoritePlaces)
		return tx.Set(ref, map[string]interface{}{
			"favorite_places": userdata.FavoritePlaces,
		}, firestore.Merge(firestore.FieldPath{"favorite_places"}))

	})
	if err != nil {
		// Handle any errors appropriately in this section.
		log.Printf("An error has occurred: %s", err)
	}

	return err
}

//func (s User) AddReview(ctx context.Context, client *firestore.Client, data []byte) error {
//	var err error
//
//	feedback := Feedback{}
//	err = json.Unmarshal(data, &feedback)
//	fmt.Println(feedback)
//
//	ref := client.Doc("Users/" + feedback.Name + "/Feedbacks/" + feedback.PlaceId)
//	err = client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
//		return tx.Create(ref, feedback.Fields)
//	})
//	if err != nil {
//		// Handle any errors appropriately in this section.
//		log.Printf("An error has occurred: %s", err)
//	}
//
//	return err
//}
//
