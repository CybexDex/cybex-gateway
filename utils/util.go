package utils

import (
	"encoding/json"
	"net/http"
)

//Message ...
func Message(status bool, message string, d ...interface{}) map[string]interface{} {
	if len(d) == 1 {
		return map[string]interface{}{"status": status, "message": message, "data": d[0]}
	}

	return map[string]interface{}{"status": status, "message": message}
}

//Respond ...
func Respond(w http.ResponseWriter, data map[string]interface{}, s ...int) {
	status := http.StatusOK
	if len(s) == 1 {
		status = s[0]
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
