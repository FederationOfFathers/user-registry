package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

var userList = privateUsers{}

type privateUserList struct {
	Members []struct {
		Deleted         bool   `json:"deleted"`
		Bot             bool   `json:"is_bot"`
		Restricted      bool   `json:"is_restricted"`
		UltraRestricted bool   `json:"is_ultra_restricted"`
		ID              string `json:"id"`
		Name            string `json:"name"`
		TimeZone        string `json:"tz"`
		Profile         struct {
			GamerTag    string `json:"first_name"`
			DisplayName string `json:"real_name_normalized"`
			Image       string `json:"image_original"`
			Thumb       string `json:"image_24"`
		} `json:"profile"`
	} `json:"members"`
}

type privateUser struct {
	ID          string
	Name        string
	GamerTag    string
	DisplayName string
	Image       string
	Thumb       string
}

type privateUsers map[string]privateUser

func mindPrivateUserList() {
	go updateUserCacheV2()
	if u, err := getPrivateUserList(); err != nil {
		log.Fatal("Error fetching initial user list!", err.Error())
	} else {
		userList = u
		updateUserCacheV2()
	}
	t := time.Tick(10 * time.Minute)
	for {
		select {
		case <-t:
			for {
				log.Println("Updating user list")
				if u, err := getPrivateUserList(); err != nil {
					log.Println("Error fetching user list:", err.Error())
					time.Sleep(30)
				} else {
					userList = u
					updateUserCacheV2()
					break
				}
			}
		}
	}
}

func updateUserCacheV2() {
	log.Println("User list updated, updating MySQL User Cache (v2)")
	rows, err := db.getAllUsers.Query()
	if err != nil {
		log.Println("Error fetching all users", err)
		return
	}
	for rows.Next() {
		var user = new(user)
		if err := user.fromRows(rows); err != nil {
			log.Println("Error scanning user", err)
			continue
		}
		userCache[user.SlackID] = user
	}
	rows.Close()
	return
}

func getPrivateUserList() (privateUsers, error) {
	var rval = privateUsers{}
	var raw privateUserList
	rsp, err := http.Get("http://127.0.0.1:8879/users.json")
	if err != nil {
		return rval, err
	}
	defer rsp.Body.Close()
	dec := json.NewDecoder(rsp.Body)
	if err := dec.Decode(&raw); err != nil {
		return rval, err
	}
	for _, user := range raw.Members {
		go db.maybeInsert(user.ID, user.Name, user.Profile.GamerTag, user.TimeZone)
		if user.Bot {
			continue
		}
		if user.Deleted {
			continue
		}
		if user.Restricted {
			continue
		}
		if user.UltraRestricted {
			continue
		}
		rval[user.ID] = privateUser{
			ID:          user.ID,
			Name:        user.Name,
			GamerTag:    user.Profile.GamerTag,
			DisplayName: user.Profile.DisplayName,
			Image:       user.Profile.Image,
			Thumb:       user.Profile.Thumb,
		}
	}
	return rval, nil
}
