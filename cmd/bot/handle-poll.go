package main

import (
	"betting-discord-bot/internal/polls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strconv"
	"strings"
)

// handleCreatePollCommand responds to the `/create-poll` command by showing a modal.
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
	pollTitle := TextDisplay{
		Type:    10,
		Content: fmt.Sprintf("# %s\n-# Warning: You cannot change your bet after submission.", title),
	}

	option1Button := Button{
		Type:     2,
		Style:    2,
		Label:    fmt.Sprintf("Bet on %s", option1),
		CustomID: fmt.Sprintf("bet:%s:0", poll.ID),
	}

	option2Button := Button{
		Type:     2,
		Style:    2,
		Label:    fmt.Sprintf("Bet on %s", option2),
		CustomID: fmt.Sprintf("bet:%s:1", poll.ID),
	}

	endPollButton := Button{
		Type:     2,
		Style:    4,
		Label:    "End Poll",
		CustomID: fmt.Sprintf("bet:%s:2", poll.ID),
	}

	selectOutcomeButton := Button{
		Type:     2,
		Style:    4,
		Label:    "Select Outcome",
		CustomID: fmt.Sprintf("bet:%s:3", poll.ID),
	}

	buttons := ActionRow{
		Type: 1,
		Components: []interface{}{
			option1Button,
			option2Button,
			endPollButton,
			selectOutcomeButton,
		},
	}

	container := Container{
		Type:        17,
		AccentColor: 0xe32458,
		Components: []interface{}{
			pollTitle,
			buttons,
		},
	}

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

func handleEndPoll(s *discordgo.Session, i *discordgo.InteractionCreate, bot *Bot, pollID string) {
	if doesNotHaveManageMemberPerm(s, i) {
		return
	}

	if err := bot.PollService.ClosePoll(pollID); err != nil {
		if errors.Is(err, polls.ErrPollIsAlreadyClosed) {
			log.Printf("Poll \"%s\" is already closed", pollID)

			sendInteractionResponse(s, i, "The poll is already closed")

			return
		}
		log.Printf("Error closing poll: %v", err)
		return
	}

	sendInteractionResponse(s, i, "The poll is closed")

	log.Printf("User %s ended poll %s", i.Member.User.GlobalName, pollID)
}

func doesNotHaveManageMemberPerm(s *discordgo.Session, i *discordgo.InteractionCreate) bool {
	if (i.Member.Permissions & discordgo.PermissionManageMessages) != discordgo.PermissionManageMessages {
		log.Printf("User \"%s\" does not have permission to end polls", i.Member.User.GlobalName)
		sendInteractionResponse(s, i, "You do not have permission to edit polls")
		return true
	}
	return false
}

func sendInteractionResponse(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	// Empty response
	data := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	}

	// Message response
	if message != "" {
		data = &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: message,
			},
		}
	}

	if err := s.InteractionRespond(i.Interaction, data); err != nil {
		log.Printf("Error sending interaction response: %v", err)
	}
}

func (bot *Bot) handleSelectOutcomeButton(s *discordgo.Session, i *discordgo.InteractionCreate, pollID string) {
	if doesNotHaveManageMemberPerm(s, i) {
		return
	}

	poll, pollErr := bot.PollService.GetPollById(pollID)
	if pollErr != nil {
		log.Printf("Error getting poll: %v", pollErr)
		return
	}

	if poll.Status == polls.Open {
		sendInteractionResponse(s, i, "The poll is still open. You cannot select an outcome.")
	}

	textDisplay := TextDisplay{
		Type:    10,
		Content: "Choose the outcome of the poll.",
	}

	selectOutcomeDropdown := &StringSelect{
		Type:        3,
		CustomID:    fmt.Sprintf("select:%s", pollID),
		Placeholder: "Select An Outcome",
		MinValues:   1,
		MaxValues:   1,
		Options: []interface{}{
			&StringOption{
				Label:       poll.Options[0],
				Value:       "1",
				Description: "Option 1",
			},
			&StringOption{
				Label:       poll.Options[1],
				Value:       "2",
				Description: "Option 2",
			},
		},
	}

	actionRow := ActionRow{
		Type: 1,
		Components: []interface{}{
			selectOutcomeDropdown,
		},
	}

	messageContainer := Container{
		Type:        17,
		AccentColor: 0xe32458,
		Components: []interface{}{
			textDisplay,
			actionRow,
		},
	}

	const permissions = IsComponentsV2 | MessageIsEphemeral

	message := MessageSend{
		Flags: permissions,
		Components: []interface{}{
			messageContainer,
		},
	}

	response := InteractionResponse{
		Type: 4,
		Data: message,
	}

	jsonMessage, jsonErr := json.Marshal(response)
	if jsonErr != nil {
		log.Printf("Error marshaling selectOutcomeDropdown: %v", jsonErr)
		return
	}

	interactionID := i.ID
	interactionToken := i.Token
	url := fmt.Sprintf("https://discord.com/api/v10/interactions/%s/%s/callback", interactionID, interactionToken)
	sendHttpRequest(url, jsonMessage)
}

func (bot *Bot) handleSelectOutcomeDropdown(s *discordgo.Session, i *discordgo.InteractionCreate) {
	customID := i.MessageComponentData().CustomID
	messageData := strings.Split(customID, ":")
	pollID := messageData[1]
	optionIndex, err := strconv.Atoi(i.MessageComponentData().Values[0])
	if err != nil {
		log.Printf("Error parsing option index: %v", err)
		return
	}

	if _, pollErr := bot.PollService.GetPollById(pollID); pollErr != nil {
		log.Printf("Error getting poll: %v", pollErr)
		return
	}

	if err := bot.PollService.SelectOutcome(pollID, polls.OutcomeStatus(optionIndex)); err != nil {
		log.Printf("Error selecting outcome: %v", err)
		return
	}

	log.Printf("Outcome has been selected for poll %s", pollID)
	sendInteractionResponse(s, i, "The outcome of the poll has been selected.")
}
