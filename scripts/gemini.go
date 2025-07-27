package scripts

import (
	"context"
	"github.com/joho/godotenv"
	"google.golang.org/genai"
	"log"
)

func rules(displayName string, ruleType string, language string) string {
	theRules := ""
	if ruleType == "file" {
		theRules = "Hi Gemini, here are all the rules you should follow when generating content for the Largo High School" +
			" Coding Club Discord Server: \n" +
			"1.The member who sent you the message is: " + displayName + "\n" +
			"2.Use the programming language the user specifies, if not use pseudocode. The user wants you to code in" + language + "\n" +
			"3. You are generating text into a file, so code only and put in comments anything you generate that is not code" +
			"4. At the first line of the file, put a comment with the file name and file extension" +
			"5. Do not have ``` at the start or the end, they are not necessary \n" +
			"Below is the message the member sent: \n"
		return theRules
	} else if ruleType == "ask" {
		theRules = "Hi Gemini, here are all the rules you should follow when generating content for the Largo High School" +
			" Coding Club Discord Server: \n" +
			"1.The member who sent you the message is: " + displayName + "\n" +
			"2.Use Discord Markdown \n" +
			"Below is the message the member sent: \n"
		return theRules
	} else {
		print("No rules specified")
	}
	return theRules
}

func GeminiAI(prompt string, displayName string, splitString bool, ruleType string, Language string) (string, []string) {
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
		"gemini-2.5-pro",
		genai.Text(rules(displayName, ruleType, Language)+prompt),
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	if splitString == true {
		return splitStringIntoChunks(result.Text())
	} else {
		return result.Text(), nil
	}
}
