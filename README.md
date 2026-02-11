# Church Lower 3rds

A real-time Bible verse and speaker name display system for OBS (Open Broadcaster Software) designed for church services and live streaming. This version features **local AI-powered speech recognition** to suggest verses in real-time.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.23+-blue.svg)
![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20Linux%20%7C%20macOS-lightgrey.svg)

## Features

### 🎙️ AI Voice Assistant
- **Local Speech Recognition**: Uses **Vosk** to transcribe audio locally on your machine (no internet required).
- **Smart Verse Search**: Analyzes the live sermon to intelligently suggest relevant Bible verses, even if paraphrased.
- **Privacy First**: All processing happens on your local server.

### 📖 Bible Verse Display
- **Complete KJV Bible**: All 66 books with 31,000+ verses.
- **Instant Search**: Type "John 3:16" or search by keywords.
- **Visual Navigation**: Grid-based selection for touchscreens.
- **Live Preview**: See verses before sending them to OBS.

### 🎥 OBS Integration
- **Browser Source Ready**: Direct integration with OBS Studio via `obs.html`.
- **Dynamic Sizing**: Automatically adjusts layout for long verses.
- **Professional Look**: Clean, modern lower thirds graphics.

## Quick Start

### Prerequisites
1.  [Go 1.23+](https://golang.org/dl/) installed.
2.  **GCC Compiler**: Required for audio/speech libraries (CGO).
    *   **Windows**: Install [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) or MinGW.
    *   **Linux**: `sudo apt install build-essential` or `sudo pacman -S base-devel`.
3.  **Vosk Library**:
    *   **Windows**: Download `vosk-win64-x.x.x.zip` from [Vosk Releases](https://github.com/alphacep/vosk-api/releases). Extract `libvosk.dll` to the project root (or system path).
    *   **Linux**: Download `vosk-linux-x.x.x.zip`. Extract `libvosk.so` to `/usr/local/lib` or project root (and set `LD_LIBRARY_PATH`).

### Setup Instructions

1.  **Clone the repository**
    ```bash
    git clone https://github.com/RudeCoffee/ChurchLower3rds.git
    cd ChurchLower3rds
    ```

2.  **Download the Speech Model**
    *   Download a lightweight English model (e.g., `vosk-model-small-en-us-0.15`) from the [Vosk Models page](https://alphacep.com/vosk/models).
    *   Extract the zip file.
    *   Rename the extracted folder to `model` and place it in the project root.
    *   *Structure check:* You should have a folder `ChurchLower3rds/model/` containing files like `am`, `conf`, `graph`, etc.

3.  **Install Go Dependencies**
    ```bash
    go mod tidy
    ```

4.  **Run the Server**
    ```bash
    go run main.go
    ```
    *Note: On Windows, ensure `libvosk.dll` is in the same folder or in your PATH.*

5.  **Access the Interfaces**
    *   **Control Panel**: http://localhost:8080/client.html
    *   **OBS Source**: http://localhost:8080/obs.html

## Usage Guide

### Voice Assistant
1.  Open the Control Panel (`client.html`).
2.  Click the **Gear Icon (⚙️)** in the "Voice Assistant" section to select your microphone.
3.  Click **"Start Listening"**.
4.  As you speak, the system will transcribe text and pop up a **"Suggested Verse"** card when a match is found.
5.  Click **"Preview"** or **"Display Live"** on the card.

### Manual Control
-   **Select Verse**: Click "Select Book", then "Chapter", then "Verse" in the grid.
-   **Search**: Type keywords in the "Text Search" box.
-   **Display/Hide**: Use the buttons to toggle visibility on OBS.

## Troubleshooting

-   **"Vosk model directory not found"**: Ensure you downloaded the model and renamed the folder to `model`.
-   **"Audio Input initialization failed"**: Ensure you have a working microphone and the GCC compiler is installed.
-   **"missing vosk_api.h"**: This usually means CGO cannot find the Vosk headers. Ensure you have the C libraries setup correctly if compiling manually, or rely on `go run` with the DLL/SO present.

## License
MIT License.
