package main

import (
	"LHSCC-Discord-Bot/main/scripts"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"strings"
)

var s *discordgo.Session

func init() {
	scripts.CreateFiles()
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	BotToken := os.Getenv("BOT_TOKEN")
	s, err = discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
	s.Identify.Intents = discordgo.IntentsGuildMembers | discordgo.IntentsGuilds | discordgo.IntentsGuildMessages
	s.AddHandler(GuildMemberAdd)
}

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

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ping": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Latency is: " + s.HeartbeatLatency().String(),
				},
			})
			if err != nil {
				log.Panic("Unable to send message pong, Error: ", err)
			}
		},

		"ping-gemini": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}
			prompt := optionMap["prompt"].StringValue()
			model := optionMap["model"].StringValue()
			author := i.Member.User.DisplayName()
			userID := i.Member.User.ID
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Please wait for Gemini AI to respond",
				},
			})
			if err != nil {
				log.Panic("Unable to send message, Error: ", err)
			}

			first, rest, latency := scripts.PingGemini(prompt, author, userID, true, "ask", "none", model)
			_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &first,
			})
			if err != nil {
				log.Panic("Unable to send message from Gemini, Error: ", err)
			}
			for a := range rest {
				_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: rest[a],
				})
				if err != nil {
					log.Panic("Unable to send message from Gemini, Error: ", err)
				}
			}
			_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Content: "Latency: " + latency,
			})
		},

		"hi": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			author := i.Member.DisplayName()
			userID := i.Member.User.ID
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Please wait for Gemini AI to respond",
				},
			})
			if err != nil {
				log.Panic("Unable to send message, Error: ", err)
			}
			prompt := "Send out a message to the author giving them a unique greeting that is less than 2000 characters"
			content, _ := scripts.GeminiAI(prompt, author, userID, false, "ask", "none", "lite")
			_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &content,
			})
		},

		"help": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			author := i.Member.DisplayName()
			userID := i.Member.User.ID
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Please wait for Gemini AI to respond",
				},
			})
			if err != nil {
				log.Panic("Unable to send message, Error: ", err)
			}
			commandNames := ""
			commandDescriptions := ""
			for i := range commands {
				commandNames += commands[i].Name + ","
				commandDescriptions += commands[i].Description + ","
			}
			firstContent, restOfContent := scripts.Help(commandNames, commandDescriptions, author, userID)
			_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &firstContent,
			})
			if err != nil {
				log.Panic("Unable to send message, Error: ", err)
			}
			for a := range restOfContent {
				_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: restOfContent[a],
				})
				if err != nil {
					log.Panic("Unable to send message from Gemini, Error: ", err)
				}
				a++
			}
		},

		"bye": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			author := i.Member.DisplayName()
			username := i.Member.User.Username
			id := i.Member.User.ID
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Please wait for Gemini AI to respond",
				},
			})
			if err != nil {
				log.Panic("Unable to send message, Error: ", err)
			}
			content, restOfContent, isDev := scripts.Bye(author, username, id)
			_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &content,
			})
			if err != nil {
				log.Panic("Unable to send message, Error: ", err)
			}
			for a := range restOfContent {
				_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: restOfContent[a],
				})
				if err != nil {
					log.Panic("Unable to send message from Gemini, Error: ", err)
				}
			}
			if isDev == true {
				err := s.Close() // Later will allow program to close fully from bye command...
				if err != nil {
					log.Panic("Unable to close Discord Websocket, Error: ", err)
				}
			}
		},

		"version": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			author := i.Member.DisplayName()
			userID := i.Member.User.ID
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Please wait for Gemini AI to respond",
				},
			})
			if err != nil {
				log.Panic("Unable to send message, Error: ", err)
			}
			content, restOfContent := scripts.Version(author, userID)
			_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &content,
			})
			if err != nil {
				log.Panic("Unable to send message, Error: ", err)
			}
			for a := range restOfContent {
				_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: restOfContent[a],
				})
				if err != nil {
					log.Panic("Unable to send message from Gemini, Error: ", err)
				}
			}
		},

		"new-stuff": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			author := i.Member.DisplayName()
			userID := i.Member.User.ID
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Please wait for Gemini AI to respond",
				},
			})
			if err != nil {
				log.Panic("Unable to send message, Error: ", err)
			}
			first, restOfContent := scripts.NewStuff(author, userID)
			_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &first,
			})
			if err != nil {
				log.Panic("Unable to send message, Error: ", err)
			}
			for a := range restOfContent {
				_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: restOfContent[a],
				})
				if err != nil {
					log.Panic("Unable to send message from Gemini, Error: ", err)
				}
			}
		},

		"permission-check": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			perms, err := s.ApplicationCommandPermissions(s.State.User.ID, i.GuildID, i.ApplicationCommandData().ID)
			var restError *discordgo.RESTError
			if errors.As(err, &restError) && restError.Message != nil && restError.Message.Code == discordgo.ErrCodeUnknownApplicationCommandPermissions {
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: ":x: No permission overwrites",
					},
				})

				if err != nil {
					log.Panic("Unable to send message permissions, Error: ", err)
				}
				return

			} else if err != nil {
				panic(err)
			}

			if err != nil {
				panic(err)
			}

			format := "- %s %s\n"
			channels := ""
			users := ""
			roles := ""

			for _, o := range perms.Permissions {
				emoji := "❌"
				if o.Permission {
					emoji = "☑"
				}

				switch o.Type {
				case discordgo.ApplicationCommandPermissionTypeUser:
					users += fmt.Sprintf(format, emoji, "<@!"+o.ID+">")

				case discordgo.ApplicationCommandPermissionTypeChannel:
					allChannels, _ := discordgo.GuildAllChannelsID(i.GuildID)

					if o.ID == allChannels {
						channels += fmt.Sprintf(format, emoji, "All channels")
					} else {
						channels += fmt.Sprintf(format, emoji, "<#"+o.ID+">")
					}

				case discordgo.ApplicationCommandPermissionTypeRole:
					if o.ID == i.GuildID {
						roles += fmt.Sprintf(format, emoji, "@everyone")
					} else {
						roles += fmt.Sprintf(format, emoji, "<@&"+o.ID+">")
					}
				}
			}

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Permissions overview",
							Description: "Overview of permissions for this command",
							Fields: []*discordgo.MessageEmbedField{
								{
									Name:  "Users",
									Value: users,
								},
								{
									Name:  "Channels",
									Value: channels,
								},
								{
									Name:  "Roles",
									Value: roles,
								},
							},
						},
					},
					AllowedMentions: &discordgo.MessageAllowedMentions{},
				},
			})

			if err != nil {
				print("Unable to send message permissions, Error: ", err)
			}
		},

		"code": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}
			chosenLanguage := optionMap["chosen_language"].StringValue()
			fileExtension, fileType := scripts.FileSupportCheck(chosenLanguage)
			task := optionMap["task"].StringValue()

			DisplayName := i.Member.DisplayName()
			userID := i.Member.User.ID

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Please wait for Gemini API to process the file",
				},
			})
			if err != nil {
				print("Unable to send message: ", err)
			}
			content, _ := scripts.GeminiAI(task, DisplayName, userID, false, "file", chosenLanguage, "pro")
			_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: nil,
				Files: []*discordgo.File{
					{
						ContentType: "text/" + fileType,
						Name:        "response" + fileExtension,
						Reader:      strings.NewReader(content),
					},
				},
			})
			if err != nil {
				print("Unable to send Gemini message Error: ", err)
			}
		},

		"ask": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}
			question := optionMap["question"].StringValue()

			DisplayName := i.Member.DisplayName()
			userID := i.Member.User.ID

			content := "Please wait for Gemini API to process the text you just sent..."
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
				},
			})
			if err != nil {
				log.Panic("Unable to send Gemini message Error: ", err)
			}
			first, rest := scripts.GeminiAI(question, DisplayName, userID, true, "ask", "none", "pro")
			_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &first,
			})
			if err != nil {
				log.Panic("Unable to send Gemini message Error: ", err)
			}
			for a := range rest {
				_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: rest[a],
				})
				if err != nil {
					log.Panic("Unable to send Gemini message Error: ", err)
				}
				a++
			}
		},
	}
)

