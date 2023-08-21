package dadjoke

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

var CommandDefinition discordgo.ApplicationCommand = discordgo.ApplicationCommand{
	Name:        "dadjoke",
	Description: "Get a dad joke.",
}

func Command(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
