package handlers

import (
	"SS-Database/utils"
	"cloud.google.com/go/firestore"
	"fmt"
	"net/http"
	"strings"
)

type HandlerInstance struct {
	UserController  UserInterface
	PlaceController PlaceInterface
	Client          *firestore.Client
}

const userReqTypes = "user|update|getall|favorite"
const placeReqTypes = "place|review|placeinfo"

func (h HandlerInstance) HandleRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	reqPath := strings.ToLower(r.URL.Path)

	if utils.SubstrInList(reqPath, userReqTypes) {
		h.UserHandler(w, r)
	} else if utils.SubstrInList(reqPath, placeReqTypes) {
		h.PlaceHandler(w, r)
	}
}
