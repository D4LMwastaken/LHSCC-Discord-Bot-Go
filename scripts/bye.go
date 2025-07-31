package scripts

import "log"

func Bye(author string, username string, id string) (string, []string, bool) {
	content := ""
	var restOfContent []string
	isDev := false
	if username == "d4lm." {
		if id == "1125963396177723452" {
			isDev = true
			prompt := "Generate a unique goodbye as the detected person is the developer and the bot is shutting down " +
				"shortly and mention the author"
			content, restOfContent = GeminiAI(prompt, author, true, "sudo", "none", "pro")
			return content, restOfContent, isDev
		} else {
			log.Panic("ID Verification failed, someone is took my username.")
			prompt := "Generate a unique goodbye and the detected person is the not the devloper"
			content, restOfContent = GeminiAI(prompt, author, true, "ask", "none", "lite")
			return content, restOfContent, isDev
		}
		return content, restOfContent, isDev
	} else {
		prompt := "Generate a unique goodbye for the person"
		content, restOfContent = GeminiAI(prompt, author, true, "ask", "none", "lite")
		return content, restOfContent, isDev
	}
	return content, restOfContent, isDev
}
