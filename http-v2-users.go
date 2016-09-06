package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func init() {
	router.Handle("/v2/users.json", mw(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(userCache)
	}))
}

func init() {
	router.Path("/v2/users/{userid}.json").Methods("PUT").Handler(mw(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.RemoteAddr, "127.0.0.1:") {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		userid, ok := mux.Vars(r)["userid"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		user, ok := userCache[userid]
		if !ok {
			http.NotFound(w, r)
			return
		}
		var values map[string]string
		err := json.NewDecoder(r.Body).Decode(&values)
		if err != nil || values == nil || len(values) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		for key, value := range values {
			switch strings.ToLower(key) {
			case "tz":
				if user.TZ != value {
					_, err := db.updateTz.Exec(value, user.SlackID)
					if err == nil {
						user.TZ = value
					}
				}
			case "xbl":
				if user.XBL != value {
					_, err := db.updateXbl.Exec(value, user.SlackID)
					if err != nil {
						user.XBL = value
					}
				}
			case "name", "destinyid", "userid", "memberid", "seentimestamp", "seen":
				continue
			default:
				// TODO: membermeta
			}
		}
		enc := json.NewEncoder(w)
		enc.Encode(user)
		enc.Encode(values)
	}))
}
