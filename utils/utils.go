package utils

import "net/http"

func SendResponse(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(message))
}
