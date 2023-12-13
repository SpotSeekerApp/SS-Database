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
	"strconv"
	"strings"
)

type PlaceInterface interface {
	AddPlace(ctx context.Context, client *firestore.Client, data []byte) codes.Code
	GetPlaceInfo(ctx context.Context, client *firestore.Client, data []byte) ([]byte, codes.Code)
	GetPlaceName(ctx context.Context, client *firestore.Client, data []byte) ([]byte, codes.Code)
	GetAllPlaces(ctx context.Context, client *firestore.Client) ([]byte, codes.Code)
}

func (h HandlerInstance) PlaceHandler(w http.ResponseWriter, r *http.Request) {
	reqFunction := strings.Split(r.URL.Path, "/")[1]

	resp := make(map[string]interface{})
	var err http.ConnState
	var res []byte

	if r.Method == "POST" {
		err = h.POSTPlaceHandler(w, r, reqFunction)
	} else if r.Method == "GET" {
		res, err = h.GETPlaceHandler(w, r, reqFunction)
		_ = json.Unmarshal(res, &resp)
	} else {
		err = http.StatusNotImplemented
	}

	w.Header().Set("Content-Type", "application/json")

	resp["StatusCode"] = err
	jsonResp, _ := json.Marshal(resp)
	_, _ = w.Write(jsonResp)

}

func (h HandlerInstance) POSTPlaceHandler(w http.ResponseWriter, r *http.Request, function string) http.ConnState {

	req := new(types.PlaceRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return http.StatusBadRequest
	}
	data, _ := json.Marshal(req)
	fmt.Println(data)

	var err http.ConnState
	if function == "AddPlace" {
		err = utils.MapErrorCode(h.PlaceController.AddPlace(context.Background(), h.Client, data))
	} else {
		return http.StatusNotImplemented
	}

	return err
}

func (h HandlerInstance) GETPlaceHandler(w http.ResponseWriter, r *http.Request, function string) ([]byte, http.ConnState) {
	req := make(map[string]string)
	if strings.Contains(r.URL.RawQuery, "place_id") {
		req["place_id"] = r.URL.Query()["place_id"][0]
	}
	data, _ := json.Marshal(req)
	fmt.Println(req)

	var err http.ConnState
	var res []byte
	if function == "GetPlaceInfo" {
		resTemp, errTemp := h.PlaceController.GetPlaceInfo(context.Background(), h.Client, data)
		err = utils.MapErrorCode(errTemp)
		res = resTemp
	} else if function == "GetAllPlaces" {
		resTemp, errTemp := h.PlaceController.GetAllPlaces(context.Background(), h.Client)
		err = utils.MapErrorCode(errTemp)
		res = resTemp
	} else if function == "GetPlaceName" {
		resTemp, errTemp := h.PlaceController.GetPlaceName(context.Background(), h.Client, data)
		err = utils.MapErrorCode(errTemp)
		res = resTemp
	} else {
		return []byte{}, http.StatusNotImplemented
	}

	resp := make(map[string]map[string]interface{})
	respData := make(map[string]interface{})
	_ = json.Unmarshal(res, &respData)
	resp["Data"] = respData
	res, _ = json.Marshal(resp)

	return res, err
}
