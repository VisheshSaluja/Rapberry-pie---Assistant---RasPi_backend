package core

import (
	"fmt"
	"strings"
	"time"
)

// StartContinuousListening runs an infinite loop that listens for the hotword.
func StartContinuousListening() {
	hotword := GetHotword()
	fmt.Println("🔊 Initializing EchoPi Phase 7: Continuous Listening Mode...")
	fmt.Println("🎧  EchoPi continuous mode activated.")
	fmt.Printf("🟢  Say 'Hey %s' anytime to trigger conversation. Press Ctrl+C to exit.\n", strings.Title(hotword))

	for {
		fmt.Println("🎙️  Recording 4s of audio snippet...")
		err := RecordAudio("tmp_listen.wav", 4)
		if err != nil {
			fmt.Printf("❌  Mic capture failed: %v\n", err)
			time.Sleep(2 * time.Second)
			continue
		}

		transcript, err := Transcribe("tmp_listen.wav")
		if err != nil {
			fmt.Printf("⚠️  STT failed: %v\n", err)
			time.Sleep(2 * time.Second)
			continue
		}

		lower := strings.ToLower(strings.TrimSpace(transcript))
		if lower == "" {
			continue
		}

		fmt.Printf("📝  Transcript snippet: %s\n", transcript)

		// --- Hotword detection ---
		if strings.Contains(lower, "hey "+hotword) ||
			strings.Contains(lower, "hello "+hotword) ||
			strings.Contains(lower, hotword) {
			fmt.Printf("🪄  Wake-word '%s' detected! Starting full conversation…\n", hotword)

			if err := RunEchoPi(); err != nil {
				fmt.Printf("❌  Error in EchoPi run: %v\n", err)
			}

			// Graceful delay before resuming listening
			fmt.Println("🔁  Returning to background listening...")
			time.Sleep(1 * time.Second)
		}
	}
}
