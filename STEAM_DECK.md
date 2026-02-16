# Steam Deck Setup Guide for Church Lower Thirds with Voice Assistant

This guide will help you set up and run the application with the new Voice Assistant features on your Steam Deck. The Steam Deck runs SteamOS (Arch Linux), so we need to set up a few development tools and the Vosk speech recognition library.

## Prerequisites

1.  **Switch to Desktop Mode**: Press the Steam button > Power > Switch to Desktop.
2.  **Open Konsole**: Click the Steam Deck icon in the bottom left, search for "Konsole", and open it.
3.  **Set a User Password** (if you haven't already):
    *   Type `passwd` in Konsole and press Enter.
    *   Set a password (it won't show as you type). This is needed for `sudo` commands.

## Step 1: Prepare the System

By default, the Steam Deck file system is read-only. To install the necessary build tools (like `gcc`), we need to disable this temporarily.

1.  **Disable Read-Only Mode**:
    ```bash
    sudo steamos-readonly disable
    ```

2.  **Initialize Pacman Keys**:
    ```bash
    sudo pacman-key --init
    sudo pacman-key --populate archlinux
    ```

3.  **Install Base Development Tools**:
    We need `gcc` to compile the Go code with CGO (required for the microphone and speech libraries).
    ```bash
    sudo pacman -S --noconfirm base-devel git
    ```

    *Note: If you prefer not to modify your system partition, you can try to build the binary on another Linux machine and copy it over, but you will still need the `libvosk.so` library on the Deck.*

## Step 2: Download Vosk Library and Model

The application uses **Vosk** for offline speech-to-text. We need to download the shared library and the speech model.

1.  **Create a `vosk` directory** inside the project folder:
    ```bash
    mkdir -p vosk
    cd vosk
    ```

2.  **Download and Extract the Vosk Library**:
    ```bash
    # Download 64-bit Linux library
    wget https://github.com/alphacep/vosk-api/releases/download/v0.3.45/vosk-linux-x86_64-0.3.45.zip
    unzip vosk-linux-x86_64-0.3.45.zip
    mv vosk-linux-x86_64-0.3.45 lib
    rm vosk-linux-x86_64-0.3.45.zip
    ```

3.  **Download and Extract a Speech Model**:
    We recommend the small US English model for performance.
    ```bash
    # Download model
    wget https://alphacephei.com/vosk/models/vosk-model-small-en-us-0.15.zip
    unzip vosk-model-small-en-us-0.15.zip
    mv vosk-model-small-en-us-0.15 model
    rm vosk-model-small-en-us-0.15.zip
    ```

    *Your project structure inside `vosk/` should now look like:*
    ```
    vosk/
    ├── lib/
    │   └── libvosk.so
    │   └── vosk_api.h
    └── model/
        └── (model files)
    ```

    Go back to the project root:
    ```bash
    cd ..
    ```

## Step 3: Build and Run

Now we need to tell Go where to find the Vosk headers and library during the build.

1.  **Build the Application**:
    Copy and paste this entire block to build:

    ```bash
    # Enable CGO
    export CGO_ENABLED=1

    # Point to the header files (vosk_api.h)
    export CGO_CFLAGS="-I$PWD/vosk/lib"

    # Point to the library files (libvosk.so)
    export CGO_LDFLAGS="-L$PWD/vosk/lib -lvosk"

    # Build the binary
    go build -o church-lower-thirds .
    ```

2.  **Run the Application**:
    When running, the application needs to know where `libvosk.so` is located.

    ```bash
    # Set library path and run
    export LD_LIBRARY_PATH="$PWD/vosk/lib:$LD_LIBRARY_PATH"
    ./church-lower-thirds
    ```

## Troubleshooting

*   **"error while loading shared libraries: libvosk.so"**:
    *   Make sure you ran the `export LD_LIBRARY_PATH...` command before running the app.
*   **"gcc: command not found"**:
    *   Ensure you ran `sudo pacman -S base-devel`.
*   **Microphone not working**:
    *   Check the Steam Deck audio settings in Desktop Mode.
    *   Ensure the application has permission to access the microphone (though on Linux/SteamOS this is usually open).
    *   Verify the correct input device is selected in the web interface (click the settings gear icon).

## Quick Start Script (Optional)

You can create a `run_deck.sh` file to make this easier in the future:

```bash
#!/bin/bash
export CGO_ENABLED=1
export CGO_CFLAGS="-I$PWD/vosk/lib"
export CGO_LDFLAGS="-L$PWD/vosk/lib -lvosk"
export LD_LIBRARY_PATH="$PWD/vosk/lib:$LD_LIBRARY_PATH"

# Check if built, if not build it
if [ ! -f ./church-lower-thirds ]; then
    echo "Building..."
    go build -o church-lower-thirds .
fi

echo "Starting Church Lower Thirds..."
./church-lower-thirds
```

Make it executable: `chmod +x run_deck.sh`
