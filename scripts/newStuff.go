package scripts

var newStuffAdded = "1. Gemini can now generate images"

func NewStuff(author string, userID string) (string, []string) {
	prompt := "Generate a unique message to the author, " + author + " about all of the new stuff that the Discord Bot (you)" +
		" now have.\n" + "NewStuff: " + newStuffAdded
	first, restOfContent, _ := GeminiAI(prompt, author, userID, true, "ask", "none", "lite")
	return first, restOfContent
}
