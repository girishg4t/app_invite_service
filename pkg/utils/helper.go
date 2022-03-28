package httputils

import (
	"encoding/json"
	"net/http"
)

// SendSuccessHeader send's the successful message after correct request
func SendSuccessHeader(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(data)
}
