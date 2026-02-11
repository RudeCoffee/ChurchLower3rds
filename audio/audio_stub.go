//go:build stub

package audio

type DeviceInfo struct {
	Index int    `json:"index"`
	Name  string `json:"name"`
}

type AudioInput struct {
}

func NewAudioInput() (*AudioInput, error) {
	return &AudioInput{}, nil
}

func (ai *AudioInput) Close() {
}

func (ai *AudioInput) ListDevices() ([]DeviceInfo, error) {
	return []DeviceInfo{{Index: 0, Name: "Stub Device"}}, nil
}

func (ai *AudioInput) StartStream(deviceIndex int, callback func([]byte)) error {
	return nil
}

func (ai *AudioInput) StopStream() {
}
