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
			Name:        "file-send-test",
			Description: "Checks if Discord Bot can send files on the Discord Server",
		},
		{
			Name:        "google-gemini-test",
			Description: "Checks if Discord Bot can receive messages from Google Gemini API",
		},
		{
			Name:        "google-gemini-test-followup",
			Description: "Requirement so Gemini API can function due to time constraints",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ping": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Pong!",
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

		"file-send-test": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Here is a file",
					Files: []*discordgo.File{
						{
							ContentType: "text/plain",
							Name:        "test.txt",
							Reader:      strings.NewReader(scripts.GeminiAI("Create random code as fast as you can, it is going to be in a file.")),
						},
					},
				},
			})

			if err != nil {
				print("Unable to send message with attachment Error: ", err)
			}
		},

		"google-gemini-test": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
			content = scripts.GeminiAI("What is the hackclub in less than 2000 characters?")
			_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &content,
			})
			if err != nil {
				print("Unable to send Gemini message Error: ", err)
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
