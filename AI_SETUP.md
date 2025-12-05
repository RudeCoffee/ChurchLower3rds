# Local AI Setup Guide

 This application uses a Local AI model to provide smart verse suggestions. This runs entirely on your machine and does not send data to the cloud.

 ## Prerequisites

 1.  **Install Ollama**:
     -   Download and install Ollama from [ollama.com](https://ollama.com).

 2.  **Pull the Model**:
     -   Open your terminal/command prompt.
     -   Run the following command to download the model (approx 2GB):
         ```bash
         ollama pull llama3.2
         ```
     -   Alternatively, you can use `phi3` or `tinyllama` if you have limited RAM, but `llama3.2` is recommended for best results.
         ```bash
         ollama pull phi3
         ```
     -   *Note: If you change the model name, update `main.go` line where `NewAIClient` is called.*

 3.  **Start Ollama**:
     -   Ensure Ollama is running (it usually starts automatically in the background).
     -   You can verify it's running by visiting `http://localhost:11434` in your browser.

 ## Server-Side Audio Transcription (Faster)

 The browser-based speech recognition can sometimes be slow or unreliable. For faster response times, you can run the audio listener directly on the server (the media computer).

 ### Setup

 1.  **Install Python Dependencies**:
     ```bash
     pip install -r requirements.txt
     ```
     *Note: You may need to install PortAudio system libraries (e.g., `sudo apt-get install portaudio19-dev` on Linux or `brew install portaudio` on Mac) for `pyaudio` to work.*

 2.  **Run the Audio Listener**:
     -   Make sure your microphone is connected and set as default.
     -   Run the script:
         ```bash
         python3 audio_listener.py
         ```
     -   The script will listen to the audio, transcribe it locally using Whisper, and send it to the Go server.
     -   The Control Client web interface will show "Server Audio Active".

 ## Usage

 1.  Start the Church Lower 3rds application (`go run main.go`).
 2.  (Optional) Start the server-side audio listener (`python3 audio_listener.py`).
 3.  Open the Control Client (`http://localhost:8080/client.html`).
 4.  Click the **"Local AI Assistant"** toggle switch.
 5.  Start speaking or listening to the sermon.
 6.  When the AI detects a verse reference or a transition, a blue suggestion box will appear.
 7.  Click **"Apply Verse"** to immediately display it.

 ## Troubleshooting

 -   **"Local AI (Ollama) not detected"**: Make sure Ollama is running on port 11434.
 -   **AI suggestions are slow**: Local AI depends on your computer's speed (GPU/CPU).
 -   **AI is hallucinating verses**: This can happen with smaller models. Verify the suggestion before clicking Apply.
