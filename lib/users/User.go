package users

import (
	places "SS-Database/lib/places"
	types "SS-Database/lib/types"
	"SS-Database/utils"
	"cloud.google.com/go/firestore"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"slices"
	"strconv"
)

const SEARCH_LIMIT = 10
const UPDATE_STEP = 0.1

type UserController struct {
	UserID         string
	FeedbackId     int
	UserName       string
	Email          string
	FavoritePlaces map[string]string
}

func (s UserController) findSelectedTags(tagList []string) [][]float64 {
	var tagVals [][]float64
	for _, tag := range types.TAGLIST {
		if slices.Contains(tagList, tag) {
			tagVals = append(tagVals, []float64{1.0})
		} else {
			tagVals = append(tagVals, []float64{0.0})
		}
	}
	return tagVals
}

func (s UserController) findNextID(ctx context.Context, client *firestore.Client, path string, id string) int {
	var userdata UserController

	iter := client.Collection(path).OrderBy(id, firestore.Desc).Limit(1).Documents(ctx)
	q, _ := iter.GetAll()
	if q == nil {
		return 0
	} else {
		_ = q[0].DataTo(&userdata)

		return userdata.FeedbackId + 1
	}
}

func (s UserController) checkEmail(ctx context.Context, client *firestore.Client, newEmail string) codes.Code {
	iter := client.Collection("Users").Select("email").Documents(ctx)
	q, _ := iter.GetAll()
	for _, s := range q {
		if s.Data()["email"] == newEmail {
			return codes.AlreadyExists
		}
	}
	return codes.OK
}

func (s UserController) FilterPlaces(ctx context.Context, client *firestore.Client, data []byte) ([]byte, codes.Code) {
	filterInfo := new(types.FilterRequest)
	userInfo := new(types.UserRequest)
	placeController := places.NewPlaceController()
	fmt.Println(filterInfo)
	_ = json.Unmarshal(data, filterInfo)

	userInfoByte, _ := s.GetUserInfo(ctx, client, data)
	_ = json.Unmarshal(userInfoByte, userInfo)

	allTags := make(map[string]map[string]float64)
	allTagsByte, _ := placeController.GetAllTagVals(ctx, client)
	_ = json.Unmarshal(allTagsByte, &allTags)

	allTagsMatrix, placeIdList := utils.MapToMatrix(allTags)
	userPrefMatrix := utils.ListToMatrix(userInfo.Tags)
	selectedTags := s.findSelectedTags(filterInfo.Tags)

	filterMatrix := utils.MatrixSum(utils.MatrixConstMul(selectedTags, 0.8), utils.MatrixConstMul(userPrefMatrix, 0.2))

	filterRes := utils.MatrixTranspose(utils.MatrixMul(allTagsMatrix, filterMatrix))[0]

	placeFilter := utils.MapLists(filterRes, placeIdList)

	slices.Sort(filterRes)
	slices.Reverse(filterRes)

	var res []string
	for idx, vals := range filterRes {

		if idx > SEARCH_LIMIT {
			break
		}
		res = append(res, placeFilter[vals])
	}

	data, _ = json.Marshal(map[string][]string{"filtered_places": res})

	return data, codes.OK
}

func (s UserController) AddUser(ctx context.Context, client *firestore.Client, data []byte) codes.Code {
	userInfo := new(types.UserRequest)
	_ = json.Unmarshal(data, userInfo)
	fmt.Println(userInfo)

	errCode := s.checkEmail(ctx, client, userInfo.Email)
	if errCode != codes.OK {
		return errCode
	}

	userInfo.Tags = make(map[string]float64)
	// create default tags
	for _, tag := range types.TAGLIST {
		userInfo.Tags[tag] = 0.0
	}

	_, err := client.Collection("Users").Doc(userInfo.UserId).Create(ctx, map[string]interface{}{
		"userID":   userInfo.UserId,
		"userName": userInfo.UserName,
		"email":    userInfo.Email,
		"userType": userInfo.UserType,
		"tags":     userInfo.Tags,
	})
	if err != nil {
		log.Fatalf("Failed adding users: %v", err)
		return codes.Aborted
	}

	return codes.OK
}

