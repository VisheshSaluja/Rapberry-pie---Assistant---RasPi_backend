package core

import (
	"fmt"
	"strings"
	"time"
)

// RunEchoPi orchestrates the complete STT → Skill → LLM → TTS pipeline with memory control
func RunEchoPi() error {
	fmt.Println("🎙️  Listening... Recording input.wav")

	// --- Ensure memory file exists ---
	if err := EnsureMemoryFile(); err != nil {
		fmt.Printf("⚠️  Could not prepare memory file: %v\n", err)
	}

	// --- Load memory ---
	if err := LoadMemory("data/memory.txt"); err != nil {
		fmt.Printf("⚠️  Could not load memory: %v\n", err)
	}

	// --- Record user input ---
	if err := RecordAudio("input.wav", 8); err != nil {
		return fmt.Errorf("mic capture failed: %v", err)
	}

	// --- Transcribe via Whisper ---
	fmt.Println("🧠  Transcribing speech via Whisper...")
	transcript, err := Transcribe("input.wav")
	if err != nil {
		return fmt.Errorf("STT failed: %v", err)
	}

	// --- Handle blank or low-volume audio ---
	lower := strings.ToLower(strings.TrimSpace(transcript))
	if lower == "" || strings.Contains(lower, "[blank_audio]") {
		fmt.Println("⚠️  Blank audio detected. Retrying with volume boost (x2.5)...")
		if err := RecordAudio("input.wav", 8); err != nil {
			return fmt.Errorf("retry mic capture failed: %v", err)
		}
		transcript, err = Transcribe("input.wav")
		if err != nil {
			return fmt.Errorf("retry STT failed: %v", err)
		}
		lower = strings.ToLower(strings.TrimSpace(transcript))
		fmt.Printf("📝  Retry Transcript: %s\n", transcript)
	} else {
		fmt.Printf("📝  Transcript: %s\n", transcript)
	}

	// --- Hotword detection (dynamic from .env) ---
	hotword := GetHotword()
	if !strings.Contains(lower, "hey "+hotword) &&
		!strings.Contains(lower, "hello "+hotword) &&
		!strings.Contains(lower, hotword) {
		fmt.Println("🤫  Hotword not detected — ignoring input.")
		return nil
	}

	// --- Memory voice commands ---
	if strings.Contains(lower, "clear memory") {
		if err := ClearMemory("data/memory.txt"); err != nil {
			return fmt.Errorf("failed to clear memory: %v", err)
		}
		fmt.Println("🧹 Memory cleared successfully.")
		_ = Speak("I've cleared all previous memory.")
		return nil
	}

	if strings.Contains(lower, "what do you remember") || strings.Contains(lower, "recall memory") {
		summary := SummarizeMemory(5)
		fmt.Println("🧠 Memory summary:\n", summary)
		_ = Speak("Here's what I remember: " + summary)
		return nil
	}

	// --- Skill routing before LLM ---
	handled, err := HandleSkill(transcript)
	if err != nil {
		if err.Error() == "sleep" {
			return nil // graceful stop
		}
		return fmt.Errorf("skill error: %v", err)
	}
	if handled {
		return nil // plugin already handled it
	}

	// --- Query LLM (only if no skill matched) ---
	fmt.Println("🤖  Querying Ollama...")
	prompt := GetMemoryContext() + "\nUser: " + transcript + "\nAssistant:"
	response, err := QueryLLM(prompt)
	if err != nil {
		return fmt.Errorf("LLM query failed: %v", err)
	}
	fmt.Printf("💬  Response: %s\n", response)

	// --- Log and update memory ---
	if err := LogInteraction(transcript, response); err != nil {
		fmt.Printf("⚠️  Could not log conversation: %v\n", err)
	}
	AddToMemory("USER: "+transcript, 6)
	AddToMemory("ECHO: "+response, 6)
	if err := SaveMemory("data/memory.txt"); err != nil {
		fmt.Printf("⚠️  Could not save memory: %v\n", err)
	}

	// --- Speak result ---
	fmt.Println("🔊  Speaking response...")
	if err := Speak(response); err != nil {
		return fmt.Errorf("TTS failed: %v", err)
	}

	fmt.Println("✅  EchoPi conversation complete.")
	time.Sleep(1 * time.Second)
	return nil
}
