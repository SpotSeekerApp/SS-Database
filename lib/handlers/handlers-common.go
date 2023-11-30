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

const userReqTypes = "user"
const placeReqTypes = "place|review"

func (h HandlerInstance) HandleRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	reqPath := strings.ToLower(r.URL.Path)
	var err http.ConnState

	if utils.SubstrInList(reqPath, userReqTypes) {
		err = h.UserHandler(w, r)
	} else if utils.SubstrInList(reqPath, placeReqTypes) {
		err = h.PlaceHandler(w, r)
	} else {
		err = http.StatusNotImplemented
	}
	w.WriteHeader(int(err))
}
