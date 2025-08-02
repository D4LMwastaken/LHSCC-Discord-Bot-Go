package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

// Message struct represents the data for a single Discord bot interaction.
// The `json` struct tags are used for mapping fields during I/O.
type Message struct {
	Prompt    string    `json:"prompt"`
	Response  string    `json:"response"`
	Timestamp time.Time `json:"timestamp"`
	UserID    string    `json:"userId"`
}

func main() {
	// Define a new message to add to the file.
	newMessage := Message{
		Prompt:    "What is your favorite color?",
		Response:  "I don't have a favorite color, but I do like the color of a starry night sky!",
		Timestamp: time.Now(),
		UserID:    "987654321098765432",
	}

	// --- JSON FILE HANDLING ---
	fmt.Println("--- JSON File Example ---")
	jsonFilename := "messages.json"

	// Create an initial JSON file with one message.
	initialMessages := []Message{
		{
			Prompt:    "Hello, bot!",
			Response:  "Hello there!",
			Timestamp: time.Now().Add(-1 * time.Minute),
			UserID:    "123456789012345678",
		},
	}
	fmt.Printf("Writing initial messages to %s...\n", jsonFilename)
	if err := writeJSON(initialMessages, jsonFilename); err != nil {
		log.Fatalf("Error writing initial JSON file: %v", err)
	}

	// Append the new message to the existing JSON file.
	fmt.Printf("Appending a new message to %s...\n", jsonFilename)
	if err := appendJSON(newMessage, jsonFilename); err != nil {
		log.Fatalf("Error appending to JSON file: %v", err)
	}

	// Read all messages from the JSON file.
	messagesJSON, err := readJSON(jsonFilename)
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}
	fmt.Println("Successfully read all messages from the JSON file.")
	fmt.Printf("Messages found: %+v\n\n", messagesJSON)

	// --- NEW FUNCTIONALITY: READING AND FILTERING ---
	// Define the user ID we want to filter by.
	targetUserID := "123456789012345678"
	fmt.Printf("Reading and filtering messages for UserID: %s...\n", targetUserID)

	// Use the new function to read and filter the messages.
	filteredMessages, err := readJSONByUserID(jsonFilename, targetUserID)
	if err != nil {
		log.Fatalf("Error filtering JSON file: %v", err)
	}
	fmt.Println("Successfully filtered messages.")
	fmt.Printf("Filtered messages found: %+v\n", filteredMessages)
}

// writeJSON writes a slice of messages to a file, overwriting its contents.
func writeJSON(messages []Message, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a JSON encoder that will write to the file.
	encoder := json.NewEncoder(file)
	// Set indenting for a human-readable format.
	encoder.SetIndent("", "  ")
	// Encode the slice of messages and write it to the file.
	return encoder.Encode(messages)
}

// appendJSON reads all existing messages, adds a new one, and then
// writes the complete, updated slice back to the file.
func appendJSON(newMessage Message, filename string) error {
	// Read the existing messages from the file.
	messages, err := readJSON(filename)
	if err != nil {
		// If the file doesn't exist, we'll start with a new, empty slice.
		if os.IsNotExist(err) {
			messages = []Message{}
		} else {
			return err
		}
	}

	// Add the new message to the end of the slice.
	messages = append(messages, newMessage)

	// Overwrite the file with the full, updated slice of messages.
	return writeJSON(messages, filename)
}

// readJSON reads all messages from a JSON file and decodes them into a slice.
func readJSON(filename string) ([]Message, error) {
	// Open the file.
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

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

// readJSONByUserID reads all messages from a JSON file and returns only
// those that match the specified userID.
func readJSONByUserID(filename string, userID string) ([]Message, error) {
	// First, read all messages from the file.
	allMessages, err := readJSON(filename)
	if err != nil {
		return nil, err
	}

	// Create a new slice to hold the filtered messages.
	var filteredMessages []Message
	// Iterate through all messages and append the ones that match the userID.
	for _, msg := range allMessages {
		if msg.UserID == userID {
			filteredMessages = append(filteredMessages, msg)
		}
	}

	return filteredMessages, nil
}
