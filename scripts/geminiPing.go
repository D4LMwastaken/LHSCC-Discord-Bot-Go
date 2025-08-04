package scripts

import "time"

func PingGemini(prompt string, displayName string, userID string, splitString bool, ruleType string, Language string, model string) (string, []string, string) {
	now := time.Now()
	first, rest, _ := GeminiAI(prompt, displayName, userID, splitString, ruleType, Language, model)
	latency := time.Since(now).String()
	return first, rest, latency
}
