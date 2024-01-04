package handlers

import (
	"SS-Database/utils"
	"cloud.google.com/go/firestore"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type HandlerInstance struct {
	UserController  UserInterface
	PlaceController PlaceInterface
	Client          *firestore.Client
}

const userReqTypes = "user|feedback|favorite|password|filter|review"
const placeReqTypes = "place|tag"

func (h HandlerInstance) VerifyAPICall(w http.ResponseWriter, r *http.Request) bool {

	fileData := make(map[string]string, 0)
	jsonFile, _ := os.ReadFile("api-key.json")
	json.Unmarshal(jsonFile, &fileData)
	return r.Header.Get("API-Key") == fileData["api-key"]
}

func (h HandlerInstance) HandleRequest(w http.ResponseWriter, r *http.Request) {
	if h.VerifyAPICall(w, r) != true {
		w.Header().Set("Content-Type", "application/json")

		jsonResp, _ := json.Marshal(map[string]http.ConnState{
			"StatusCode": http.StatusUnauthorized,
		})
		_, _ = w.Write(jsonResp)
		return
	}
	fmt.Println(r.URL.Path)
	reqPath := strings.ToLower(r.URL.Path)

	if utils.SubstrInList(reqPath, userReqTypes) {
		h.UserHandler(w, r)
	} else if utils.SubstrInList(reqPath, placeReqTypes) {
		h.PlaceHandler(w, r)
	}
}
