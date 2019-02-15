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
			userList[id].DisplayName = nickname
			userList[id].Name = nickname
			_, err := updateName.Exec(nickname, id)
			if err != nil {
				log.Println("error updating nickname", err.Error(), id, nickname)
			} else {
				log.Println("updated nickname", id, nickname)
			}
		}
		if u.Thumb != thumb {
			userList[id].Thumb = thumb
			_, err := insertUserMeta.Exec(id, "thumb", thumb)
			if err != nil {
				log.Println("error updating nickname", err.Error(), id, thumb)
			} else {
				log.Println("updated nickname", id, thumb)
			}

		}
		if u.Image != image {
			userList[id].Image = image
			_, err := insertUserMeta.Exec(id, "image", image)
			if err != nil {
				log.Println("error updating image", err.Error(), id, image)
			} else {
				log.Println("updated image", id, image)
			}
		}
		return
	}
	_, err := insertUser.Exec(discordID, nickname)
	if err != nil {
		log.Println("insertUser.Exec error", err.Error(), "for", discordID, nickname)
		return
	}
	u := new(user)
	if err := u.fromRow(getDiscordUser.QueryRow(discordID)); err != nil {
		log.Println("u.fromRows error", err.Error(), discordID)
		return
	}
	p := u.privateUser()
	userList[p.ID] = p
	log.Println("new user", p)
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
		_, err := updateSeen.Exec(now, id)
		if err != nil {
			log.Println("error updating seen", id, now)
		} else {
			log.Println("updated seen", userList[id].Name, now)
		}
		return
	}
}

func loadUsersFromDatabase() {
	for {
		if rows, err := getAllUsers.Query(); err == nil {
			privateUsersLock.Lock()
			defer privateUsersLock.Unlock()
			defer rows.Close()
			for rows.Next() {
				u := new(user)
				if err := u.fromRows(rows); err != nil {
					log.Println("u.fromRow error", err.Error())
					continue
				}
				p := u.privateUser()
				userList[p.ID] = p
			}
			return
		}
		time.Sleep(5 * time.Second)
	}
}
