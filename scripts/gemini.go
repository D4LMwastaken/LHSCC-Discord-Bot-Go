package scripts

import (
	"context"
	"github.com/joho/godotenv"
	"google.golang.org/genai"
	"log"
)

func rules(displayName string, userID string, ruleType string, language string) string {
	history := geminiPast(userID)
	theRules := ""
	if ruleType == "file" {
		theRules = "Hi Gemini, here are all the rules you should follow when generating content for the Largo High School" +
			" Coding Club Discord Server: \n" +
			"1.The member who sent you the message is: " + displayName + "\n" +
			"2.Use the programming language the user specifies, if not use pseudocode. The user wants you to code in" + language + "\n" +
			"3. You are generating text into a file, so code only and put in comments anything you generate that is not code" +
			"4. At the first line of the file, put a comment with the file name and file extension" +
			"5. Do not have ``` at the start or the end, they are not necessary \n" +
			"6. Here is the messages the author has sent before as context: " + history + "\n" +
			"Below is the message the member sent: \n"
		return theRules
	} else if ruleType == "ask" {
		theRules = "Hi Gemini, here are all the rules you should follow when generating content for the Largo High School" +
			" Coding Club Discord Server: \n" +
			"1.The member who sent you the message is: " + displayName + "\n" +
			"2.Use Discord Markdown, Discord markdown does not have table support \n" +
			"3.If your message is longer than 2000 characters (message = everything that you sent), it will be " +
			"split up, make sure that you let it split so that words are not split apart or markdown is broken" +
			"4. Here is the messages the author has sent before as context: " + history + "\n" +
			"Below is the message the member sent: \n"
		return theRules
	} else if ruleType == "sudo" {
		theRules = "Hi Gemini, this is your developer, D4LM, here are rules that you should follow when speaking to him: \n" +
			"1. This you developer that made this Discord Bot in Golang \n" +
			"2. Use Discord markdown to format your message, know that it does not have table support \n" +
			"3. Here is the messages the author has sent before as context: " + history + "\n" +
			"Below is the message the developer has sent: "
		return theRules
	} else if ruleType == "image" {
		theRules = "Hi Gemini, here are all the rules you should follow when generating images for the Largo High School " +
			"Coding Club Discord Server: \n" +
			"1.The member who sent you the message is: " + displayName + "\n" +
			"2.Use Discord Markdown, Discord markdown does not have table support \n" +
			"3.Here is the messages the author has sent before as context: " + history + "\n" +
			"4. Do not include the rules when generating images \n" +
			"Below is the message the member sent: \n"
	} else {
		print("No specific rules specified")
	}
	return theRules
}

func ModelCheck(modelName string) string {
	modelVersionName := ""
	if modelName == "pro" {
		modelVersionName = "gemini-2.5-pro"
	} else if modelName == "flash" {
		modelVersionName = "gemini-2.5-flash"
	} else if modelName == "lite" {
		modelVersionName = "gemini-2.5-flash-lite"
	} else if modelName == "image" {
		modelVersionName = "gemini-2.0-flash-preview-image-generation"
	} else {
		modelVersionName = "gemini-2.5-flash"
	}
	return modelVersionName
}

func GeminiAI(prompt string, displayName string, userID string, splitString bool, ruleType string, Language string, model string) (string, []string, []byte) {
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

	var config *genai.GenerateContentConfig
	if model == "image" {
		config = &genai.GenerateContentConfig{
			ResponseModalities: []string{"TEXT", "IMAGE"},
		}
	}

	result, err := client.Models.GenerateContent(
		ctx,
		modelVersionName,
		genai.Text(rules(displayName, userID, ruleType, Language)+prompt),
		config,
	)
	if err != nil {
		log.Fatal(err)
	}

	if model == "image" {
		content := ""
		var image []byte
		for _, part := range result.Candidates[0].Content.Parts {
			if part.Text != "" {
				content += part.Text
			} else if part.InlineData != nil {
				image = part.InlineData.Data // Assumed to be only one image
			}
		}
		go geminiSaver(prompt, content, userID, displayName, image)
		return content, nil, image
	} else {
		go geminiSaver(prompt, result.Text(), userID, displayName, nil)
		// ^ Goroutine in case someone decides to have a long prompt or long response
		if splitString == true {
			firstContent, restOfContent := splitStringIntoChunks(result.Text())
			return firstContent, restOfContent, nil
		} else {
			return result.Text(), nil, nil
		}
	}
}
