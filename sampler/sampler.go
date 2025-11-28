package sampler

import (
	"time"

	"github.com/gordonklaus/portaudio"
)

type Sampler struct {
	Stream    *portaudio.Stream
	RecSample *Sample
	Samples   []*Sample
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
}

func (s *Sampler) Quit() {
	if s.RecSample != nil {
		s.RecSample.Dispose()
	}

	if err := s.Stream.Close(); err != nil {
		panic(err)
	}

	if err := portaudio.Terminate(); err != nil {
		panic(err)
	}
}

func (s *Sampler) Record() {
	if err := s.Stream.Start(); err != nil {
		panic(err)
	}
}

func (s *Sampler) StopRecording() int {
	err := s.Stream.Stop()
	if err != nil {
		panic(err)
	}

	s.Samples = append(s.Samples, s.RecSample)
	return len(s.Samples) - 1
}

func (s *Sampler) CaptureAudio(in []float32) {
	s.RecSample.Buf = append(s.RecSample.Buf, in...)
}

func (s *Sampler) FreshSample() {
	if s.RecSample != nil {
		_ = s.RecSample.Store()
		//	s.RecSample.Dispose()
	}

	s.RecSample = NewSample(NewSampleFileName())
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
