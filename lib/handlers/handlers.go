package handlers

import (
	"cloud.google.com/go/firestore"
	"fmt"
	"golang.org/x/net/context"
	"net/http"
	"strings"
)

type HandlerInstance struct {
	UserController UserInterface
	Client         *firestore.Client
}

type UserInterface interface {
	AddUser(ctx context.Context, client *firestore.Client, data []byte) error
}

func (h HandlerInstance) UserHandler(w http.ResponseWriter, r *http.Request) {
	reqFunction := strings.Split(r.URL.Path, "/")[1]
	var err error
	if strings.HasPrefix(reqFunction, "AddUser") {
		err = h.UserController.AddUser(context.Background(), h.Client, []byte(r.Header["Data"][0]))
	}
	fmt.Fprintf(w, "%s", err)
}
