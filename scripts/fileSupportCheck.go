package scripts

import "strings"

/* Will be used later for another slash command,
var currentLanguageSupport = "Here are all the languages supported by file: \n" +
	"* Python (python3, py, python)\n" +
	"* C (c, c++, cpp)\n" +
	"* C# (c#, csharp, csh)\n" +
	"* Rust (rust, rs)\n" +
	"* \n" +
	"* \n" +
	"* \n" +
	"* \n" +
	"* \n" +
	"* \n"
*/

func FileSupportCheck(language string) (string, string) {
	undercaseLanguage := strings.ToLower(language)
	fileExtension := ""
	fileType := ""
	if undercaseLanguage == "python" || language == "python3" || language == "py" {
		fileExtension = ".py"
		fileType = "python"
		return fileExtension, fileType
	} else if undercaseLanguage == "c" || language == "c++" || language == "cpp" {
		fileExtension = ".cpp"
		fileType = "cpp"
		return fileExtension, fileType
	} else if undercaseLanguage == "c#" || language == "csharp" || language == "csh" {
		fileExtension = ".csh"
		fileType = "csh"
		return fileExtension, fileType
	} else if undercaseLanguage == "rust" || language == "rs" {
		fileExtension = ".rs"
		fileType = "rs"
		return fileExtension, fileType
	} else if undercaseLanguage == "java" {
		fileExtension = ".java"
		fileType = "java"
		return fileExtension, fileType
	} else if undercaseLanguage == "js" || language == "javascript" {
		fileExtension = ".js"
		fileType = "js"
		return fileExtension, fileType
	} else if undercaseLanguage == "go" || language == "golang" {
		fileExtension = ".go"
		fileType = "go"
		return fileExtension, fileType
	} else if undercaseLanguage == "lua" {
		fileExtension = ".lua"
		fileType = "lua"
		return fileExtension, fileType
	} else if undercaseLanguage == "typescript" || language == "ts" {
		fileExtension = ".ts"
		fileType = "ts"
	} else if undercaseLanguage == "html" {
		fileExtension = ".html"
		fileType = "html"
	} else {
		fileExtension = ".txt"
		fileType = "plain"
	}
	return fileExtension, fileType
}
