package sampler

import (
	"fmt"
	"os"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

type SampleStatus int

const (
	SampleStatusIdle SampleStatus = iota
	SampleStatusRecording
	SampleStatusPlaying
)

type Sample struct {
	Status            SampleStatus
	FileName          string
	Buf               []float32
	sound             rl.Sound
	soundLoaded       bool
	cachedSampleCount int
	mu                sync.Mutex
}

func (t *Sample) Store() error {
	f, err := os.Create(t.FileName)
	if err != nil {
		return err
	}

	defer f.Close()

	enc := wav.NewEncoder(f, 44100, 16, 1, 1)
	defer enc.Close()

	intBuf := make([]int, len(t.Buf))
	for i, sample := range t.Buf {
		intBuf[i] = int(sample * 32768.0)
	}

	err = enc.Write(&audio.IntBuffer{Data: intBuf, Format: &audio.Format{SampleRate: 44100, NumChannels: 1}})

	return err
}

func NewSample(fileName string) *Sample {
	return &Sample{
		FileName: fileName,
		Status:   SampleStatusIdle,
		Buf:      make([]float32, 0),
	}
}

func (t *Sample) Play() error {
	if len(t.Buf) == 0 {
		return fmt.Errorf("Bufor sampla jest pusty")
	}

	if !rl.IsAudioDeviceReady() {
		return fmt.Errorf("Błąd: Dźwięk Raylib nie jest dostępny!")
	}

	if !t.soundLoaded || t.cachedSampleCount != len(t.Buf) {
		if err := t.loadSoundFromBuffer(); err != nil {
			return err
		}
	}

	rl.PlaySound(t.sound)
	t.Status = SampleStatusPlaying

	// @todo should use mutex to protect Status?
	go func() {
		for rl.IsSoundPlaying(t.sound) {

		}
		t.Status = SampleStatusIdle
	}()

	return nil
}

func (t *Sample) loadSoundFromBuffer() error {
	if len(t.Buf) == 0 {
		return fmt.Errorf("Błąd ładowania dźwięku z bufora: bufor jest pusty")
	}

	if t.soundLoaded {
		rl.UnloadSound(t.sound)
		t.soundLoaded = false
	}

	waveBytes, err := encodeBufferToWAV(t.Buf)
	if err != nil {
		return err
	}

	wave := rl.LoadWaveFromMemory(".wav", waveBytes, int32(len(waveBytes)))
	if wave.FrameCount == 0 {
		return fmt.Errorf("Nie udało się załadować audio WAV z danych %s", t.FileName)
	}

	t.sound = rl.LoadSoundFromWave(wave)
	rl.UnloadWave(wave)

	t.soundLoaded = true
	t.cachedSampleCount = len(t.Buf)

	return nil
}

func (t *Sample) Dispose() {
	if !t.soundLoaded {
		return
	}

	rl.UnloadSound(t.sound)
	t.soundLoaded = false
	t.cachedSampleCount = 0
	t.Status = SampleStatusIdle
}
