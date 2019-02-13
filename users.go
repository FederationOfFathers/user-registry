package main

import (
	"log"
	"sync"
	"time"
)

var userList = privateUsers{}

var privateUsersLock sync.RWMutex

type privateUser struct {
	ID          string
	DiscordID   string
	Name        string
	GamerTag    string
	DisplayName string
	Image       string
	Thumb       string
	Seen        int64
}

type privateUsers map[string]*privateUser

func maybeInsertUser(discordID string, nickname string) {
	log.Println("maybeInsertUser", discordID, nickname)
	privateUsersLock.Lock()
	defer privateUsersLock.Unlock()
	for id, u := range userList {
		if u.DiscordID != discordID {
			continue
		}
		if u.Name == nickname || u.DisplayName == nickname {
			return
		}
		userList[id].DisplayName = nickname
		userList[id].Name = nickname
		// update database
		return
	}
	_, err := insertUser.Exec(discordID, nickname)
	if err != nil {
		log.Println("insertUser.Exec error", err.Error(), "for", discordID, nickname)
		return
	}
	// New User, insert get ID, add to map
}

func maybeUpdateSeen(discordID string) {
	log.Println("maybeUpdateSeen", discordID)
	privateUsersLock.Lock()
	defer privateUsersLock.Unlock()
	for id, u := range userList {
		if u.DiscordID != discordID {
			continue
		}
		now := time.Now().Unix()
		if now > (u.Seen + 60) {
			return
		}
		userList[id].Seen = now
		// update database
		return
	}
}

func loadUsersFromDatabase() {
	// Do the thing. On startup only
}
