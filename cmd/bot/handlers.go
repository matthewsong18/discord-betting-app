package main

import (
	"betting-discord-bot/internal/bets"
	"betting-discord-bot/internal/polls"
	"betting-discord-bot/internal/users"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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

func (bot *Bot) handleButtonPress(s *discordgo.Session, i *discordgo.InteractionCreate) {
	customID := i.MessageComponentData().CustomID
	betData := strings.Split(customID, ":")

	if len(betData) != 3 {
		log.Printf("Invalid custom ID received: %s", customID)
		return
	}

	pollID := betData[1]
	optionIndex, err := strconv.Atoi(betData[2])
	if err != nil {
		log.Printf("failed to convert option index to int: %s", betData[2])
		return
	}

	userDiscordID := i.Member.User.ID

	user, getUserErr := bot.UserService.GetUserByDiscordID(userDiscordID)
	if getUserErr != nil {
		if errors.Is(getUserErr, users.ErrUserNotFound) {
			var createUserErr error
			user, createUserErr = bot.UserService.CreateUser(userDiscordID)
			if createUserErr != nil {
				log.Printf("Error creating user: %v", createUserErr)
				return
			}
		} else {
			log.Printf("Error getting user: %v", getUserErr)
			return
		}
	}

	bet, betErr := bot.BetService.CreateBet(pollID, user.ID, optionIndex)
	if betErr != nil {
		if errors.Is(betErr, bets.ErrUserAlreadyBet) {
			if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprint("You have already bet on this poll"),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			}); err != nil {
				log.Printf("Error sending bet confirmation: %v", err)
			}
		}

		log.Printf("Error creating bet: %v", betErr)
		return
	}

	log.Printf("Bet created: %v", bet)

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprint("Bet submitted"),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	}); err != nil {
		log.Printf("Error sending bet confirmation: %v", err)
	}

}

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
		newTextInputRow(
			"title",
			"Poll Title",
			"Who will win the grand finals?",
			discordgo.TextInputShort,
			300,
			true,
		),
		newTextInputRow(
			"option1",
			"First Option",
			"",
			discordgo.TextInputShort,
			100,
			true,
		),
		newTextInputRow(
			"option2",
			"Second Option",
			"",
			discordgo.TextInputShort,
			100,
			true,
		),
	}
}

func newTextInputRow(customID, label, placeholder string, style discordgo.TextInputStyle, maxLength int, required bool) discordgo.ActionsRow {
	return discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			&discordgo.TextInput{
				CustomID:    customID,
				Label:       label,
				Placeholder: placeholder,
				Style:       style,
				Required:    required,
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

	sendPollMessage(title, option1, option2, poll, i)
}

func sendPollMessage(title string, option1 string, option2 string, poll *polls.Poll, i *discordgo.InteractionCreate) {
	pollTitle := TextDisplay{
		Type:    10,
		Content: fmt.Sprintf("# %s\n", title),
	}

	button1 := Button{
		Type:     2,
		Style:    2,
		Label:    fmt.Sprintf("Bet on %s", option1),
		CustomID: fmt.Sprintf("bet:%s:0", poll.ID),
	}

	button2 := Button{
		Type:     2,
		Style:    2,
		Label:    fmt.Sprintf("Bet on %s", option2),
		CustomID: fmt.Sprintf("bet:%s:1", poll.ID),
	}

	buttons := ActionRow{
		Type: 1,
		Components: []interface{}{
			button1,
			button2,
		},
	}

	message := MessageSend{
		Flags: IsComponentsV2,
		Components: []interface{}{
			pollTitle,
			buttons,
		},
	}

	jsonMessage, jsonErr := json.Marshal(message)
	if jsonErr != nil {
		log.Printf("Error marshaling message: %v", jsonErr)
		return
	}

	channelID := i.ChannelID
	url := fmt.Sprintf("https://discord.com/api/v10/channels/%s/messages", channelID)
	request, requestErr := http.NewRequest("POST", url, bytes.NewBuffer(jsonMessage))
	if requestErr != nil {
		//return fmt.Errorf("error creating request: %w", requestErr)
	}

	request.Header.Set("Content-Type", "application/json")
	botToken := os.Getenv("TOKEN")
	request.Header.Set("Authorization", fmt.Sprintf("Bot %s", botToken))

	log.Println("Sending manual HTTP request to Discord API...")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Printf("error sending HTTP request to Discord: %v", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		// If it's not a success, we read the error message Discord sent back.
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("discord API returned a non-success status code %d: %s", resp.StatusCode, string(bodyBytes))
		return
	}

	log.Println("Successfully sent message with custom components.")
}
