package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"google.golang.org/genai"
	"log"
	"os"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	os.Getenv("GEMINI_API_KEY")
	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	config := &genai.GenerateContentConfig{
		ResponseModalities: []string{"TEXT", "IMAGE"},
	}

	result, _ := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash-preview-image-generation",
		genai.Text("Hi, can you create a 3d rendered image of the Largo High School Coding Club Discord Bot?"),
		config,
	)

	for _, part := range result.Candidates[0].Content.Parts { // 0 currently but very useful later on...
		if part.Text != "" {
			fmt.Println(part.Text)
		} else if part.InlineData != nil {
			imageBytes := part.InlineData.Data
			outputFilename := "gemini_generated_image.png"
			_ = os.WriteFile(outputFilename, imageBytes, 0644)
		}
	}
}
