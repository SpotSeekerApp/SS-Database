package users

import (
	types "SS-Database/lib/types"
	"cloud.google.com/go/firestore"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (s UserController) AddUser(ctx context.Context, client *firestore.Client, data []byte) codes.Code {
	userInfo := new(types.UserRequest)
	err := json.Unmarshal(data, userInfo)
	fmt.Println(userInfo)

	_, err = client.Collection("Users").Doc(userInfo.Id).Create(ctx, map[string]interface{}{
		"name":    userInfo.Name,
		"surname": userInfo.Surname,
		"born":    userInfo.Born,
		"email":   userInfo.Email,
	})
	if err != nil {
		log.Fatalf("Failed adding users: %v", err)
		return codes.Aborted
	}

	return codes.OK
}

func (s UserController) AddFavoritePlace(ctx context.Context, client *firestore.Client, data []byte) codes.Code {
	favReqInfo := new(types.FavoritePlaceRequest)
	err := json.Unmarshal(data, &favReqInfo)
	fmt.Println(favReqInfo)

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
		}
		return tx.Set(ref, map[string]interface{}{
			"favorite_places": userdata.FavoritePlaces,
		}, firestore.MergeAll)

	})
	if status.Code(err) == codes.NotFound {
		return codes.NotFound
	}
	if err != nil {
		// Handle any errors appropriately in this section.
		log.Printf("An error has occurred: %s", err)
		return codes.Aborted
	}

	return codes.OK
}

func (s UserController) RemoveFavoritePlace(ctx context.Context, client *firestore.Client, data []byte) codes.Code {
	favReqInfo := new(types.FavoritePlaceRequest)
	err := json.Unmarshal(data, favReqInfo)
	fmt.Println(favReqInfo)

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
	if status.Code(err) == codes.NotFound {
		return codes.NotFound
	}
	if err != nil {
		// Handle any errors appropriately in this section.
		log.Printf("An error has occurred: %s", err)
		return codes.Aborted
	}

	return codes.OK
}
