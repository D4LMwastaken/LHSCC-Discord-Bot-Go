package scripts

var newStuffAdded = "1. Gemini now stores all prompts and responses to now be personalized to each Discord user "

func NewStuff(author string, userID string) (string, []string) {
	prompt := "Generate a unique message to the author, " + author + " about all of the new stuff that the Discord Bot (you)" +
		" now have.\n" + "NewStuff: " + newStuffAdded
	first, restOfContent := GeminiAI(prompt, author, userID, true, "ask", "none", "lite")
	return first, restOfContent
}
