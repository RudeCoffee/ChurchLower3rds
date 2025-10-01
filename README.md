# Church Lower 3rds

A real-time Bible verse and speaker name display system for OBS (Open Broadcaster Software) designed for church services and live streaming.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)
![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20macOS%20%7C%20Linux-lightgrey.svg)

**[Note: The screenshot below is outdated. Please replace it with the new `new_ui.png` screenshot.]**
<img width="1986" height="1288" alt="image" src="https://github.com/user-attachments/assets/e8ce40da-fd31-430f-9324-0c65e2b1f6d6" />


## Features

### 🎤 Live Sermon Scripture Suggestions
- **Real-time Transcription**: Listens to the live sermon audio and transcribes it in real-time.
- **Automatic Suggestions**: Intelligently suggests relevant Bible verses based on the spoken words.
- **Keyword-Based and Direct References**: Detects both direct scripture mentions (e.g., "John 3:16") and keyword-based themes (e.g., "love," "faith").
- **One-Click Display**: Click on a suggestion to instantly display the verse on the broadcast.

### 📖 Bible Verse Display
- **Complete KJV Bible**: All 66 books with 31,000+ verses.
- **Smart Quick Input**: Type "Genesis 1:16" for instant verse selection.
- **Real-time Autocomplete**: Intelligent book name suggestions with validation.
- **Auto-Advance Toggle**: Automatically switch to the next verse when the speaker reads it.
- **Live & Next Verse Preview**: See the current and upcoming verses before displaying them.

### 🎤 Speaker Management
- **Speaker Names Display**: Professional lower third graphics for speaker identification.
- **Autocomplete System**: Quick speaker selection with predefined names.
- **Customizable List**: Easily manage the list of speakers in the `speakers.txt` file.

### ⏱️ Countdown Timer
- **Live Countdown**: Display a real-time countdown for service starts or events.
- **"Starting Soon" Button**: Manually trigger a "Starting Soon" message at any time.
- **Preset Times**: Quickly set the countdown to common service start times.

### 🎥 OBS Integration
- **Browser Source Ready**: Direct integration with OBS Studio.
- **Fade Animations**: Smooth transitions for a professional appearance.
- **Dynamic Sizing**: Automatically adjusts for long verses.

### 🎨 User Interface
- **Dark/Light Themes**: Toggle between visual modes for comfort.
- **Responsive Design**: Works on desktop and tablet devices.
- **Cleaner Layout**: A redesigned interface places key controls side-by-side for easier access.
- **Keyboard Shortcuts**: Efficient navigation with arrow keys and Enter.

## Quick Start

### Prerequisites
- [Go 1.21+](https://golang.org/dl/) installed on your system
- [OBS Studio](https://obsproject.com/) for broadcasting
- A modern web browser that supports the Web Speech API (Chrome, Edge) for the sermon suggestion feature.

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/RudeCoffee/ChurchLower3rds.git
   cd ChurchLower3rds
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Start the server**
   ```bash
   go run main.go
   ```

4. **Access the interfaces**
   - **Control Interface**: http://localhost:8080/client.html
   - **OBS Lower Thirds Source**: http://localhost:8080/obs.html
   - **OBS Countdown Source**: http://localhost:8080/countdown.html

## Usage Guide

### Live Sermon Suggestions
1.  Click the **"🎤 Start Listening"** button. Your browser will ask for microphone permission.
2.  As the sermon is preached, the system will transcribe the audio and display relevant scripture suggestions.
3.  Click on any suggestion to instantly preview it and display it on the OBS output.

### Quick Verse Selection
- **Smart Input**: Type a reference like `Genesis 1:16` and press Enter.
- **Auto-Advance**: Enable the "Auto-Advance" toggle switch. The system will listen to the sermon audio and automatically advance to the next verse when it detects that the current verse has been read.
- **Traditional Dropdowns**: Use the Book, Chapter, and Verse dropdowns for manual selection.

### Countdown Timer
- **Set a Time**: Use one of the preset buttons or enter a custom time.
- **Start Countdown**: Click "Start Countdown" to begin the timer on the OBS output.
- **Starting Soon**: Click the "Starting Soon" button to immediately display the "Starting soon" message, overriding the timer.

## File Structure
```
ChurchLower3rds/
├── main.go                 # Go server with WebSocket handling
├── client.html            # Control interface for operators
├── obs.html               # OBS browser source display
├── countdown.html         # OBS browser source for countdown timer
├── kjv.json               # Complete KJV Bible database
├── speakers.txt           # Customizable speaker names list
└── ...
```

## Contributing
Contributions are welcome! Please fork the repository, create a feature branch, and open a pull request.

## License
This project is licensed under the MIT License. See the [LICENSES.md](LICENSES.md) file for details.