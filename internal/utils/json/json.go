package json

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteResponse(w http.ResponseWriter, resp interface{}) {
    jsonBody, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("json.WriteResponse - failed to marshal response body:", err)
		return
	}
    w.Write(jsonBody)
}

func WriteError(w http.ResponseWriter, err error, code int) {
    // w.Header().Set("Content-Type", "application/json; charset=utf-8")
    // w.Header().Set("X-Content-Type-Options", "nosniff")
    jsonErr := struct {
        Text string `json:"text"`
    }{
        Text: err.Error(),
    }
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(jsonErr)
}
