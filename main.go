package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	discord, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	} else {
		fmt.Println("Discord session created")
	}
	discord.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds
	discord.AddHandler(ready)
	discord.AddHandler(interactionCreate)

	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	discord.Close()
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "test",
			Description: "A test command.",
		},
		{
			Name:        "secondtest",
			Description: "A second test command.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "other_testing_value",
					Description: "Value that the bot repeats back to you.",
					Required:    true,
				},
			},
		},
	}

	for _, guild := range event.Guilds {
		for _, command := range commands {
			_, err := s.ApplicationCommandCreate(s.State.User.ID, guild.ID, command)
			if err != nil {
				fmt.Printf("error creating command for %s: %v\n", guild.Name, err)
			}
		}
	}
}

func interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		if i.ApplicationCommandData().Name == "test" {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "You tested me!",
				},
			})
		}
		if i.ApplicationCommandData().Name == "secondtest" {
			if len(i.ApplicationCommandData().Options) > 0 {
				for _, option := range i.ApplicationCommandData().Options {
					if option.Name == "other_testing_value" {
						value := option.Value.(string)
						s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: fmt.Sprintf("You provided the testing value: %s", value),
							},
						})
					}
				}
			}
		}
	}
}
