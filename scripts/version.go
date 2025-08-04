package scripts

var (
	version  = "v0.90-alpha"
	language = "Go"
	dev      = "D4LM"
	codev    = "No one is the apprentice yet!"
)

func Version(author string, userID string) (string, []string) {
	prompt := "Generate a unique message telling " + author + "what version and these specific information about this Discord Bot \n" +
		"Version: " + version + "\n" +
		"Language: " + language + "\n" +
		"Dev: " + dev + "\n" +
		"Codev: " + codev + "\n"
	content, rest, _ := GeminiAI(prompt, author, userID, true, "ask", "none", "lite")
	return content, rest
}
