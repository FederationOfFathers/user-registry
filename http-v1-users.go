package main

import (
	"encoding/json"
	"net/http"
	"time"
)

func init() {
	router.HandleFunc("/v1/users.json", func(w http.ResponseWriter, r *http.Request) {
		//var u = userList
		// var s = seenList
		var rval = map[string]struct {
			User privateUser
			Seen time.Time
		}{}
		/*
			var maxAge = time.Now().Add(0 - (30 * 24 * time.Hour))
			for id, t := range s {
				if id == "USLACKBOT" {
					continue
				}
				if t.Before(maxAge) {
					continue
				}
				rval[id] = struct {
					User privateUser
					Seen time.Time
				}{
					User: u[id],
					Seen: t,
				}
			}
		*/
		json.NewEncoder(w).Encode(rval)
	})
}
