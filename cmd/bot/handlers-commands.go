package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

func (bot *Bot) handleCreatePollCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
