package core

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func Speak(text string) error {
	model := os.Getenv("PIPER_MODEL")
	if model == "" {
		return fmt.Errorf("PIPER_MODEL not set in .env")
	}

	cmd := exec.Command("piper", "-m", model, "-f", "output.wav")
	cmd.Stdin = bytes.NewBufferString(text)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("❌ Piper command failed: %v", err)
	}

	play := exec.Command("ffplay", "-nodisp", "-autoexit", "output.wav")
	play.Stdout = os.Stdout
	play.Stderr = os.Stderr
	return play.Run()
}
