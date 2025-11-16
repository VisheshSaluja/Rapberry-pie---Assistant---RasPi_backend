
package core

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Transcribe converts an audio file to text using Whisper.cpp
func Transcribe(audioFile string) (string, error) {
	binary := os.Getenv("STT_BINARY")
	model := os.Getenv("STT_MODEL")
	if binary == "" || model == "" {
		return "", fmt.Errorf("STT_BINARY or STT_MODEL not set in .env")
	}

	var out bytes.Buffer
	cmd := exec.Command(binary, "-m", model, "-f", audioFile, "-nt", "-np")
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("Whisper failed: %v", err)
	}

	text := strings.TrimSpace(out.String())
	return text, nil
}
