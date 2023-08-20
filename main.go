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
	}
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
					Name:        "query",
					Description: "The query to search for.",
					Required:    true,
				},
			},
		},
	}

	for _, guild := range event.Guilds {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, guild.ID, commands[0])
		if err != nil {
			fmt.Println("error creating command,", err)
			continue // Continue to the next guild
		}

		_, err = s.ApplicationCommandCreate(s.State.User.ID, guild.ID, commands[1])
		if err != nil {
			fmt.Println("error creating command,", err)
			continue // Continue to the next guild
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
			// Check if the command has options
			if len(i.ApplicationCommandData().Options) > 0 {
				// Loop through the options and handle them
				for _, option := range i.ApplicationCommandData().Options {
					switch option.Name {
					case "query":
						value := option.Value.(string)
						response := fmt.Sprintf("You provided the query: %s", value)
						s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: response,
							},
						})
					}
				}
			}
		}
	}
}
