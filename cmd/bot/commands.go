package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
)

func (bot *Bot) RegisterCommands() error {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "create-poll",
			Description: "Create a new poll",
		},
		{
			Name:        "bet",
			Description: "Place a bet on an open poll",
		},
	}

	AppID = os.Getenv("APP_ID")
	GuildID = os.Getenv("GUILD_ID")
	_, err := bot.DiscordSession.ApplicationCommandBulkOverwrite(AppID, GuildID, commands)
	if err != nil {
		log.Printf("Error overwriting commands: %v", err)
		return err
	}

	log.Println("Commands successfully registered.")
	return nil
}

func pollModal() func(*discordgo.Session, *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID: "poll_modal",
				Title:    "Create Poll",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:    "title",
								Label:       "What is the title of the poll?",
								Style:       discordgo.TextInputShort,
								Placeholder: "Who wins between A or B blah blah blah",
								Required:    true,
								MaxLength:   300,
								MinLength:   10,
							},
						},
					},
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:  "option-1",
								Label:     "What's the first option?",
								Style:     discordgo.TextInputParagraph,
								Required:  true,
								MaxLength: 2000,
							},
						},
					},
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:  "option-2",
								Label:     "What's the second option?",
								Style:     discordgo.TextInputParagraph,
								Required:  true,
								MaxLength: 2000,
							},
						},
					},
				},
			},
		})

		if err != nil {
			panic(err)
		}
	}
}
