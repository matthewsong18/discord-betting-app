package main

import (
	"betting-discord-bot/internal/storage"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

// Bot parameters
var (
	GuildID       = os.Getenv("GUILD_ID")
	BotToken      = os.Getenv("BOT_TOKEN")
	AppID         = os.Getenv("APP_ID")
	DbPath        = os.Getenv("DB_PATH")
	EncryptionKey = os.Getenv("ENCRYPTION_KEY")
)

var discordSession *discordgo.Session

func init() {
	var err error
	discordSession, err = discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

func run() (err error) {
	flagErr := false
	if GuildID == "" {
		log.Println("GUILD_ID environment variable is not set")
		flagErr = true
	}
	if BotToken == "" {
		log.Println("BOT_TOKEN environment variable is not set")
		flagErr = true
	}
	if AppID == "" {
		log.Println("APP_ID environment variable is not set")
		flagErr = true
	}
	if DbPath == "" {
		log.Println("DB_PATH environment variable is not set")
		flagErr = true
	}

	if flagErr {
		return fmt.Errorf("invalid environment variables")
	}

	if EncryptionKey == "" {
		log.Println("ENCRYPTION_KEY environment variable is not set, using unencrypted database")
	}

	db, initDbError := storage.InitializeDatabase(DbPath, EncryptionKey)

	if initDbError != nil {
		return fmt.Errorf("failed to initialize database: %w", initDbError)
	}

	log.Println("Database initialized successfully")

	// Initialize internal services here

	// Run discord logic here
	discordSession.AddHandler(func(discordSession *discordgo.Session, ready *discordgo.Ready) {
		log.Println("Bot is up")
	})

	err = discordSession.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer func(discordSession *discordgo.Session) {
		_ = discordSession.Close()
	}(discordSession)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Graceful shutdown")

	defer func() {
		if closeError := db.Close(); closeError != nil {
			fmt.Println("Error closing database", closeError)
			if err == nil {
				err = closeError
			}
		}
	}()

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("application failed to start: %v", err)
	}
}
