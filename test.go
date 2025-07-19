package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"google.golang.org/genai"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		return
	}
	ctx := context.Background()
	// The client gets the API key from the environment variable `GEMINI_API_KEY`.
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text("Explain how AI works in a few words"),
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.Text())
}
