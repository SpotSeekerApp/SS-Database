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
	AddFavoritePlace(ctx context.Context, client *firestore.Client, data []byte) codes.Code
	RemoveFavoritePlace(ctx context.Context, client *firestore.Client, data []byte) codes.Code
}

func (h HandlerInstance) UserHandler(w http.ResponseWriter, r *http.Request) http.ConnState {
	reqFunction := strings.Split(r.URL.Path, "/")[1]
	var err http.ConnState
	if r.Method == "POST" {
		err = h.POSTUserHandler(w, r, reqFunction)
	}
	fmt.Fprintf(w, "%s", err)
	return err
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
	} else {
		return http.StatusNotImplemented
	}

	return err
}
