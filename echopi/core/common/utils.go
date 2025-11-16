package common

import (
	"fmt"
	"os/exec"
)

// Speak uses ffplay to play a generated audio file.
func Speak(text string) error {
	fmt.Printf("🗣️ Speaking: %s\n", text)
	cmd := exec.Command("say", text) // macOS; on Linux replace with 'espeak'
	return cmd.Run()
}

// (Optional) Add other shared helpers here later,
// e.g., GetTime(), SaveToFile(), LoadConfig(), etc.
