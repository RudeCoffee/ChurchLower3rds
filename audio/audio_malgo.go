//go:build !stub

package audio

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"unsafe"

	"github.com/gen2brain/malgo"
)

// DeviceInfo holds information about an audio input device
type DeviceInfo struct {
	Index int    `json:"index"` // Use int index for simplicity
	Name  string `json:"name"`
}

// AudioInput manages audio capture using malgo (miniaudio)
type AudioInput struct {
	ctx        *malgo.AllocatedContext
	device     *malgo.Device
	streamOpen bool
	mu         sync.Mutex
	callback   func([]byte)
	pinner     *runtime.Pinner // Keep pinner to unpin later
}

// NewAudioInput creates a new audio input manager
func NewAudioInput() (*AudioInput, error) {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		log.Printf("Malgo Log: %v\n", message)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init malgo context: %v", err)
	}

	return &AudioInput{
		ctx: ctx,
	}, nil
}

// Close cleans up resources
func (ai *AudioInput) Close() {
	ai.mu.Lock()
	defer ai.mu.Unlock()

	if ai.streamOpen && ai.device != nil {
		ai.device.Uninit()
	}
	if ai.pinner != nil {
		ai.pinner.Unpin()
		ai.pinner = nil
	}
	if ai.ctx != nil {
		ai.ctx.Free() // Changed Uninit to Free for allocated context
	}
}

// ListDevices returns a list of available input devices
func (ai *AudioInput) ListDevices() ([]DeviceInfo, error) {
	devices, err := ai.ctx.Devices(malgo.Capture)
	if err != nil {
		return nil, fmt.Errorf("failed to list devices: %v", err)
	}

	var deviceList []DeviceInfo
	for i, info := range devices {
		deviceList = append(deviceList, DeviceInfo{
			Index: i,
			Name:  info.Name(),
		})
	}
	return deviceList, nil
}

// StartStream starts capturing audio from the specified device index
// callback is called with audio chunks (PCM 16-bit, 16kHz, Mono)
func (ai *AudioInput) StartStream(deviceIndex int, callback func([]byte)) error {
	ai.mu.Lock()
	defer ai.mu.Unlock()

	// Stop existing stream if running
	if ai.streamOpen {
		ai.device.Uninit()
		if ai.pinner != nil {
			ai.pinner.Unpin()
			ai.pinner = nil
		}
		ai.streamOpen = false
	}

	ai.callback = callback

	devices, err := ai.ctx.Devices(malgo.Capture)
	if err != nil {
		return fmt.Errorf("failed to list devices for stream: %v", err)
	}

	if deviceIndex < 0 || deviceIndex >= len(devices) {
		return fmt.Errorf("invalid device index: %d", deviceIndex)
	}

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	deviceConfig.Capture.Format = malgo.FormatS16
	deviceConfig.Capture.Channels = 1
	deviceConfig.SampleRate = 16000
	deviceConfig.Alsa.NoMMap = 1

	// Set specific device ID
	if deviceIndex == 0 {
		deviceConfig.Capture.DeviceID = nil // Default device
	} else {
		// Use runtime.Pinner to ensure the Go pointer passed to CGO is pinned
		id := devices[deviceIndex].ID
		ai.pinner = new(runtime.Pinner)
		ai.pinner.Pin(&id)
		deviceConfig.Capture.DeviceID = unsafe.Pointer(&id)
		log.Printf("Pinned device ID for device index %d", deviceIndex)
	}

	// Callback function for malgo
	onRecv := func(pOutputSample, pInputSamples []byte, framecount uint32) {
		// Create a copy to avoid race conditions
		data := make([]byte, len(pInputSamples))
		copy(data, pInputSamples)

		if ai.callback != nil {
			ai.callback(data)
		}
	}

	deviceCallbacks := malgo.DeviceCallbacks{
		Data: onRecv,
	}

	device, err := malgo.InitDevice(ai.ctx.Context, deviceConfig, deviceCallbacks)
	if err != nil {
		if ai.pinner != nil {
			ai.pinner.Unpin()
			ai.pinner = nil
		}
		return fmt.Errorf("failed to init device: %v", err)
	}

	err = device.Start()
	if err != nil {
		device.Uninit()
		if ai.pinner != nil {
			ai.pinner.Unpin()
			ai.pinner = nil
		}
		return fmt.Errorf("failed to start device: %v", err)
	}

	ai.device = device
	ai.streamOpen = true
	log.Printf("Started audio stream on device: %s", devices[deviceIndex].Name())

	return nil
}

// StopStream stops the current audio stream
func (ai *AudioInput) StopStream() {
	ai.mu.Lock()
	defer ai.mu.Unlock()

	if ai.streamOpen && ai.device != nil {
		ai.device.Stop()
		ai.device.Uninit()
		ai.streamOpen = false
		log.Println("Stopped audio stream")
	}

	if ai.pinner != nil {
		ai.pinner.Unpin()
		ai.pinner = nil
	}
}
