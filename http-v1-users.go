package main

import (
	"encoding/json"
	"net/http"
	"time"
)

func init() {
	router.HandleFunc("/v1/users.json", func(w http.ResponseWriter, r *http.Request) {
		var rval = map[string]struct {
			User privateUser
			Seen time.Time
		}{}
		var maxAge = time.Now().Add(0 - (30 * 24 * time.Hour))
		for id, user := range userList {
			t := time.Unix(user.Seen, 0)
			if t.Before(maxAge) {
				continue
			}
			rval[id] = struct {
				User privateUser
				Seen time.Time
			}{
				User: *user,
				Seen: t,
			}
		}
		json.NewEncoder(w).Encode(rval)
	})
}
