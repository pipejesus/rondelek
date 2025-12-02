package sampler

import (
	"fmt"
	"time"

	"github.com/gordonklaus/portaudio"
)

type Sampler struct {
	Stream      *portaudio.Stream
	RecSample   *Sample
	Samples     []*Sample
	isRecording bool
}

func (s *Sampler) Init() {
	err := portaudio.Initialize()
	if err != nil {
		panic(err)
	}

	stream, err := portaudio.OpenDefaultStream(1, 0, 44100, 0, s.CaptureAudio)
	if err != nil {
		panic(err)
	}

	s.Stream = stream

	// Start the stream once during initialization; keep it running so we don't
	// have to warm it up for the first recording.
	if err := s.Stream.Start(); err != nil {
		panic(err)
	}
}

func (s *Sampler) Quit() {
	if s.RecSample != nil {
		s.RecSample.Dispose()
	}

	for _, sample := range s.Samples {
		sample.Dispose()
	}

	if s.Stream != nil {
		_ = s.Stream.Stop()
		if err := s.Stream.Close(); err != nil {
			panic(err)
		}
	}

	if err := portaudio.Terminate(); err != nil {
		panic(err)
	}
}

func (s *Sampler) PlaySample(idx int) error {
	if idx < 0 || idx >= len(s.Samples) {
		return fmt.Errorf("invalid sample index: %d", idx)
	}
	return s.Samples[idx].Play()
}

func (s *Sampler) Record() {
	s.RecSample = NewSample(NewSampleFileName())
	s.isRecording = true
}

func (s *Sampler) StopRecording() int {
	s.isRecording = false
	s.Samples = append(s.Samples, s.RecSample)
	return len(s.Samples) - 1
}

func (s *Sampler) CaptureAudio(in []float32) {
	if !s.isRecording || s.RecSample == nil {
		return
	}

	s.RecSample.mu.Lock()
	s.RecSample.Buf = append(s.RecSample.Buf, in...)
	s.RecSample.mu.Unlock()
}

func (s *Sampler) SaveCurrentSample() {
	if s.RecSample != nil {
		_ = s.RecSample.Store()
		s.RecSample.Dispose()
	}
}

func NewSampleFileName() string {
	return time.Now().Format("2006-01-02_15-04-05") + ".wav"
}

func NewSampler() *Sampler {
	return &Sampler{
		RecSample: NewSample(NewSampleFileName()),
		Samples:   []*Sample{},
	}
}
