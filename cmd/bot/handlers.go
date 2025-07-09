package main

import (
	"fmt"
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

// handleCreatePollCommand responds to the `/create-poll` command by showing a modal.
func (bot *Bot) handleCreatePollCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "poll_modal", // The ID we'll check for on submission
			Title:    "Create a New Poll",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "title",
							Label:       "Poll Title",
							Style:       discordgo.TextInputShort,
							Placeholder: "Who will win the grand finals?",
							Required:    true,
							MaxLength:   300,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:  "option1",
							Label:     "First Option",
							Style:     discordgo.TextInputShort,
							Required:  true,
							MaxLength: 100,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:  "option2",
							Label:     "Second Option",
							Style:     discordgo.TextInputShort,
							Required:  true,
							MaxLength: 100,
						},
					},
				},
			},
		},
	})

	if err != nil {
		log.Printf("Error showing modal: %v", err)
	}
}

// handlePollModalSubmit processes the data from the poll creation modal.
func (bot *Bot) handlePollModalSubmit(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ModalSubmitData()

	// Safely parse the data from the modal.
	var title, option1, option2 string
	for _, row := range data.Components {
		input := row.(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput)
		switch input.CustomID {
		case "title":
			title = input.Value
		case "option1":
			option1 = input.Value
		case "option2":
			option2 = input.Value
		}
	}

	log.Printf("Poll submitted: Title='%s', Option1='%s', Option2='%s'", title, option1, option2)

	_, err := bot.PollService.CreatePoll(title, []string{option1, option2})
	if err != nil {
		log.Printf("Error creating poll: %v", err)
		return
	}

	// Send a confirmation message back to the user who submitted the modal.
	// This message is "ephemeral," so only they can see it.
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprint("Poll submitted"),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	}); err != nil {
		log.Printf("Error sending modal confirmation: %v", err)
	}
}
