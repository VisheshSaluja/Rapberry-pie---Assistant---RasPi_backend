package plugins

import (
    "encoding/json"
    "fmt"
    "strings"
    "time"
)

type PomodoroSkill struct {
    running bool
    start   time.Time
}

func init() { Register(&PomodoroSkill{}) }

// --------------------------------------------------------------
// CanHandle — trigger when user talks about pomodoro/timer
// --------------------------------------------------------------
func (p *PomodoroSkill) CanHandle(input string) bool {
    l := strings.ToLower(input)

    return strings.Contains(l, "pomodoro") ||
        strings.Contains(l, "timer") ||
        strings.Contains(l, "focus session") ||
        strings.Contains(l, "focus mode")
}

// --------------------------------------------------------------
// Handle Pomodoro voice commands
// --------------------------------------------------------------
func (p *PomodoroSkill) Handle(input string) error {
    l := strings.ToLower(input)

    // -------------------------
    // 1️⃣ STOP / CANCEL
    // -------------------------
    if strings.Contains(l, "stop") ||
        strings.Contains(l, "cancel") ||
        strings.Contains(l, "end") {

        p.running = false

        payload := map[string]interface{}{
            "page":    "pomodoro",
            "running": false,
            "minutes": "25",
            "seconds": "00",
        }
        uiJson, _ := json.Marshal(payload)

        return fmt.Errorf("ui::%s|speak::Pomodoro stopped. Back to idle mode.", string(uiJson))
    }

    // -------------------------
    // 2️⃣ PAUSE
    // -------------------------
    if strings.Contains(l, "pause") {
        if !p.running {
            return fmt.Errorf("speak::There's no active pomodoro to pause.")
        }

        p.running = false

        mins, secs := p.remaining()

        payload := map[string]interface{}{
            "page":    "pomodoro",
            "running": false,
            "minutes": fmt.Sprintf("%02d", mins),
            "seconds": fmt.Sprintf("%02d", secs),
        }
        uiJson, _ := json.Marshal(payload)

        return fmt.Errorf("ui::%s|speak::Pomodoro paused.", string(uiJson))
    }

    // -------------------------
    // 3️⃣ RESUME
    // -------------------------
    if strings.Contains(l, "resume") ||
        strings.Contains(l, "continue") {

        p.running = true
        p.start = time.Now()

        go p.runCountdown() // 🔥 LIVE TIMER LOOP

        mins, secs := p.remaining()
        payload := map[string]interface{}{
            "page":    "pomodoro",
            "running": true,
            "minutes": fmt.Sprintf("%02d", mins),
            "seconds": fmt.Sprintf("%02d", secs),
        }
        uiJson, _ := json.Marshal(payload)

        return fmt.Errorf("ui::%s|speak::Resuming pomodoro.", string(uiJson))
    }

    // -------------------------
    // 4️⃣ START POMODORO
    // -------------------------
    if strings.Contains(l, "start") ||
        strings.Contains(l, "begin") ||
        strings.Contains(l, "focus") {

        p.running = true
        p.start = time.Now()

        go p.runCountdown() // 🔥 LIVE TIMER LOOP

        payload := map[string]interface{}{
            "page":    "pomodoro",
            "running": true,
            "minutes": "25",
            "seconds": "00",
        }
        uiJson, _ := json.Marshal(payload)

        return fmt.Errorf("ui::%s|speak::Starting a 25 minute pomodoro. Focus mode activated.", string(uiJson))
    }

    // Fallback
    return fmt.Errorf("speak::I recognize pomodoro commands, but I wasn't sure what you wanted.")
}

// --------------------------------------------------------------
// Helper: remaining countdown
// --------------------------------------------------------------
func (p *PomodoroSkill) remaining() (int, int) {
    if !p.running {
        return 25, 0
    }

    elapsed := time.Since(p.start)
    totalSeconds := 25*60 - int(elapsed.Seconds())

    if totalSeconds < 0 {
        totalSeconds = 0
    }

    min := totalSeconds / 60
    sec := totalSeconds % 60

    return min, sec
}

// --------------------------------------------------------------
// Countdown Loop — emits UI updates every second
// --------------------------------------------------------------
func (p *PomodoroSkill) runCountdown() {
    for p.running {
        mins, secs := p.remaining()

        ui := map[string]interface{}{
            "page":    "pomodoro",
            "running": true,
            "minutes": fmt.Sprintf("%02d", mins),
            "seconds": fmt.Sprintf("%02d", secs),
        }

        EmitUI(ui)

        if mins == 0 && secs == 0 {
            p.running = false
            EmitUI(map[string]interface{}{
                "page":    "pomodoro",
                "running": false,
                "minutes": "25",
                "seconds": "00",
            })
            return
        }

        time.Sleep(1 * time.Second)
    }
}
