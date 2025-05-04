package main

import (
	"discordBot/command"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/servusdei2018/shards/v2"
)

var (
	mgr      *shards.Manager
	token    string
	Commands = make(map[string]command.Command)
)

func registerCommands(s *discordgo.Session, cmds []command.Command) {
	for _, cmd := range cmds {
		Commands[cmd.ApplicationCommand.Name] = cmd
		if _, err := s.ApplicationCommandCreate(s.State.User.ID, "", &cmd.ApplicationCommand); err != nil {
			log.Printf("Cannot create command: %v\n", err)
		}
	}
}

func init() {
	_ = godotenv.Load()
	token = os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_TOKEN not set")
	}
}

func main() {
	var err error
	if mgr, err = shards.New("Bot " + token); err != nil {
		log.Fatalf("Error creating shard manager: %v\n", err)
	}

	mgr.AddHandler(onReady)
	mgr.AddHandler(messageCreate)
	mgr.AddHandler(onInteractionCreate)

	mgr.RegisterIntent(discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent)

	if err = mgr.Start(); err != nil {
		log.Fatalf("Error starting shard manager: %v\n", err)
	}
	log.Println("Bot is running. Press Ctrl+C to exit.")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	if err = mgr.Shutdown(); err != nil {
		log.Printf("Error during shutdown: %v\n", err)
	} else {
		log.Println("Bot shut down cleanly.")
	}
}

func onReady(s *discordgo.Session, _ *discordgo.Ready) {
	registerCommands(s, []command.Command{
		command.Ping,
		command.Echo,
	})
}
