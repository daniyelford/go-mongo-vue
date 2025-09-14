package handlers

import (
	"encoding/json"
	"net/http"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users := []map[string]string{
		{"id": "1", "name": "Daniyal"},
		{"id": "2", "name": "Ali"},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
