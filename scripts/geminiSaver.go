package scripts

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

var filename = "./files/messages.json"

type Message struct {
	Username  string `json:"username"`
	UserID    string `json:"userid"`
	Prompt    string `json:"prompt"`
	Response  string `json:"response"`
	Image     []byte `json:"image"`
	Timestamp string `json:"timestamp"`
}

func geminiSaver(prompt string, response string, UserID string, author string, image []byte) {
	newMessage := Message{
		Username:  author,
		UserID:    UserID,
		Prompt:    prompt,
		Response:  response,
		Image:     image,
		Timestamp: time.Now().Local().String(), // Default is hopefully eastern time, but should not matter as server is on eastern time
	}
	err := appendJSON(newMessage, filename)
	if err != nil {
		log.Printf("Error appending message to file: %v", err) // Not panic to avoid slowing down
	}
}

func writeJSON(messages []Message, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}(file)

	// Create a JSON encoder that will write to the file.
	encoder := json.NewEncoder(file)
	// Set indenting for a human-readable format.
	encoder.SetIndent("", "  ")
	// Encode the slice of messages and write it to the file.
	return encoder.Encode(messages)
}

func appendJSON(newMessage Message, filename string) error {
	messages, err := readJSON(filename)
	if err != nil {
		if os.IsNotExist(err) {
			messages = []Message{}
		} else {
			return err
		}
	}

	messages = append(messages, newMessage)

	return writeJSON(messages, filename)
}

func readJSON(filename string) ([]Message, error) {
	// Open the file.
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}(file)

	// Define a slice to hold the decoded messages.
	var messages []Message
	// Create a JSON decoder that will read from the file.
	decoder := json.NewDecoder(file)
	// Decode the file content into the messages slice.
	if err := decoder.Decode(&messages); err != nil {
		return nil, err
	}

	return messages, nil
}
