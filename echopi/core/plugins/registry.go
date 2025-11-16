package plugins

import "strings"

// Skill defines the interface every plugin must implement.
type Skill interface {
	CanHandle(input string) bool
	Handle(input string) error // returns special formatted errors like:
	// "ui::pomodoro_start|speak::Starting session"
}

// Global registry
var skills []Skill

// Register plugin during init()
func Register(skill Skill) {
	skills = append(skills, skill)
}

// List number of registered plugins
func List() int {
	return len(skills)
}

// -----------------------------------------------------------------------------
// ❗ Old TryHandle (USED BY MAC LOCAL HANDLING) — KEEP
// -----------------------------------------------------------------------------
func TryHandle(input string) (bool, error) {
	in := strings.ToLower(input)
	for _, skill := range skills {
		if skill.CanHandle(in) {
			err := skill.Handle(in)
			return true, err
		}
	}
	return false, nil
}

// -----------------------------------------------------------------------------
// NEW: Unified plugin response struct
// -----------------------------------------------------------------------------
type UnifiedResult struct {
	Status string `json:"status"` // always "handled"
	Speak  string `json:"speak,omitempty"`
	UI     string `json:"ui,omitempty"`
	Data   any    `json:"data,omitempty"`
}

// -----------------------------------------------------------------------------
// NEW: Unified router for Pi backend
// -----------------------------------------------------------------------------
func TryHandleUnified(input string) (bool, UnifiedResult) {

	for _, skill := range skills {
		if skill.CanHandle(strings.ToLower(input)) {

			err := skill.Handle(strings.ToLower(input))

			// If plugin returned no structured response:
			if err == nil {
				return true, UnifiedResult{Status: "handled"}
			}

			msg := err.Error()

			var uiCmd string
			var speak string

			// Extract speak::
			if strings.Contains(msg, "speak::") {
				speak = extract(msg, "speak::")
			}

			// Extract ui::
			if strings.Contains(msg, "ui::") {
				uiCmd = extract(msg, "ui::")
			}

			return true, UnifiedResult{
				Status: "handled",
				Speak:  speak,
				UI:     uiCmd,
			}
		}
	}

	return false, UnifiedResult{}
}

// Helper: extract substring after prefix until next "|"
func extract(text, prefix string) string {
	parts := strings.Split(text, prefix)
	if len(parts) < 2 {
		return ""
	}
	after := parts[1]

	// Cut at next section if exists
	if idx := strings.Index(after, "|"); idx != -1 {
		return strings.TrimSpace(after[:idx])
	}
	return strings.TrimSpace(after)
}
