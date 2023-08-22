package TEMPLATE

import (
	"github.com/bwmarrin/discordgo"
)

var CommandDefinition discordgo.ApplicationCommand = discordgo.ApplicationCommand{
	Name:        "TEMPLATE",
	Description: "This is a slash command.",
}

func Command(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Answer!",
		},
	})
}
