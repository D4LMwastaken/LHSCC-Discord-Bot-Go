package main

import (
	"github.com/bwmarrin/discordgo"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "Sends out pong when running this slash command along with latency",
		},
		{
			Name:        "ping-gemini",
			Description: "Pings Google Gemini API and see how long it takes for a specific model to respond",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "model",
					Description: "Name of the model Gemini API will use",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "prompt",
					Description: "Prompt that you want Google Gemini to follow and measure latency for",
					Required:    true,
				},
			},
		},
		{
			Name:        "hi",
			Description: "Sends out a friendly greeting powered by Google Gemini",
		},
		{
			Name: "help",
			Description: "Sends out help message with all the possible commands powered in a message powered by Google " +
				"Gemini",
		},
		{
			Name:        "bye",
			Description: "Sends out a goodbye message powered by Google Gemini and shut down if user is a specific person",
		},
		{
			Name:        "version",
			Description: "Sends out a version information on what version the Discord Bot is on",
		},
		{
			Name:        "new-stuff",
			Description: "Sends out all the new additions to the Discord Bot added, powered by Google Gemini",
		},
		{
			Name:        "permission-check",
			Description: "Checks if Discord Bot has enough permissions",
		},
		{
			Name:        "code",
			Description: "Checks if Discord Bot can send files on the Discord Server",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "chosen_language",
					Description: "Chosen programming that Google Gemini will code in",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "task",
					Description: "Task you want Google Gemini to code",
					Required:    true,
				},
			},
		},
		{
			Name:        "ask",
			Description: "Ask Google Gemini something",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "question",
					Description: "Question you want to ask Google Gemini",
					Required:    true,
				},
			},
		},
	}
)

func main() {
	commandNames := ""
	commandDescriptions := ""
	for a := range commands {
		commandNames = commandNames + commands[a].Name + ","
		commandDescriptions = commandDescriptions + commands[a].Description + ","
	}
	println(commandNames)
	println(commandDescriptions)
}
