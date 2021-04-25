package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/JoshLampen/fiddle/api/db"
	"github.com/JoshLampen/fiddle/api/db/model"
)

// UsersCreate is an HTTP handler for inserting a user into the database
func UsersCreate(w http.ResponseWriter, r *http.Request, store *db.Store) {
    w.Header().Set("Content-Type", "application/json")

	// Read the request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("handler.UsersCreate - failed to read request body:", err)
		return
	}

	var user model.User
	if err := json.Unmarshal(body, &user); err != nil {
		fmt.Println("handler.UsersCreate - failed to unmarshal request body:", err)
		return
	}

	// Insert into database
	result, err := store.UserStore.Create(user)
	if err != nil {
		fmt.Println("handler.UsersCreate - failed to create user:", err)
		return
	}

	// Send a response
	user.ID = result.ID
	jsonBody, err := json.Marshal(user)
	if err != nil {
		fmt.Println("handler.UsersCreate - failed to marshal response body:", err)
		return
	}
	w.Write(jsonBody)
}
