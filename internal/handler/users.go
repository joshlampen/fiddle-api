package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/JoshLampen/fiddle/api/db"
	"github.com/JoshLampen/fiddle/api/db/model"
	"github.com/JoshLampen/fiddle/api/internal/constant"
	jsonWriter "github.com/JoshLampen/fiddle/api/internal/utils/json"
)

func UsersGet(w http.ResponseWriter, r *http.Request, store *db.Store) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,OPTIONS")

    // Get auth ID from url
    authID := r.URL.Query().Get(constant.URLParamAuthID)

    user, err := store.UserStore.GetByAuthID(authID)
    if err != nil {
        err := fmt.Errorf("Failed to get user: %w", err)
        jsonWriter.WriteError(w, err, http.StatusInternalServerError)
        return
    }

    jsonWriter.WriteResponse(w, user)
}

func UsersGetFriends(w http.ResponseWriter, r *http.Request, store *db.Store) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,OPTIONS")

    userID := r.URL.Query().Get(constant.URLParamUserID)

    friends, err := store.UserStore.GetFriendsByUserID(userID)
    if err != nil {
        err := fmt.Errorf("Failed to get friends: %w", err)
        jsonWriter.WriteError(w, err, http.StatusInternalServerError)
        return
    }

    jsonWriter.WriteResponse(w, friends)
}

// UsersCreate is an HTTP handler for inserting a user into the database
func UsersCreate(w http.ResponseWriter, r *http.Request, store *db.Store) {
    w.Header().Set("Content-Type", "application/json")

	// Read the request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
        err := fmt.Errorf("Failed to read user request: %w", err)
        jsonWriter.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	var user model.User
	if err := json.Unmarshal(body, &user); err != nil {
        err := fmt.Errorf("Failed to process user request: %w", err)
        jsonWriter.WriteError(w, err, http.StatusInternalServerError)
		return
	}

    // Return user from the database if it exists
    result, err := store.UserStore.GetBySpotifyIDIfExists(user.SpotifyID)
	if err != nil {
        err := fmt.Errorf("Failed to get user: %w", err)
        jsonWriter.WriteError(w, err, http.StatusInternalServerError)
		return
	}
    if result == nil {
        // Insert into database
        result, err = store.UserStore.Create(user)
        if err != nil {
            err := fmt.Errorf("Failed to create user: %w", err)
            jsonWriter.WriteError(w, err, http.StatusInternalServerError)
            return
        }
    }

    jsonWriter.WriteResponse(w, result)
}
