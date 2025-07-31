package scripts

var newStuffAdded = "1. Migrated all non-calendar commands to the Discord Bot from the original python Discord Bot"

func NewStuff(author string) (string, []string) {
	prompt := "Generate a unique message to the author, " + author + " about all of the new stuff that the Discord Bot (you)" +
		" now have.\n" + "NewStuff: " + newStuffAdded
	first, restOfContent := GeminiAI(prompt, author, true, "ask", "none", "lite")
	return first, restOfContent
}
