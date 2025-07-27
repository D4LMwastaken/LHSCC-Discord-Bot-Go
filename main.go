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
}

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "Sends out pong when running this slash command",
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
			author := i.Member.User.DisplayName()
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Pong! " + author + " ",
				},
			})

			if err != nil {
				print("Unable to send message pong, Error: ", err)
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
					print("Unable to send message permissions, Error: ", err)
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

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Please wait for Gemini API to process the file",
				},
			})
			if err != nil {
				print("Unable to send message: ", err)
			}
			content, _ := scripts.GeminiAI(task, DisplayName, false, "file", chosenLanguage)
			_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
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
			// Username := i.Member.User.Username // Will be used for other purposes later

			content := "Please wait for Gemini API to process the text you just sent..."
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
				},
			})
			if err != nil {
				print("Unable to send Gemini message Error: ", err)
			}
			first, rest := scripts.GeminiAI(question, DisplayName, true, "ask", "none")
			_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &first,
			})
			if err != nil {
				print("Unable to send Gemini message Error: ", err)
			}
			for a := range rest {
				_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: rest[a],
				})
				if err != nil {
					print("Unable to send Gemini message Error: ", err)
				}
				a++
			}
		},
	}
)

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {
	GuildID := os.Getenv("GUILD_ID")
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	err := s.Open()

	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	// Remove all existing commands at startup

	log.Println("Removing existing commands...")

	existingCommands, err := s.ApplicationCommands(s.State.User.ID, GuildID)

	if err != nil {
		log.Printf("Could not fetch existing commands: %v", err)
	} else {
		for _, cmd := range existingCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, GuildID, cmd.ID)
			if err != nil {
				log.Printf("Could not delete command '%v': %v", cmd.Name, err)
			} else {
				log.Printf("Deleted command: %v", cmd.Name)
			}
		}
	}

	log.Println("Adding commands...")

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))

	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer func(s *discordgo.Session) {
		err := s.Close()
		if err != nil {
			log.Printf("Error closing session: %v", err)
		}
	}(s)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Removing commands...")
	for _, v := range registeredCommands {
		err := s.ApplicationCommandDelete(s.State.User.ID, GuildID, v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
	log.Println("Gracefully shutting down.")
}
