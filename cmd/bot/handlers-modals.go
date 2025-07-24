package main

import (
	"betting-discord-bot/internal/polls"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
)

func pollModalComponents() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		newTextInputRow("title", "Poll Title", "Who will win the grand finals?", 50),
		newTextInputRow("option1", "First Option", "", 20),
		newTextInputRow("option2", "Second Option", "", 20),
	}
}

func newTextInputRow(customID, label, placeholder string, maxLength int) discordgo.ActionsRow {
	return discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			&discordgo.TextInput{
				CustomID:    customID,
				Label:       label,
				Placeholder: placeholder,
				Style:       discordgo.TextInputShort,
				Required:    true,
				MaxLength:   maxLength,
			},
		},
	}
}

// handlePollModalSubmit processes the data from the poll creation modal.
func (bot *Bot) handlePollModalSubmit(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Println("A user submitted a poll modal")

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

	poll, err := bot.PollService.CreatePoll(title, []string{option1, option2})
	if err != nil {
		log.Printf("Error creating poll: %v", err)
		return
	}

	responseMessage := "Reminder: You must end the poll before you are allowed to select an outcome."
	sendInteractionResponse(s, i, responseMessage)

	sendPollMessage(title, option1, option2, poll, i)
}

func sendPollMessage(title string, option1 string, option2 string, poll *polls.Poll, i *discordgo.InteractionCreate) {
	pollString := fmt.Sprintf("# %s\n-# Warning: You cannot change your bet after submission.", title)
	pollTitle := NewTextDisplay(pollString)

	option1Button := NewButton(
		2,
		fmt.Sprintf("Bet on %s", option1),
		fmt.Sprintf("bet:%s:0", poll.ID),
	)

	option2Button := NewButton(
		2,
		fmt.Sprintf("Bet on %s", option2),
		fmt.Sprintf("bet:%s:1", poll.ID),
	)

	endPollButton := NewButton(
		4,
		"End Poll",
		fmt.Sprintf("bet:%s:2", poll.ID),
	)

	selectOutcomeButton := NewButton(
		4,
		"Select Outcome",
		fmt.Sprintf("bet:%s:3", poll.ID),
	)

	buttons := NewActionRow(
		[]interface{}{
			option1Button,
			option2Button,
			endPollButton,
			selectOutcomeButton,
		},
	)

	container := NewContainer(
		0xe32458,
		[]interface{}{
			pollTitle,
			buttons,
		},
	)

	message := MessageSend{
		Flags: IsComponentsV2,
		Components: []interface{}{
			container,
		},
	}

	jsonMessage, jsonErr := json.Marshal(message)
	if jsonErr != nil {
		log.Printf("Error marshaling message: %v", jsonErr)
		return
	}

	channelID := i.ChannelID
	url := fmt.Sprintf("https://discord.com/api/v10/channels/%s/messages", channelID)
	sendHttpRequest(url, jsonMessage)
}
