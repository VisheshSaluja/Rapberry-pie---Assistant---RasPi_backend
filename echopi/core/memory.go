package core

import (
	"fmt"
	"os"
	"strings"
)

const memoryFile = "data/memory.txt"

// GetMemoryContext safely loads memory or creates it if missing
func GetMemoryContext() string {
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		os.MkdirAll("data", 0755)
	}
	if _, err := os.Stat(memoryFile); os.IsNotExist(err) {
		os.WriteFile(memoryFile, []byte(""), 0644)
	}
	data, err := os.ReadFile(memoryFile)
	if err != nil {
		fmt.Printf("⚠️  Could not load memory: %v\n", err)
		return ""
	}
	return string(data)
}

// AddToMemory appends a new line and trims file to maxLines
func AddToMemory(entry string, maxLines int) {
	content := GetMemoryContext()
	lines := strings.Split(content, "\n")
	lines = append(lines, entry)
	if len(lines) > maxLines {
		lines = lines[len(lines)-maxLines:]
	}
	newData := strings.Join(lines, "\n")
	if err := os.WriteFile(memoryFile, []byte(newData), 0644); err != nil {
		fmt.Printf("⚠️  Could not update memory: %v\n", err)
	}
}

// ---- Compatibility helpers for integrator.go ----

// LoadMemory just ensures the file exists
func LoadMemory(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.WriteFile(path, []byte(""), 0644)
	}
	return nil
}

// ClearMemory wipes the memory file
func ClearMemory(path string) error {
	return os.WriteFile(path, []byte(""), 0644)
}

// SummarizeMemory returns the last N entries as a summary
func SummarizeMemory(n int) string {
	data := GetMemoryContext()
	lines := strings.Split(strings.TrimSpace(data), "\n")
	if len(lines) == 0 {
		return "I don't have any stored memory yet."
	}
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	return strings.Join(lines, "\n")
}

// SaveMemory rewrites the current memory file (noop here)
func SaveMemory(path string) error {
	// Already saved in AddToMemory, but keep for compatibility
	return nil
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