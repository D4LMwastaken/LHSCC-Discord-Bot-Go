package scripts

var (
	version  = "v0.88"
	language = "Go"
	dev      = "D4LM"
	codev    = "No one is the apprentice yet!"
)

func Version(author string) (string, []string) {
	prompt := "Generate a unique message telling " + author + "what version and these specific information about this Discord Bot \n" +
		"Version: " + version + "\n" +
		"Language: " + language + "\n" +
		"Dev: " + dev + "\n" +
		"Codev: " + codev + "\n"
	content, rest := GeminiAI(prompt, author, true, "ask", "none", "lite")
	return content, rest
}
