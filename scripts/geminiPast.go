package scripts

import "log"

func geminiPast(userID string) string {
	history := ""
	allMessages, err := readJSON(filename)
	if err != nil {
		log.Printf("Error reading messages.json: %v", err)
	}
	for message := range allMessages {
		if allMessages[message].UserID == userID {
			history = history + "Prompt: " + allMessages[message].Prompt + "\n" +
				"Responses: " + allMessages[message].Response + "\n" +
				"Timestamp: " + allMessages[message].Timestamp
		}
	}
	return history
}
