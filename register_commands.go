package main

import (
	"discord-go-test/commands/dadjoke"
	"discord-go-test/commands/takefruit"

	"github.com/bwmarrin/discordgo"
)

func registerCommands(s *discordgo.Session, event *discordgo.Ready) {
	commands := []*discordgo.ApplicationCommand{
		&dadjoke.CommandDefinition,
		&takefruit.CommandDefinition,
	}
	for _, command := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", command)
		if err != nil {
			panic(err)
		}
	}
}

func interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		if i.ApplicationCommandData().Name == "dadjoke" {
			dadjoke.Command(s, i)
		}
		if i.ApplicationCommandData().Name == "take_fruit" {
			takefruit.Command(s, i)
		}
	}
}
