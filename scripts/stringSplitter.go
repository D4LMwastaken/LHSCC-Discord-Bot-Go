package scripts

/*
	Splits strings into chunks that are under 2000 characters so that they avoid
	the 2000 characters limit set by Discord.
*/

import "unicode/utf8"

const maxChunkLength = 2000 // Need to be changed if Discord increases character limit

func splitStringIntoChunks(s string) (firstChunk string, remainingChunks []string) {
	var allChunks []string

	if len(s) == 0 {
		return "", []string{}
	}

	for len(s) > 0 {
		runeCount := utf8.RuneCountInString(s)

		if runeCount <= maxChunkLength {
			allChunks = append(allChunks, s)
			break
		}

		currentLen := 0
		byteIndex := 0
		for i, r := range s {
			if currentLen >= maxChunkLength {
				byteIndex = i
				break
			}
			currentLen++
			byteIndex = i + utf8.RuneLen(r)
		}

		if byteIndex == 0 && len(s) > 0 {
			_, size := utf8.DecodeRuneInString(s)
			byteIndex = size
		} else if byteIndex > len(s) {
			byteIndex = len(s)
		}

		allChunks = append(allChunks, s[:byteIndex])
		s = s[byteIndex:]
	}

	if len(allChunks) > 0 {
		firstChunk = allChunks[0]
		if len(allChunks) > 1 {
			remainingChunks = allChunks[1:]
		}
	}

	return firstChunk, remainingChunks
}
