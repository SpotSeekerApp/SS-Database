package places

import (
	types "SS-Database/lib/types"
	"cloud.google.com/go/firestore"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"os"
)

type PlaceController struct {
	PlaceId    string `firebase:"placeId"`
	Name       string `firebase:"placeName"`
	Link       string `firebase:"location"`
	Link2Photo string `firebase:"link2Photo"`
	InitReview string `firebase:"initReview"`
	Category   string
	Rating     float32
}

func NewPlaceController() *PlaceController {
	return &PlaceController{}
}

func (s PlaceController) AddPlace(ctx context.Context, client *firestore.Client, data []byte) codes.Code {
	placeInfo := new(types.PlaceRequest)
	placeInfo.OwnerId = -1 // default
	err := json.Unmarshal(data, placeInfo)
	fmt.Println(placeInfo)

	in := map[string]interface{}{
		"placeId":      placeInfo.PlaceId,
		"placeName":    placeInfo.PlaceName,
		"mainCategory": placeInfo.MainCategory,
		"link":         placeInfo.Link,
		"tags":         placeInfo.Tags,
	}
	if placeInfo.OwnerId != -1 {
		in["ownerId"] = placeInfo.OwnerId
	}

	ref := client.Collection("Places").NewDoc()
	in["placeId"] = ref.ID
	_, err = ref.Create(ctx, in)
	if err != nil {
		log.Fatalf("Failed adding users: %v", err)
		return codes.Aborted
	}

	return codes.OK
}

func (s PlaceController) checkIfExists(ctx context.Context, client *firestore.Client, id string) bool {

	iter := client.Collection("Places").Select("place_id").Documents(ctx)
	q, _ := iter.GetAll()
	for _, s := range q {
		if s.Data()["place_id"] == id {
			return true
		}
	}
	return false
}

func (s PlaceController) AddPlaceBatch(ctx context.Context, client *firestore.Client, data []map[string]interface{}) codes.Code {

	for _, place := range data {
		placeInfo := types.PlaceRequest{}
		in := map[string]interface{}{
			"placeId":      place["place_id"],
			"placeName":    place["name"],
			"mainCategory": place["main_category"],
			"link":         place["link"],
		}
		_ = mapstructure.Decode(in, &placeInfo)
		ref := client.Collection("Places")
		err := client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {

			docRef, _ := tx.DocumentRefs(ref).GetAll()
			docSnaps, _ := tx.GetAll(docRef)
			for _, s := range docSnaps {
				if s.Data()["placeId"] == placeInfo.PlaceId {
					return os.ErrExist
				}
			}
			err := tx.Create(ref.Doc(placeInfo.PlaceId), in)
			return err
		})
		if status.Code(err) == codes.NotFound {
			return codes.NotFound
		}
	}

	return codes.OK
}

func (s PlaceController) GetPlaceName(ctx context.Context, client *firestore.Client, data []byte) ([]byte, codes.Code) {
	placeInfo := new(types.PlaceRequest)
	err := json.Unmarshal(data, placeInfo)
	fmt.Println(placeInfo)

	q := client.Collection("Places").Where("placeId", "==", placeInfo.PlaceId).Select("placeName")
	ref, err := q.Documents(ctx).GetAll()
	name, _ := ref[0].DataAtPath(firestore.FieldPath{"placeName"})
	if err != nil {
		//log.Fatalf("Failed adding users: %v", err)
		return []byte{}, codes.Aborted
	}
	jsonStr, _ := json.Marshal(name)
	return jsonStr, codes.OK
}

func (s PlaceController) GetPlaceInfo(ctx context.Context, client *firestore.Client, data []byte) ([]byte, codes.Code) {
	placeInfo := new(types.PlaceRequest)
	err := json.Unmarshal(data, placeInfo)
	fmt.Println(placeInfo)

	q := client.Collection("Places").Where("placeId", "==", placeInfo.PlaceId)
	ref, err := q.Documents(ctx).GetAll()
	_ = ref[0].DataTo(placeInfo)
	if err != nil {
		//log.Fatalf("Failed adding users: %v", err)
		return []byte{}, codes.Aborted
	}
	jsonStr, _ := json.Marshal(placeInfo)
	return jsonStr, codes.OK
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

func (s PlaceController) GetAllPlaces(ctx context.Context, client *firestore.Client) ([]byte, codes.Code) {
	docRefs, err := client.Collection("Places").Documents(ctx).GetAll()
	if err != nil {
		return []byte{}, codes.NotFound
	}

	if len(docRefs) == 0 {
		return []byte{}, codes.OK
	}

	resp := map[string]types.PlaceRequest{}
	for _, docRef := range docRefs {
		var placeData types.PlaceRequest
		_ = docRef.DataTo(&placeData)
		resp[placeData.PlaceId] = placeData
	}

	jsonStr, err := json.Marshal(resp)
	return jsonStr, codes.OK
}
