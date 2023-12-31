package places

import (
	types "SS-Database/lib/types"
	"SS-Database/utils"
	"cloud.google.com/go/firestore"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

const REVIEW_COUNT = 10

type PlaceController struct {
	PlaceId    string `firebase:"placeId"`
	PlaceName  string `firebase:"placeName"`
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
	placeInfo.UserId = "-1" // default
	err := json.Unmarshal(data, placeInfo)
	fmt.Println(placeInfo)

	in := map[string]interface{}{
		"placeId":      placeInfo.PlaceId,
		"placeName":    placeInfo.PlaceName,
		"mainCategory": placeInfo.MainCategory,
		"link":         placeInfo.Link,
	}
	if placeInfo.UserId != "-1" {
		in["userId"] = placeInfo.UserId
	}

	ref := client.Collection("Places").NewDoc()
	in["placeId"] = ref.ID
	_, err = ref.Create(ctx, in)
	if err != nil {
		log.Fatalf("Failed adding users: %v", err)
		return codes.Aborted
	}
	if placeInfo.Tags == nil {
		return codes.OK
	}
	_, err = client.Collection("Places/"+ref.ID+"/Tags").Doc("TagList").Create(ctx, map[string]interface{}{
		"tags":    placeInfo.Tags,
		"placeId": ref.ID,
	})
	if err != nil {
		log.Fatalf("Failed adding users: %v", err)
		return codes.Aborted
	}

	return codes.OK
}

func (s PlaceController) RemovePlace(ctx context.Context, client *firestore.Client, data []byte) codes.Code {
	placeInfo := new(types.PlaceRequest)
	fmt.Println(placeInfo)
	_ = json.Unmarshal(data, placeInfo)
	ref := client.Collection("Places").Doc(placeInfo.PlaceId)
	placeTags := ref.Collection("Tags")
	bulkWriter := client.BulkWriter(ctx)

	for {
		// Get a batch of documents
		iter := placeTags.Limit(utils.BATCH_SIZE).Documents(ctx)
		numDeleted := 0

		// Iterate through the documents, adding
		// a delete operation for each one to the BulkWriter.
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				break
			}
			bulkWriter.Delete(doc.Ref)
			numDeleted++
		}
		if numDeleted == 0 {
			bulkWriter.End()
			break
		}
		bulkWriter.Flush()
	}
	bulkWriter = client.BulkWriter(ctx)
	bulkWriter.Delete(ref)
	bulkWriter.End()
	bulkWriter.Flush()
	return codes.OK
}

