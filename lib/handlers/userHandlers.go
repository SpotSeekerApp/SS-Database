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

type UserInterface interface {
	AddUser(ctx context.Context, client *firestore.Client, data []byte) codes.Code
	RemoveUser(ctx context.Context, client *firestore.Client, data []byte) codes.Code
	UpdateUser(ctx context.Context, client *firestore.Client, data []byte) codes.Code
	AddFavoritePlace(ctx context.Context, client *firestore.Client, data []byte) codes.Code
	RemoveFavoritePlace(ctx context.Context, client *firestore.Client, data []byte) codes.Code
	GetAllUsers(ctx context.Context, client *firestore.Client) ([]byte, codes.Code)
}

func (h HandlerInstance) UserHandler(w http.ResponseWriter, r *http.Request) {
	reqFunction := strings.Split(r.URL.Path, "/")[1]

	resp := make(map[string]interface{})
	var err http.ConnState
	var res []byte

	if r.Method == "POST" {
		err = h.POSTUserHandler(w, r, reqFunction)
	} else if r.Method == "GET" {
		res, err = h.GETUserHandler(w, r, reqFunction)
		_ = json.Unmarshal(res, &resp)
	} else {
		err = http.StatusNotImplemented
	}

	w.Header().Set("Content-Type", "application/json")

	resp["StatusCode"] = err
	jsonResp, _ := json.Marshal(resp)
	_, _ = w.Write(jsonResp)

}

func (h HandlerInstance) POSTUserHandler(w http.ResponseWriter, r *http.Request, function string) http.ConnState {
	req := new(types.UserRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return http.StatusBadRequest
	}
	data, _ := json.Marshal(req)
	fmt.Println(data)

	var err http.ConnState
	if function == "AddUser" {
		err = utils.MapErrorCode(h.UserController.AddUser(context.Background(), h.Client, data))
	} else if function == "UpdateUser" {
		err = utils.MapErrorCode(h.UserController.UpdateUser(context.Background(), h.Client, data))
	} else if function == "RemoveUser" {
		err = utils.MapErrorCode(h.UserController.RemoveUser(context.Background(), h.Client, data))
	} else {
		return http.StatusNotImplemented
	}

	return err
}

func (h HandlerInstance) GETUserHandler(w http.ResponseWriter, r *http.Request, function string) ([]byte, http.ConnState) {
	//req := new(types.UserRequest)
	//if err := json.NewDecoder(r.Body).Decode(req); err != nil {
	//	return http.StatusBadRequest
	//}
	//data, _ := json.Marshal(req)
	//fmt.Println(data)

	var err http.ConnState
	var res []byte
	if function == "GetAllUsers" {
		resTemp, errTemp := h.UserController.GetAllUsers(context.Background(), h.Client)
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
