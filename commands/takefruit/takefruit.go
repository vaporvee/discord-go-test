package takefruit

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var CommandDefinition discordgo.ApplicationCommand = discordgo.ApplicationCommand{
	Name:        "take_fruit",
	Description: "Take a fruit.",
}

func Command(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
