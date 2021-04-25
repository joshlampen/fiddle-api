package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/JoshLampen/fiddle/api/db"
	"github.com/JoshLampen/fiddle/api/db/model"
	"github.com/JoshLampen/fiddle/api/internal/constant"
)

func TokensGet(w http.ResponseWriter, r *http.Request, store *db.Store) {
    w.Header().Set("Content-Type", "application/json")

    // Get auth ID from url
    authID := r.URL.Query().Get(constant.URLParamAuthID)

    // Get from database
    token, err := store.TokenStore.GetByID(authID)
    if err != nil {
        fmt.Println("handler.TokensGet - failed to get token from database:", err)
        return
    }

    // Send a response
	jsonBody, err := json.Marshal(token)
	if err != nil {
		fmt.Println("handler.TokensGet - failed to marshal response body:", err)
		return
	}
	w.Write(jsonBody)
}

// TokensCreate is an HTTP handler for inserting a token into the database
func TokensCreate(w http.ResponseWriter, r *http.Request, store *db.Store) {
	// Read the request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("handler.TokensCreate - failed to read request body:", err)
		return
	}

	var token model.Token
	if err := json.Unmarshal(body, &token); err != nil {
		fmt.Println("handler.TokensCreate - failed to unmarshal request body:", err)
		return
	}

    // Insert into database
    _, err = store.TokenStore.Create(token)
    if err != nil {
        fmt.Println("handler.TokensCreate - failed to create token:", err)
        return
    }
}
