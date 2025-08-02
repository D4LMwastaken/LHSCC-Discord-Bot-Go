package scripts

import (
	"errors"
	"log"
	"os"
)

func checkFileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !errors.Is(err, os.ErrNotExist)
}

func CreateFiles() {
	var envFilePath = ".env"
	var jsonFilePath = "./files/messages.json"

	isEnvFileExist := checkFileExists(envFilePath)
	isJsonFileExist := checkFileExists(jsonFilePath)

	if isEnvFileExist {
		log.Println(".env exists")
	} else {
		log.Println(".env does not exists, creating it and writing it")
		envFile, err := os.Create(envFilePath)
		if err != nil {
			log.Panicf("Error creating env file: %v", err)
		}
		_, err = envFile.WriteString("GUILD_ID=\nBOT_TOKEN=\nGEMINI_API_KEY=\nWELCOME_CHANNEL_ID=")
		if err != nil {
			log.Panicf("Error creating env file: %v", err)
		}
		err = envFile.Close()
		if err != nil {
			log.Panicf("Error closing env file: %v", err)
		}
	}

	if isJsonFileExist {
		log.Println("messages.json exist")
	} else {
		log.Println("messages.json does not exists creating it")
		jsonFile, err := os.Create(jsonFilePath)
		if err != nil {
			log.Panicf("Error while creating messages.json %v", err)
		}
		err = jsonFile.Close()
		if err != nil {
			log.Panicf("Error while closing messages.json %v", err)
		}
	}

}
