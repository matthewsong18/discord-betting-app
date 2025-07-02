package main

import (
	"betting-discord-bot/internal/bets"
	"betting-discord-bot/internal/polls"
	"betting-discord-bot/internal/users"
	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	DiscordSession *discordgo.Session
	PollService    polls.PollService
	BetService     bets.BetService
	UserService    users.UserService
}

func NewBot(session *discordgo.Session, pollService polls.PollService, betService bets.BetService, userService users.UserService) *Bot {
	return &Bot{
		DiscordSession: session,
		PollService:    pollService,
		BetService:     betService,
		UserService:    userService,
	}
}
