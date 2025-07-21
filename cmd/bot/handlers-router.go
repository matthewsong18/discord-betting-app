package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

func (bot *Bot) interactionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		bot.handleSlashCommand(s, i)
	case discordgo.InteractionModalSubmit:
		// This is a modal submission
		bot.handleModalSubmit(s, i)
	case discordgo.InteractionMessageComponent:
		bot.handleButtonPress(s, i)
	default:
		log.Printf("Unknown interaction type received: %v", i.Type)
	}
}

func (bot *Bot) handleSlashCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	commandName := i.ApplicationCommandData().Name
	switch commandName {
	case "create-poll":
		bot.handleCreatePollCommand(s, i)
	case "bet":
		// bot.handleBetCommand(s, i) // Future implementation
	default:
		log.Printf("Unknown slash command received: %s", commandName)
	}
}

func (bot *Bot) handleModalSubmit(s *discordgo.Session, i *discordgo.InteractionCreate) {
	customID := i.ModalSubmitData().CustomID
	switch customID {
	case "poll_modal":
		bot.handlePollModalSubmit(s, i)
	default:
		log.Printf("Unknown modal submission received: %s", customID)
	}
}
