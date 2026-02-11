//go:build !stub

package speech

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	vosk "github.com/alphacep/vosk-api/go"
)

// SpeechEngine manages the speech recognition process
type SpeechEngine struct {
	model      *vosk.VoskModel
	recognizer *vosk.VoskRecognizer
	sampleRate float64
}

// NewSpeechEngine initializes the speech engine with a model
func NewSpeechEngine(modelPath string, sampleRate float64) (*SpeechEngine, error) {
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("model path does not exist: %s", modelPath)
	}

	vosk.SetLogLevel(-1) // Suppress Vosk logs unless error

	model, err := vosk.NewModel(modelPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load model: %v", err)
	}

	recognizer, err := vosk.NewRecognizer(model, sampleRate)
	if err != nil {
		model.Free()
		return nil, fmt.Errorf("failed to create recognizer: %v", err)
	}

	return &SpeechEngine{
		model:      model,
		recognizer: recognizer,
		sampleRate: sampleRate,
	}, nil
}

// Close cleans up resources
func (se *SpeechEngine) Close() {
	if se.recognizer != nil {
		se.recognizer.Free()
	}
	if se.model != nil {
		se.model.Free()
	}
}

// ProcessAudio feeds audio data to the recognizer and returns text if a result is ready
// Returns: (partialResult, finalResult, isFinal)
func (se *SpeechEngine) ProcessAudio(data []byte) (string, string, bool) {
	// AcceptWaveform returns 1 (true) if a final result is available (silence detected)
	// Returns 0 (false) if partial result available
	if se.recognizer.AcceptWaveform(data) == 1 {
		resultJSON := se.recognizer.Result()

		// Parse result JSON: {"text": "some text"}
		var res struct {
			Text string `json:"text"`
		}
		if err := json.Unmarshal([]byte(resultJSON), &res); err == nil {
			return "", res.Text, true
		}
	} else {
		// Partial result
		partialJSON := se.recognizer.PartialResult()

		// Parse partial JSON: {"partial": "some t"}
		var res struct {
			Partial string `json:"partial"`
		}
		if err := json.Unmarshal([]byte(partialJSON), &res); err == nil {
			return res.Partial, "", false
		}
	}
	return "", "", false
}

// Reset clears the recognizer state
func (se *SpeechEngine) Reset() {
	se.recognizer.Reset()
}
