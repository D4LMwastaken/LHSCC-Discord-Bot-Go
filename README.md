# LHSCC-Discord-Bot-Go

## About the project
Rewrite of the original Largo High School Coding Club in the Discord Server with Go and with new features.

### Features (not yet all completed)
* Google Gemini
  * Specific profiles (probably json so results personalized)
  * AI image generation (if free)
* Faster startup and ability to send executable
* Full Google Calendar sync
* Auto update

## Project Structure

```
.
├── examples # Examples of code that I used to code this
    └── slash_commands.go # From https://github.com/bwmarrin/discordgo/
├── go.mod # All dependencies needed to run this project
    ├── go.sum # Links to all dependencies
├── LICENSE # Current License
├── main.go # Main script
├── README.md
└── scripts # Subscripts
    ├── gemini.go # Used to run functions related to Google Gemini
    └── stringspliter.go # Splits strings specifically for Google Gemini
```