package main

import (
	"encoding/json"
	"net/http"
)

func init() {
	router.HandleFunc("/v2/users.json", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(userCache)
	})
}
