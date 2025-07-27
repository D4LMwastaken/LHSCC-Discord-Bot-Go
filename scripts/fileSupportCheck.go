package scripts

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
	fileExtension := ""
	fileType := ""
	if language == "python" || language == "python3" || language == "py" {
		fileExtension = ".py"
		fileType = "python"
		return fileExtension, fileType
	} else if language == "c" || language == "c++" || language == "cpp" {
		fileExtension = ".cpp"
		fileType = "cpp"
		return fileExtension, fileType
	} else if language == "c#" || language == "csharp" || language == "csh" {
		fileExtension = ".csh"
		fileType = "csh"
		return fileExtension, fileType
	} else if language == "rust" || language == "rs" {
		fileExtension = ".rs"
		fileType = "rs"
		return fileExtension, fileType
	} else if language == "java" {
		fileExtension = ".java"
		fileType = "java"
		return fileExtension, fileType
	} else if language == "js" {
		fileExtension = ".js"
		fileType = "js"
		return fileExtension, fileType
	} else if language == "go" || language == "golang" {
		fileExtension = ".go"
		fileType = "go"
		return fileExtension, fileType
	} else if language == "lua" {
		fileExtension = ".lua"
		fileType = "lua"
		return fileExtension, fileType
	} else if language == "typescript" {
		fileExtension = ".ts"
		fileType = "ts"
	} else {
		fileExtension = ".txt"
		fileType = "plain"
	}
	return fileExtension, fileType
}