func (s PlaceController) UpdateTags(placeId string, client *firestore.Client, tags map[string]float64, change bool) codes.Code {

	ref := client.Collection("Places/" + placeId + "/Tags").Doc("TagList")
	err := client.RunTransaction(context.Background(), func(ctx context.Context, tx *firestore.Transaction) error {

		doc, err := tx.Get(ref)
		if err != nil {
			return err
		}
		var tagData types.PlaceRequest
		_ = doc.DataTo(&tagData)

		fmt.Println(tagData.Tags)

		for tag, updateVal := range tags {
			if change == true {
				tagData.Tags[tag] = updateVal
			} else {
				tagData.Tags[tag] += updateVal
			}
		}

		return tx.Set(ref, map[string]interface{}{
			"tags": tagData.Tags,
		}, firestore.Merge(firestore.FieldPath{"tags"}))

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

func (s PlaceController) UpdatePlace(ctx context.Context, client *firestore.Client, data []byte) codes.Code {
	placeInfo := new(types.PlaceRequest)
	err := json.Unmarshal(data, placeInfo)
	if err != nil {
		log.Printf("Failed removing user: %v", err)
		return codes.Aborted
	}
	tagsUpdated := false
	if placeInfo.Tags != nil {
		_ = s.UpdateTags(placeInfo.PlaceId, client, placeInfo.Tags, true)
		placeInfo.Tags = nil
		tagsUpdated = true
	}

	docSnap, err := client.Collection("Places").Doc(placeInfo.PlaceId).Get(ctx)
	if err != nil {
		log.Printf("Failed removing user: User not found")
		return codes.NotFound
	}

	_, err = docSnap.Ref.Update(ctx, utils.ExtractNonEmptyFields(*placeInfo))

	if err != nil {
		if tagsUpdated == false {
			log.Printf("Failed removing user: %v", err)
			return codes.NotFound
		}
	}
	return codes.OK
}

func (s PlaceController) GetTagByPlace(ctx context.Context, client *firestore.Client, data []byte) ([]byte, codes.Code) {
	placeInfo := new(types.PlaceRequest)
	_ = json.Unmarshal(data, placeInfo)
	fmt.Println(placeInfo)

	docSnap, _ := client.Doc("Places/" + placeInfo.PlaceId + "/Tags/TagList").Get(ctx)
	_ = docSnap.DataTo(placeInfo)
	jsonStr, _ := json.Marshal(placeInfo.Tags)
	return jsonStr, codes.OK
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

func (s PlaceController) AddPlaceBatch(ctx context.Context, client *firestore.Client, data []map[string]interface{}, tags map[string]map[string]float64) codes.Code {

	for _, place := range data {
		var placeImages []string
		for _, imgUrl := range place["images"].([]interface{}) {
			placeImages = append(placeImages, imgUrl.(map[string]interface{})["link"].(string))
		}
		placeInfo := types.PlaceRequest{}
		in := map[string]interface{}{
			"placeId":      place["place_id"],
			"placeName":    place["name"],
			"mainCategory": place["main_category"],
			"link":         place["link"],
			"images":       placeImages,
		}
		_ = mapstructure.Decode(in, &placeInfo)
		ref := client.Collection("Places")

		valMap, _ := place["detailed_reviews"].([]interface{})
		if len(valMap) > 0 {
			in["firstReview"] = valMap[0].(map[string]interface{})["review_text"].(string)
		} else {
			in["firstReview"] = "No reviews"
		}

		docSnaps, _ := ref.Documents(ctx).GetAll()
		for _, s := range docSnaps {
			if s.Data()["placeId"] == placeInfo.PlaceId {
				log.Printf("Place with id %s already exists", placeInfo.PlaceId)
			}
		}
		_, err := ref.Doc(placeInfo.PlaceId).Create(ctx, in)

		log.Printf("Creating new place with id %s", placeInfo.PlaceId)

		reviewRef := ref.Doc(placeInfo.PlaceId).Collection("Reviews")

		err = s.AddReviews(reviewRef, valMap)

		_, err = ref.Doc(placeInfo.PlaceId).Collection("Tags").Doc("TagList").Create(ctx, map[string]interface{}{
			"tags":    tags[placeInfo.PlaceId],
			"placeId": placeInfo.PlaceId,
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
	nameRes, _ := name.(string)
	placeInfo.PlaceName = nameRes
	jsonStr, _ := json.Marshal(placeInfo)
	return jsonStr, codes.OK
}

func (s PlaceController) GetPlaceInfo(ctx context.Context, client *firestore.Client, data []byte, field string) ([]byte, codes.Code) {
	placeInfo := new(types.PlaceRequest)
	err := json.Unmarshal(data, placeInfo)
	fmt.Println(placeInfo)

	var q firestore.Query
	if field == "place_id" {
		q = client.Collection("Places").Where("placeId", "==", placeInfo.PlaceId)
	} else {
		q = client.Collection("Places").Where("userId", "==", placeInfo.UserId)
	}
	ref, err := q.Documents(ctx).GetAll()
	_ = ref[0].DataTo(placeInfo)
	if err != nil {
		//log.Fatalf("Failed adding users: %v", err)
		return []byte{}, codes.Aborted
	}
	jsonStr, _ := json.Marshal(placeInfo)
	return jsonStr, codes.OK
}

func (s PlaceController) AddReviews(ref *firestore.CollectionRef, data []interface{}) error {
	for idx, place := range data {
		if idx >= REVIEW_COUNT {
			break
		}
		place, _ := place.(map[string]interface{})
		reviewInfo := types.ReviewRequest{}
		in := map[string]interface{}{
			"reviewId":     place["review_id"],
			"reviewerName": place["reviewer_name"],
			"rating":       place["rating"],
			"comment":      place["review_text"],
			"date":         place["published_at_date"],
		}
		_ = mapstructure.Decode(in, &reviewInfo)
		_, _ = ref.Doc(reviewInfo.ReviewId).Create(context.Background(), in)
	}
	return nil
}

func (s PlaceController) GetAllPlaces(ctx context.Context, client *firestore.Client, data []byte) ([]byte, codes.Code) {
	placeInfo := new(types.PlaceRequest)
	placeInfo.UserId = "-1"
	err := json.Unmarshal(data, placeInfo)
	fmt.Println(placeInfo)
	var docRefs []*firestore.DocumentSnapshot
	if placeInfo.UserId == "-1" {
		docRefs, err = client.Collection("Places").Documents(ctx).GetAll()
	} else {
		docRefs, err = client.Collection("Places").Where("userId", "==", placeInfo.UserId).Documents(ctx).GetAll()
	}
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

func (s PlaceController) GetAllTagVals(ctx context.Context, client *firestore.Client) ([]byte, codes.Code) {
	docRefs, _ := client.CollectionGroup("Tags").Documents(ctx).GetAll()

	if len(docRefs) == 0 {
		return []byte{}, codes.OK
	}

	resp := make(map[string]map[string]float64)
	for _, docRef := range docRefs {
		var placeData types.PlaceRequest
		_ = docRef.DataTo(&placeData)
		resp[placeData.PlaceId] = placeData.Tags
	}

	jsonStr, _ := json.Marshal(resp)
	return jsonStr, codes.OK
}

func (s PlaceController) GetAllTags() ([]byte, codes.Code) {

	resp := map[string][]string{"all_tags": types.TAGLIST}

	jsonStr, _ := json.Marshal(resp)
	return jsonStr, codes.OK
}
