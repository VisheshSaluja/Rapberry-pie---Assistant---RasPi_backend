package plugins

import (
	"sync"
	"time"
)

type PomodoroUI struct {
	Running     bool `json:"running"`
	SecondsLeft int  `json:"secondsLeft"`
}

type UIState struct {
	ActiveView   string     `json:"activeView"`             // "idle", "notes", "pomodoro", etc.
	Notes        []string   `json:"notes,omitempty"`        // notes lines
	Pomodoro     PomodoroUI `json:"pomodoro"`               // pomodoro info
	HotwordActive bool      `json:"hotwordActive"`          // true briefly after "hey nova"
}

var (
	uiState = UIState{
		ActiveView: "idle",
		Pomodoro:   PomodoroUI{Running: false, SecondsLeft: 0},
	}
	lastHotword time.Time
	mu          sync.RWMutex
)

func MarkHotword() {
	mu.Lock()
	defer mu.Unlock()
	lastHotword = time.Now()
}

func SetNotesUI(notes []string) {
	mu.Lock()
	defer mu.Unlock()
	uiState.ActiveView = "notes"
	uiState.Notes = notes
}

func SetPomodoroUI(running bool, seconds int) {
	mu.Lock()
	defer mu.Unlock()
	uiState.ActiveView = "pomodoro"
	uiState.Pomodoro = PomodoroUI{
		Running:     running,
		SecondsLeft: seconds,
	}
}

func SetIdleUI() {
	mu.Lock()
	defer mu.Unlock()
	uiState.ActiveView = "idle"
}

// Called by HTTP handler to expose current state
func GetUIState() UIState {
	mu.RLock()
	defer mu.RUnlock()

	state := uiState // copy

	if !lastHotword.IsZero() && time.Since(lastHotword) < 4*time.Second {
		state.HotwordActive = true
	} else {
		state.HotwordActive = false
	}

	return state
}
