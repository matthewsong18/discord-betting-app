package main

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

func (bot *Bot) RegisterCommands() error {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "create-poll",
			Description: "Create a new poll",
		},
	}

	AppID = os.Getenv("APP_ID")
	GuildID = os.Getenv("GUILD_ID")
	_, err := bot.DiscordSession.ApplicationCommandBulkOverwrite(AppID, GuildID, commands)
	if err != nil {
		log.Printf("Error overwriting commands: %v", err)
		return err
	}

	log.Println("Commands successfully registered.")
	return nil
}