func (s UserController) RemoveUser(ctx context.Context, client *firestore.Client, data []byte) codes.Code {
	userInfo := new(types.UserRequest)
	fmt.Println(userInfo)
	_ = json.Unmarshal(data, userInfo)
	ref := client.Collection("Users").Doc(userInfo.UserId)
	userReview := ref.Collection("Feedbacks")
	bulkWriter := client.BulkWriter(ctx)

	for {
		// Get a batch of documents
		iter := userReview.Limit(utils.BATCH_SIZE).Documents(ctx)
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

func (s UserController) UpdateUser(ctx context.Context, client *firestore.Client, data []byte) codes.Code {
	userInfo := new(types.UserRequest)
	err := json.Unmarshal(data, userInfo)
	if err != nil {
		log.Printf("Failed removing user: %v", err)
		return codes.Aborted
	}

	docSnap, err := client.Collection("Users").Doc(userInfo.UserId).Get(ctx)
	if err != nil {
		log.Printf("Failed removing user: User not found")
		return codes.NotFound
	}

	_, err = docSnap.Ref.Update(ctx, utils.ExtractNonEmptyFields(*userInfo))

	if err != nil {
		log.Printf("Failed removing user: %v", err)
		return codes.Aborted
	}
	return codes.OK
}

func (s UserController) GetUserInfo(ctx context.Context, client *firestore.Client, data []byte) ([]byte, codes.Code) {
	userInfo := new(types.UserRequest)
	err := json.Unmarshal(data, userInfo)
	fmt.Println(userInfo)

	q := client.Collection("Users").Where("userID", "==", userInfo.UserId)
	ref, err := q.Documents(ctx).GetAll()
	_ = ref[0].DataTo(userInfo)
	if err != nil {
		//log.Fatalf("Failed adding users: %v", err)
		return []byte{}, codes.Aborted
	}
	jsonStr, _ := json.Marshal(userInfo)
	return jsonStr, codes.OK
}

func (s UserController) UpdatePersonalization(tagMap map[string]float64, placeId string, favCount int, client *firestore.Client, add bool) map[string]float64 {
	var place = places.NewPlaceController()
	placeTags := make(map[string]float64)

	data, _ := json.Marshal(map[string]string{"place_id": placeId})

	tagMapByte, _ := place.GetTagByPlace(context.Background(), client, data)
	_ = json.Unmarshal(tagMapByte, &placeTags)

	for key, val := range tagMap {
		if add {
			tagMap[key] = (val*float64(favCount-1) + placeTags[key]) / float64(favCount)
		} else {
			if favCount == 0 {
				tagMap[key] = 0.0
			} else {
				tagMap[key] = (val*float64(favCount+1) - placeTags[key]) / float64(favCount)
			}
		}
	}
	return tagMap
}
func (s UserController) AddFavoritePlace(ctx context.Context, client *firestore.Client, data []byte) codes.Code {
	favReqInfo := new(types.UserRequest)

	var place = places.NewPlaceController()
	placeInfo := new(types.PlaceRequest)
	placeData, _ := place.GetPlaceName(ctx, client, data)
	err := json.Unmarshal(data, &favReqInfo)
	_ = json.Unmarshal(placeData, placeInfo)
	fmt.Println(favReqInfo)

	ref := client.Collection("Users").Doc(favReqInfo.UserId)
	err = client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {

		doc, err := tx.Get(ref)
		if err != nil {
			return err
		}
		var userdata types.UserRequest
		_ = doc.DataTo(&userdata)

		fmt.Println(userdata.FavoritePlaces)
		if userdata.FavoritePlaces == nil {
			userdata.FavoritePlaces = map[string]string{favReqInfo.PlaceId: placeInfo.PlaceName}
		} else {
			userdata.FavoritePlaces[favReqInfo.PlaceId] = placeInfo.PlaceName
		}
		userdata.Tags = s.UpdatePersonalization(userdata.Tags, favReqInfo.PlaceId, len(userdata.FavoritePlaces), client, true)
		return tx.Set(ref, map[string]interface{}{
			"favoritePlaces": userdata.FavoritePlaces,
			"tags":           userdata.Tags,
		}, firestore.Merge(firestore.FieldPath{"favoritePlaces"}, firestore.FieldPath{"tags"}))

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
	favReqInfo := new(types.UserRequest)
	err := json.Unmarshal(data, favReqInfo)
	fmt.Println(favReqInfo)

	ref := client.Collection("Users").Doc(favReqInfo.UserId)
	err = client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		doc, err := tx.Get(ref) // tx.Get, NOT ref.Get!
		if err != nil {
			return err
		}
		var userdata types.UserRequest
		_ = doc.DataTo(&userdata)
		fmt.Println(userdata.FavoritePlaces)

		delete(userdata.FavoritePlaces, favReqInfo.PlaceId)
		fmt.Println(userdata.FavoritePlaces)
		userdata.Tags = s.UpdatePersonalization(userdata.Tags, favReqInfo.PlaceId, len(userdata.FavoritePlaces), client, false)
		return tx.Set(ref, map[string]interface{}{
			"favoritePlaces": userdata.FavoritePlaces,
			"tags":           userdata.Tags,
		}, firestore.Merge(firestore.FieldPath{"favoritePlaces"}, firestore.FieldPath{"tags"}))

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

func (s UserController) AddFeedback(ctx context.Context, client *firestore.Client, data []byte) codes.Code {
	feedbackInfo := new(types.FeedbackRequest)

	err := json.Unmarshal(data, feedbackInfo)
	fmt.Println(feedbackInfo)

	feedbackPath := "Users/" + feedbackInfo.UserId + "/Feedbacks"
	ref := client.Collection(feedbackPath).NewDoc()
	feedbackInfo.FeedbackId = ref.ID

	_, err = ref.Create(ctx, map[string]interface{}{
		"feedbackId": feedbackInfo.FeedbackId,
		"rating":     feedbackInfo.Rating,
		"placeId":    feedbackInfo.PlaceId,
	})

	tags := make(map[string]float64)
	for tag, val := range feedbackInfo.Rating {
		satisfaction, _ := strconv.Atoi(val)
		tags[tag] = UPDATE_STEP * float64(satisfaction)
	}

	placeController := places.NewPlaceController()
	placeController.UpdateTags(feedbackInfo.PlaceId, client, tags, false)

	if err != nil {
		// Handle any errors appropriately in this section.
		log.Printf("An error has occurred: %s", err)
		return codes.Aborted
	}
	return codes.OK
}

func (s UserController) GetFeedbacks(ctx context.Context, client *firestore.Client, data []byte) ([]byte, codes.Code) {
	userInfo := new(types.UserRequest)

	_ = json.Unmarshal(data, userInfo)
	fmt.Println(userInfo)

	feedbackPath := "Users/" + userInfo.UserId + "/Feedbacks"
	feedbackSnaps, _ := client.Collection(feedbackPath).Documents(ctx).GetAll()
	res := make(map[string]types.FeedbackRequest, 0)

	for _, fbSnap := range feedbackSnaps {
		var fb types.FeedbackRequest
		_ = fbSnap.DataTo(&fb)
		fb.UserId = userInfo.UserId
		res[fb.FeedbackId] = fb
	}
	jsonStr, _ := json.Marshal(res)
	return jsonStr, codes.OK
}

func (s UserController) GetAllUsers(ctx context.Context, client *firestore.Client) ([]byte, codes.Code) {
	docRefs, err := client.Collection("Users").OrderBy("userID", firestore.Asc).Documents(ctx).GetAll()
	if err != nil {
		return []byte{}, codes.NotFound
	}

	if len(docRefs) == 0 {
		return []byte{}, codes.OK
	}

	resp := map[string]types.UserRequest{}
	for _, docRef := range docRefs {
		var userData types.UserRequest
		_ = docRef.DataTo(&userData)
		resp[userData.UserId] = userData
	}

	jsonStr, err := json.Marshal(resp)
	return jsonStr, codes.OK
}
