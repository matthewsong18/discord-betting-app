package main

import (
	"betting-discord-bot/internal/bets"
	"betting-discord-bot/internal/users"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strconv"
	"strings"
)

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
