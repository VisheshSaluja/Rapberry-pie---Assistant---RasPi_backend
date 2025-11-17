# Backend — Navo Productivity Assistant
The repoistory powers Navo, a voice-driven productivity assistant.
It runs entirely on a Raspberry Pi, acting as the brain of the device.
**The backend is responsible for:**

- Speech-to-Text using Whisper
- Reasoning and intent parsing using Ollama
- Text-to-Speech using Piper
- Notes and Pomodoro modules
- Real-time communication with the frontend using WebSockets

  ## Requirements

- Raspberry Pi 3/4
- Go 1.21+
- Whisper.cpp or Whisper API
- Piper TTS
- Ollama running on Pi or remote machine
- Speaker output connected to Pi

  ## Setup
  #### 1. Install Dependencies
  ```bash
  sudo apt install git ffmpeg build-essential -y
  ```
  #### 2. Clone repo
    ```bash
  git clone <repo-url>
  cd backend
  ```
  #### 3. Install Go modules
    ```bash
  go mod tidy
  ```
  #### 4. Start Whisper / Ollama / Piper
  Depending on your setup:
    ```bash
  ./whisper -m models/base.en --in

  ollama run mistral

  piper -m en_US-amy.low --output-raw
  ```
  #### 5. Run the Backend Server
    ```bash
  go run main.go
  ```
  The WebSocket server will be available at:
  ```bash
    ws://<raspberry-pi-ip>:9001/ws
  ```

## Testing
#### You can test the backend using:
WebSocket test extensions

CURL for REST endpoints

Svelte UI frontend in the other repo
