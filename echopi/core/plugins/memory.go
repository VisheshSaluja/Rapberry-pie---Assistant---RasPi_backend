package plugins

import (
	"os"
	"strings"
)

// SummarizeMemory reads the last few lines of memory.txt
func SummarizeMemory() string {
	data, err := os.ReadFile("data/memory.txt")
	if err != nil {
		return "I couldn't access memory."
	}

	lines := strings.Split(string(data), "\n")
	n := len(lines)
	if n > 8 {
		lines = lines[n-8:]
	}
	return strings.Join(lines, ", ")
}

// ClearMemoryFile deletes the memory.txt content
func ClearMemoryFile() error {
	if err := os.MkdirAll("data", 0755); err != nil {
		return err
	}
	return os.WriteFile("data/memory.txt", []byte{}, 0644)
}

func EnsureMemoryFile() error {
	// Make sure "data" folder exists
	if err := os.MkdirAll("data", 0755); err != nil {
		return err
	}

	// If memory.txt doesn’t exist, create an empty one
	if _, err := os.Stat("data/memory.txt"); os.IsNotExist(err) {
		return os.WriteFile("data/memory.txt", []byte(""), 0644)
	}

	return nil
}