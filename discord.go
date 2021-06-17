package main

import (
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

func updateDiscordSeen(discordUserID string) {
}

func discordTypingEventHandler(discord *discordgo.Session, event *discordgo.TypingStart) {
	maybeUpdateSeen(event.UserID)
}

func discordMessageEventHandler(discord *discordgo.Session, event *discordgo.MessageCreate) {
	maybeUpdateSeen(event.Author.ID)
}

func discordPressenceEventHandler(discord *discordgo.Session, event *discordgo.PresenceUpdate) {
	maybeUpdateSeen(event.User.ID)
}

func mindDiscord() {
	guildID := os.Getenv("GUILD")
	token := os.Getenv("APPTOKEN")
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	discord.AddHandler(discordMessageEventHandler)
	discord.AddHandler(discordPressenceEventHandler)
	discord.AddHandler(discordTypingEventHandler)

	if err := discord.Open(); err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second)

	// main loop, runs every 300s
	for {

		// get the paginated list of guild members
		members := []*discordgo.Member{}
		var after string
		for {
			ms, err := discord.GuildMembers(guildID, after, 1000)
			if err != nil {
				panic(err)
			}

			// if we have no members, then there are no more
			if len(ms) == 0 {
				break
			}

			members = append(members, ms...)

			// if we have 1000, then we may have more to retrieve, so update the after value
			// but if we have <1000, then this is the last page
			if len(ms) == 1000 {
				lastmember := ms[len(ms)-1]
				after = lastmember.User.ID
			} else {
				break
			}
		}

		for _, v := range members {
			if v.User.Bot {
				continue
			}
			nickname := v.User.Username
			if v.Nick != "" {
				nickname = v.Nick
			}
			maybeInsertUser(v.User.ID, nickname)
		}
		time.Sleep(300 * time.Second)
	}
}
