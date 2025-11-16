package plugins

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

// ⚠️ CHANGE THIS TO YOUR MAC'S LAN IP
const MacRelayURL = "http://192.168.1.169:9001/forward"

func EmitUI(payload map[string]interface{}) {
    payload["type"] = "ui"

    jsonBody, err := json.Marshal(payload)
    if err != nil {
        fmt.Println("❌ EmitUI marshal error:", err)
        return
    }

    resp, err := http.Post(
        MacRelayURL,
        "application/json",
        bytes.NewReader(jsonBody),
    )
    if err != nil {
        fmt.Println("⚠️ EmitUI error:", err)
        return
    }
    _ = resp.Body.Close()
}
