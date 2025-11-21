package sampler

import (
	"os"

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
	Status   SampleStatus
	FileName string
	Buf      []float32
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
