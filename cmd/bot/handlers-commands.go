package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func (bot *Bot) handleCreatePollCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Println("A user requested to create a poll")

	modalData := &discordgo.InteractionResponseData{
		CustomID:   "poll_modal", // The ID we'll check for on submission
		Title:      "Create a New Poll",
		Components: pollModalComponents(),
	}

	if err := s.InteractionRespond(
		i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: modalData,
		},
	); err != nil {
		log.Printf("Error showing modal: %v", err)
	}
}
