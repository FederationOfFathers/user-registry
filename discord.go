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

	time.Sleep(5)
	for {
		guild, err := discord.Guild(guildID)
		if err != nil {
			panic(err)
		}

		m := guild.Members
		for _, v := range m {
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
