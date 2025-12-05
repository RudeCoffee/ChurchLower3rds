import json
import time
import websocket
import speech_recognition as sr
import threading
import sys
import os

import argparse

# Configuration
DEFAULT_SERVER_URL = "ws://localhost:8080/ws/control"

def on_message(ws, message):
    try:
        data = json.loads(message)
        # We can handle server messages here if needed (e.g. commands to stop listening)
        # For now, we mainly ignore incoming messages as this is a producer
        pass
    except Exception as e:
        print(f"Error handling message: {e}")

def on_error(ws, error):
    print(f"WebSocket error: {error}")

def on_close(ws, close_status_code, close_msg):
    print("WebSocket connection closed")

def on_open(ws):
    print("WebSocket connection opened")

def audio_listener(ws):
    recognizer = sr.Recognizer()

    # Adjust for ambient noise
    with sr.Microphone() as source:
        print("Adjusting for ambient noise... (Please be silent)")
        recognizer.adjust_for_ambient_noise(source, duration=2)
        print("Listening...")

    # Continuous listening
    def callback(recognizer, audio):
        try:
            # Use Whisper if available locally, otherwise fallback or error
            # 'faster-whisper' is not directly supported by SpeechRecognition library out of box
            # in the standard recognize_whisper method without valid setup.
            # But we can try to use the raw audio data with a custom whisper implementation if needed.
            # However, standard SpeechRecognition has recognize_whisper() which uses OpenAI's python library.
            # If the user wants "live" local whisper, using faster-whisper directly on audio stream is better.

            # Since SpeechRecognition 'listen_in_background' provides audio data chunks,
            # we will try to transcribe them.

            # Note: recognize_whisper uses the 'openai-whisper' library (local model).
            # It might be slow if the model is large. 'base' or 'tiny' is recommended for CPU.

            print("Transcribing...")
            text = recognizer.recognize_whisper(audio, language="english", model="base")

            if text.strip():
                print(f"Heard: {text}")
                # Send to server
                msg = {
                    "type": "transcript",
                    "text": text,
                    # Context is handled by server state
                    "currentBook": "",
                    "currentChapter": 0,
                    "currentVerse": 0
                }
                ws.send(json.dumps(msg))

        except sr.UnknownValueError:
            pass # No speech detected
        except sr.RequestError as e:
            print(f"Could not request results; {e}")
        except Exception as e:
            print(f"Error during transcription: {e}")

    # Start background listening
    stop_listening = recognizer.listen_in_background(sr.Microphone(), callback, phrase_time_limit=10)

    # Keep the thread alive
    try:
        while True:
            time.sleep(0.1)
    except KeyboardInterrupt:
        stop_listening(wait_for_stop=False)

def run_audio_listener(server_url):
    # Loop to auto-reconnect
    while True:
        try:
            ws = websocket.WebSocketApp(server_url,
                                      on_open=on_open,
                                      on_message=on_message,
                                      on_error=on_error,
                                      on_close=on_close)

            # Run the websocket in a separate thread so we can run audio loop
            wst = threading.Thread(target=ws.run_forever)
            wst.daemon = True
            wst.start()

            # Wait for connection
            time.sleep(1)

            if ws.sock and ws.sock.connected:
                audio_listener(ws)
            else:
                print("Could not connect to server. Retrying in 5 seconds...")
                time.sleep(5)

        except KeyboardInterrupt:
            print("Exiting...")
            break
        except Exception as e:
            print(f"Connection error: {e}")
            time.sleep(5)

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Church Lower Thirds Audio Listener")
    parser.add_argument("--url", default=DEFAULT_SERVER_URL, help="WebSocket URL of the server (default: ws://localhost:8080/ws/control)")
    args = parser.parse_args()

    print("Starting Church Lower Thirds Audio Listener")
    print(f"Connecting to: {args.url}")
    print("This script listens to the microphone and sends transcripts to the local server.")

    run_audio_listener(args.url)
