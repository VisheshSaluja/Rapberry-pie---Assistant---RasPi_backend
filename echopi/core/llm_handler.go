

package core

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// QueryLLM sends prompt to Ollama and returns the model’s response
func QueryLLM(prompt string) (string, error) {
	model := os.Getenv("LLM_MODEL")
	if model == "" {
		model = "llama3.2:1b"
	}

	cmd := exec.Command("ollama", "run", model, prompt)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("Ollama run failed: %v", err)
	}
	return out.String(), nil
}
