package scripts

func Help(commandNames string, commandDescriptions string, author string) (string, []string) {
	prompt := "Generate a unique help message telling the " + author + " what these commands do.\n" +
		"Here are all the command names seperated by comma: " + commandNames + "\n" +
		"Here are all the command descriptions seperated by comma: " + commandDescriptions
	firstContent, restOfContent := GeminiAI(prompt, author, true, "ask", "none", "lite")
	return firstContent, restOfContent
}
