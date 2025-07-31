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
			"2.Use Discord Markdown, Discord markdown does not have table support \n" +
			"3.If your message is longer than 2000 characters (message = everything that you sent), it will be " +
			"split up, make sure that you let it split so that words are not split apart or markdown is broken" +
			"Below is the message the member sent: \n"
		return theRules
	} else if ruleType == "sudo" {
		theRules = "Hi Gemini, this is your developer, D4LM, here are rules that you should follow when speaking to him: \n" +
			"1. This you developer that made this Discord Bot in Golang \n" +
			"2. Use Discord markdown to format your message, know that it does not have table support \n" +
			"3. "
	} else {
		print("No specific rules specified")
	}
	return theRules
}

func ModelCheck(modelName string) string {
	modelVersionName := ""
	if modelName == "pro" {
		modelVersionName = "gemini-2.5-pro"
		return modelVersionName
	} else if modelName == "flash" {
		modelVersionName = "gemini-2.5-flash"
		return modelVersionName
	} else if modelName == "lite" {
		modelVersionName = "gemini-2.5-flash-lite"
		return modelVersionName
	}
	return modelVersionName
}

func GeminiAI(prompt string, displayName string, splitString bool, ruleType string, Language string, model string) (string, []string) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ctx := context.Background() // Assumes API Key is set as GEMINI_API_KEY
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	modelVersionName := ModelCheck(model)

	result, err := client.Models.GenerateContent(
		ctx,
		modelVersionName,
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
