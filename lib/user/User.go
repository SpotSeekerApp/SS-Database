package user

import (
	"cloud.google.com/go/firestore"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"log"
)

type User struct {
	Name            string                 `firestore:"name"`
	Surname         string                 `firestore:"surname"`
	Born            string                 `firestore:"born"`
	Favorite_places map[string]interface{} `firestore:"favorite_places"`
}
type FavoritePlace struct {
	Name      string `json:"name"`
	PlaceId   string `json:"place_id"`
	PlaceName string `json:"place_name"`
}

type Feedback struct {
	Name    string            `json:"name"`
	PlaceId string            `json:"place_id"`
	Fields  map[string]string `json:"fields"`
}

func (s User) AddUser(ctx context.Context, client *firestore.Client, data []byte) error {
	var err error

	m := make(map[string]interface{})
	err = json.Unmarshal(data, &m)
	fmt.Println(m)
	val := m["name"]
	id := fmt.Sprintf("%v", val)

	_, err = client.Collection("Users").Doc(id).Create(ctx, m["fields"])
	if err != nil {
		log.Fatalf("Failed adding aturing: %v", err)
	}

	return err
}

func (s User) AddFavoritePlace(ctx context.Context, client *firestore.Client, data []byte) error {
	var err error

	favPlace := FavoritePlace{}
	err = json.Unmarshal(data, &favPlace)
	fmt.Println(favPlace)

	ref := client.Collection("Users").Doc(favPlace.Name)
	err = client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {

		doc, err := tx.Get(ref) // tx.Get, NOT ref.Get!
		if err != nil {
			return err
		}
		var userdata User
		_ = doc.DataTo(&userdata)
		fmt.Println(userdata.Favorite_places)
		if userdata.Favorite_places == nil {
			return tx.Set(ref, map[string]interface{}{
				"favorite_places": map[string]string{favPlace.PlaceId: favPlace.PlaceName},
			}, firestore.MergeAll)
		} else {
			userdata.Favorite_places[favPlace.PlaceId] = favPlace.PlaceName
			if err != nil {
				return err
			}
			return tx.Set(ref, map[string]interface{}{
				"favorite_places": userdata.Favorite_places,
			}, firestore.MergeAll)
		}

	})
	if err != nil {
		// Handle any errors appropriately in this section.
		log.Printf("An error has occurred: %s", err)
	}

	return err
}

func (s User) RemoveFavoritePlace(ctx context.Context, client *firestore.Client, data []byte) error {
	var err error

	favPlace := FavoritePlace{}
	err = json.Unmarshal(data, &favPlace)
	fmt.Println(favPlace)

	ref := client.Collection("Users").Doc(favPlace.Name)
	err = client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {

		doc, err := tx.Get(ref) // tx.Get, NOT ref.Get!
		if err != nil {
			return err
		}
		var userdata User
		_ = doc.DataTo(&userdata)
		fmt.Println(userdata.Favorite_places)

		delete(userdata.Favorite_places, favPlace.PlaceId)
		fmt.Println(userdata.Favorite_places)
		return tx.Set(ref, map[string]interface{}{
			"favorite_places": userdata.Favorite_places,
		}, firestore.Merge(firestore.FieldPath{"favorite_places"}))

	})
	if err != nil {
		// Handle any errors appropriately in this section.
		log.Printf("An error has occurred: %s", err)
	}

	return err
}

func (s User) AddReview(ctx context.Context, client *firestore.Client, data []byte) error {
	var err error

	feedback := Feedback{}
	err = json.Unmarshal(data, &feedback)
	fmt.Println(feedback)

	ref := client.Doc("Users/" + feedback.Name + "/Feedbacks/" + feedback.PlaceId)
	err = client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		return tx.Create(ref, feedback.Fields)
	})
	if err != nil {
		// Handle any errors appropriately in this section.
		log.Printf("An error has occurred: %s", err)
	}

	return err
}
