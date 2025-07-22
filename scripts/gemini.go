package scripts

import (
	"context"
	"github.com/joho/godotenv"
	"google.golang.org/genai"
	"log"
)

func GeminiAI(prompt string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ctx := context.Background() // Assumes API Key is set as GEMINI_API_KEY
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	return result.Text()
}
