package plugins

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type NotesSkill struct{}

func init() { Register(&NotesSkill{}) }

// --------------------------------------------------------------
// CanHandle
// --------------------------------------------------------------
func (n *NotesSkill) CanHandle(input string) bool {
	return strings.Contains(strings.ToLower(input), "note")
}

// --------------------------------------------------------------
// Handle
// --------------------------------------------------------------
func (n *NotesSkill) Handle(input string) error {
	lower := strings.ToLower(input)

	// 1️⃣ CLEAR NOTES
	if strings.Contains(lower, "clear note") ||
		strings.Contains(lower, "clear my notes") {

		os.WriteFile("data/notes.txt", []byte(""), 0644)

		ui := map[string]interface{}{
			"page":  "notes",
			"notes": "",
		}
		uiJson, _ := json.Marshal(ui)

		return fmt.Errorf("ui::%s|speak::All your notes have been cleared.", string(uiJson))
	}

	// 2️⃣ READ NOTES
	if strings.Contains(lower, "read") ||
		strings.Contains(lower, "show") ||
		strings.Contains(lower, "list") {

		text, lines := loadNotes()

		if len(lines) == 0 {
			return fmt.Errorf("speak::You don't have any notes yet.")
		}

		ui := map[string]interface{}{
			"page":  "notes",
			"notes": text,
		}
		uiJson, _ := json.Marshal(ui)

		spoken := "Here are your notes. " + strings.ReplaceAll(text, "\n", " ")

		return fmt.Errorf("ui::%s|speak::%s", string(uiJson), spoken)
	}

	// 3️⃣ ADD NOTE
	for _, phrase := range []string{
		"make a note of",
		"note that",
		"take a note of",
		"add a note",
		"set a note to",
		"set note to",
	} {
		if strings.Contains(lower, phrase) {
			content := strings.TrimSpace(strings.Replace(lower, phrase, "", 1))

			if content == "" {
				return fmt.Errorf("speak::What should I note down?")
			}

			appendNote(content)
			text, _ := loadNotes()

			ui := map[string]interface{}{
				"page":  "notes",
				"notes": text,
			}
			uiJson, _ := json.Marshal(ui)

			return fmt.Errorf("ui::%s|speak::Got it. I've added that to your notes.", string(uiJson))
		}
	}

	return nil
}

// --------------------------------------------------------------
// Helpers
// --------------------------------------------------------------
func appendNote(note string) error {
	os.MkdirAll("data", 0755)

	f, err := os.OpenFile("data/notes.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("- %s (%s)\n",
		note,
		time.Now().Format("Jan 02, 2006 3:04 PM")))
	return err
}

func loadNotes() (string, []string) {
	data, err := os.ReadFile("data/notes.txt")
	if err != nil || len(data) == 0 {
		return "", []string{}
	}

	text := string(data)
	scanner := bufio.NewScanner(strings.NewReader(text))
	var lines []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}

	return text, lines
}
