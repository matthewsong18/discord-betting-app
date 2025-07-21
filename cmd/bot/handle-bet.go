package main

import (
	"betting-discord-bot/internal/bets"
	"betting-discord-bot/internal/users"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
)

func handleBet(s *discordgo.Session, i *discordgo.InteractionCreate, bot *Bot, pollID string, user *users.User, optionIndex int) {
	bet, betErr := bot.BetService.CreateBet(pollID, user.ID, optionIndex)
	if betErr != nil {
		if errors.Is(betErr, bets.ErrUserAlreadyBet) {
			if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprint("You have already bet on this poll."),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			}); err != nil {
				log.Printf("Error sending invalid action response: %v", err)
			}
		}

		if errors.Is(betErr, bets.ErrPollIsClosed) {
			if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprint("This poll is closed. You cannot place a bet."),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			}); err != nil {
				log.Printf("Error sending poll is closed response: %v", err)
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
