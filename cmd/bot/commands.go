package main

import "github.com/bwmarrin/discordgo"

func (bot *Bot) RegisterCommands() error {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "create-poll",
			Description: "Create a new poll",
		},
		{
			Name:        "bet",
			Description: "Place a bet on an open poll",
		},
	}

	_, err := bot.DiscordSession.ApplicationCommandBulkOverwrite(bot.DiscordSession.State.Application.ID, bot.DiscordSession.State.Application.GuildID, commands)
	return err
}
