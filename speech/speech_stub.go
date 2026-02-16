//go:build stub

package speech

// SpeechEngine manages the speech recognition process (STUB)
type SpeechEngine struct {
	sampleRate float64
}

// NewSpeechEngine initializes the speech engine (STUB)
func NewSpeechEngine(modelPath string, sampleRate float64) (*SpeechEngine, error) {
	return &SpeechEngine{sampleRate: sampleRate}, nil
}

// Close cleans up resources (STUB)
func (se *SpeechEngine) Close() {
}

// ProcessAudio feeds audio data to the recognizer (STUB)
func (se *SpeechEngine) ProcessAudio(data []byte) (string, string, bool) {
	return "", "", false
}

// Reset clears the recognizer state (STUB)
func (se *SpeechEngine) Reset() {
}
