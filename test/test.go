package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

func main() {
	var err error
	err = godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	GUILDID := os.Getenv("GUILD_ID")
	GUILDIDS := strings.Split(GUILDID, ",")
	for i := range GUILDIDS {
		println(GUILDIDS[i])
	}
}
