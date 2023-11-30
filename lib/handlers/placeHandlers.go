package handlers

import (
	types "SS-Database/lib/types"
	"SS-Database/utils"
	"cloud.google.com/go/firestore"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"net/http"
	"strings"
)

type PlaceInterface interface {
	AddPlace(ctx context.Context, client *firestore.Client, data []byte) codes.Code
	AddReview(ctx context.Context, client *firestore.Client, data []byte) codes.Code
}

func (h HandlerInstance) PlaceHandler(w http.ResponseWriter, r *http.Request) http.ConnState {
	reqFunction := strings.Split(r.URL.Path, "/")[1]
	var err http.ConnState

	if reqFunction == "AddPlace" {
		err = h.AddPlaceHandler(w, r)
	} else if reqFunction == "AddReview" {
		err = h.AddReviewHandler(w, r)
	} else {
		return http.StatusNotImplemented
	}
	fmt.Fprintf(w, "%s", err)
	return err
}

func (h HandlerInstance) AddPlaceHandler(w http.ResponseWriter, r *http.Request) http.ConnState {
	req := new(types.PlaceRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return http.StatusBadRequest
	}
	data, _ := json.Marshal(req)
	fmt.Println(data)
	err := utils.MapErrorCode(h.PlaceController.AddPlace(context.Background(), h.Client, data))

	return err
}

func (h HandlerInstance) AddReviewHandler(w http.ResponseWriter, r *http.Request) http.ConnState {
	req := new(types.ReviewRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return http.StatusBadRequest
	}
	data, _ := json.Marshal(req)
	fmt.Println(data)
	err := utils.MapErrorCode(h.PlaceController.AddReview(context.Background(), h.Client, data))

	return err
}
