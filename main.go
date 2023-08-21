package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	fmt.Printf("Bot is now running as \"%s\"!", discord.State.User.Username)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	fmt.Println("\nShutting down...")
	defer removeCommandsFromAllGuilds(discord)
	discord.Close()
}

func removeCommandsFromAllGuilds(s *discordgo.Session) {
	for _, guild := range s.State.Guilds {
		existingCommands, err := s.ApplicationCommands(s.State.User.ID, guild.ID)
		if err != nil {
			fmt.Printf("error fetching existing commands for guild %s: %v\n", guild.Name, err)
			continue
		}

		for _, existingCommand := range existingCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, guild.ID, existingCommand.ID)
			if err != nil {
				fmt.Printf("error deleting command %s for guild %s: %v\n", existingCommand.Name, guild.Name, err)
			}
		}
	}
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "dadjoke",
			Description: "Get a dad joke.",
		},
		{
			Name:        "take_fruit",
			Description: "A test command with auto-completion.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "fruit",
					Description: "The fruit you are taking.",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Apple",
							Value: "apple",
						},
						{
							Name:  "Banana",
							Value: "banana",
						},
						{
							Name:  "Cherry",
							Value: "cherry",
						},
					},
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
		if i.ApplicationCommandData().Name == "dadjoke" {
			joke, err := getDadJoke()
			if err != nil {
				fmt.Printf("error fetching dad joke: %v\n", err)
				return
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: joke,
				},
			})
		}
		if i.ApplicationCommandData().Name == "take_fruit" {
			if len(i.ApplicationCommandData().Options) > 0 {
				for _, option := range i.ApplicationCommandData().Options {
					if option.Name == "fruit" {
						value := option.Value.(string)
						var response string
						response = fmt.Sprintf("You took a %s!", value)
						channel, _ := s.UserChannelCreate(i.Member.User.ID)
						s.ChannelMessageSend(channel.ID, response)
						s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "DM Sent!",
								Flags:   1 << 6, // ephemeral message
							},
						})
					}
				}
			}
		}
	}
}

func getDadJoke() (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://icanhazdadjoke.com/", nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var jokeResponse struct {
		Joke string `json:"joke"`
	}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&jokeResponse); err != nil {
		return "", err
	}
	return jokeResponse.Joke, nil
}