func main() {
	GuildIDs := strings.Split(os.Getenv("GUILD_ID"), ",")
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err := s.Open()

	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	// Remove all existing commands at startup

	log.Println("Removing existing commands...")

	for i := range GuildIDs {
		existingCommands, err := s.ApplicationCommands(s.State.User.ID, GuildIDs[i])
		if err != nil {
			log.Printf("Could not fetch existing commands: %v", err)
		} else {
			for _, cmd := range existingCommands {
				err := s.ApplicationCommandDelete(s.State.User.ID, GuildIDs[i], cmd.ID)
				if err != nil {
					log.Panicf("Could not delete command '%v': %v", cmd.Name, err)
				} else {
					log.Printf("Deleted command: '%v' on %v", cmd.Name, GuildIDs[i])
				}
			}
		}
	}

	log.Println("Adding commands...")

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))

	for i, v := range commands {
		for a := range GuildIDs {
			cmd, err := s.ApplicationCommandCreate(s.State.User.ID, GuildIDs[a], v)
			if err != nil {
				log.Panicf("Cannot create '%v' command on: %v", v.Name, err)
			} else {
				log.Printf("Created command '%v' on %v", cmd.Name, GuildIDs[a])
			}
			registeredCommands[i] = cmd
		}
	}

	defer func(s *discordgo.Session) {
		err := s.Close()
		if err != nil {
			log.Panicf("Error closing session: %v", err)
		}
	}(s)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Removing commands...")
	for _, v := range registeredCommands {
		for i := range GuildIDs {
			err := s.ApplicationCommandDelete(s.State.User.ID, GuildIDs[i], v.ID)
			if err != nil {
				log.Printf("Cannot delete '%v' command on Guild '%v': %v", v.Name, GuildIDs[i], err)
			}
		}
	}

	log.Println("Gracefully shutting down.")
}

func GuildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	channelID := os.Getenv("WELCOME_CHANNEL_ID")
	println("A new member has joined!")
	displayName := m.DisplayName()
	userID := m.User.ID
	prompt := "Generate a unique Discord message on the LHSCC Discord Server greeting the new user who just joined by mentioning them who's name is: " + displayName + "\n" +
		"Then, tell the new user to check out readme-md, where the rules are and wait for a member to make them a member of the Discord Server.\n" +
		"Mention them by: " + m.User.Mention()
	firstContent, restOfContent := scripts.GeminiAI(prompt, displayName, userID, true, "ask", "none", "flash")
	_, err := s.ChannelMessageSend(channelID, firstContent)
	for a := range restOfContent {
		_, err = s.ChannelMessageSend(channelID, restOfContent[a])
	}
	if err != nil {
		log.Panic("Error sending message when message joins: ", err)
	}
	DMChannel, err := s.UserChannelCreate(m.User.ID)
	if err != nil {
		log.Panic("Error creating DM channel: ", err)
	}
	_, err = s.ChannelMessageSend(DMChannel.ID, firstContent)
	if err != nil {
		log.Panic("Error sending message to DM: ", err)
	}
	for a := range restOfContent {
		_, err = s.ChannelMessageSend(channelID, restOfContent[a])
	}
}
