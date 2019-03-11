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

func maybeInsertUser(discordID, nickname, image, thumb string) {
	privateUsersLock.Lock()
	defer privateUsersLock.Unlock()
	for id, u := range userList {
		if u.DiscordID != discordID {
			continue
		}
		if u.Name != nickname || u.DisplayName != nickname {
			oldName := userList[id].Name
			userList[id].DisplayName = nickname
			userList[id].Name = nickname
			if _, err := updateName.Exec(nickname, id); err != nil {
				log.Printf("Error updating name for %s to %s: %s", userList[id].Name, nickname, err.Error())
			} else {
				log.Printf("Updated name for %s to %s", oldName, userList[id].Name)
			}
		}
		if u.Thumb != thumb {
			userList[id].Thumb = thumb
			if _, err := insertUserMeta.Exec(id, "thumb", thumb); err != nil {
				log.Printf("Error updating thumb for %s to %s: %s", userList[id].Name, thumb, err.Error())
			} else {
				log.Printf("Updated thumb for %s to %s", userList[id].Name, thumb)
			}
		}
		if u.Image != image {
			userList[id].Image = image
			if _, err := insertUserMeta.Exec(id, "image", image); err != nil {
				log.Printf("Error updating image for %s to %s: %s", userList[id].Name, image, err.Error())
			} else {
				log.Printf("Updated image for %s to %s", userList[id].Name, image)
			}
		}
		return
	}
	if _, err := insertUser.Exec(discordID, nickname); err != nil {
		log.Printf("Error inserting user %s (%s): %s", discordID, nickname, err.Error())
		return
	}
	u := new(user)
	if err := u.fromRow(getDiscordUser.QueryRow(discordID)); err != nil {
		log.Printf("Error querying user %s (%s) after insert: %s", discordID, nickname, err.Error())
		return
	}
	p := u.privateUser()
	userList[p.ID] = p
	userList[p.ID].Thumb = thumb
	if _, err := insertUserMeta.Exec(p.ID, "thumb", thumb); err != nil {
		log.Printf("Error updating thumb for inserted user %s to %s: %s", userList[p.ID].Name, thumb, err.Error())
	} else {
		log.Printf("Updated thumb for inserted user %s to %s", userList[p.ID].Name, thumb)
	}
	userList[p.ID].Image = image
	if _, err := insertUserMeta.Exec(p.ID, "image", image); err != nil {
		log.Printf("Error updating image for inserted user %s to %s: %s", userList[p.ID].Name, image, err.Error())
	} else {
		log.Printf("Updated image for inserted user %s to %s", userList[p.ID].Name, image)
	}
	log.Printf("Created new user id %s: %s, discord: %s, image: %s", p.ID, p.Name, p.DiscordID, p.Image)
}

func maybeUpdateSeen(discordID string) {
	privateUsersLock.Lock()
	defer privateUsersLock.Unlock()
	for id, u := range userList {
		if u.DiscordID != discordID {
			continue
		}
		now := time.Now().Unix()
		if now < (u.Seen + 60) {
			return
		}
		userList[id].Seen = now
		if _, err := updateSeen.Exec(now, id); err != nil {
			log.Printf("Error updating seen for %s (%s): %s", userList[id].Name, id, err.Error())
		} else {
			log.Printf("Updated seen for %s (%s): %d", userList[id].Name, id, now)
		}
		return
	}
}

func loadUsersFromDatabase() {
	for {
		go func(){
			if rows, err := getAllUsers.Query(); err == nil {
				privateUsersLock.Lock()
				defer privateUsersLock.Unlock()
				defer rows.Close()
				for rows.Next() {
					u := new(user)
					if err := u.fromRows(rows); err != nil {
						log.Printf("Error loading user from database: %s", err.Error())
						continue
					}
					p := u.privateUser()
					userList[p.ID] = p
				}
				return
			}
		}()
		time.Sleep(5 * time.Second) // should probably use a tick/chan?
	}
}
