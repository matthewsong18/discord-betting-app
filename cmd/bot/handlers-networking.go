package main

import (
	"bytes"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"log"
	"net/http"
	"os"
)

func sendHttpRequest(url string, jsonMessage []byte) {
	request, requestErr := http.NewRequest("POST", url, bytes.NewBuffer(jsonMessage))
	if requestErr != nil {
		log.Printf("error creating request: %v", requestErr)
		return
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

	log.Println("Successfully sent http request to Discord.")
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
