/*
 * Memory WriteSeeker courtesy of GPT-5.1-codex
 * Greg thanks!
 */
package sampler

import (
	"fmt"
	"io"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

type memWriteSeeker struct {
	buf []byte
	pos int64
}

func newMemWriteSeeker() *memWriteSeeker {
	return &memWriteSeeker{buf: make([]byte, 0, 4096)}
}

func (m *memWriteSeeker) ensureSize(size int) {
	if size <= len(m.buf) {
		return
	}
	newBuf := make([]byte, size)
	copy(newBuf, m.buf)
	m.buf = newBuf
}

func (m *memWriteSeeker) Write(p []byte) (int, error) {
	end := int(m.pos) + len(p)
	m.ensureSize(end)
	copy(m.buf[int(m.pos):end], p)
	m.pos += int64(len(p))
	return len(p), nil
}

func (m *memWriteSeeker) Seek(offset int64, whence int) (int64, error) {
	var base int64
	switch whence {
	case io.SeekStart:
		base = 0
	case io.SeekCurrent:
		base = m.pos
	case io.SeekEnd:
		base = int64(len(m.buf))
	default:
		return 0, fmt.Errorf("invalid whence: %d", whence)
	}

	newPos := base + offset
	if newPos < 0 {
		return 0, fmt.Errorf("negative position")
	}

	m.pos = newPos
	m.ensureSize(int(m.pos))
	return m.pos, nil
}

func (m *memWriteSeeker) Bytes() []byte {
	return append([]byte(nil), m.buf[:m.pos]...)
}

func encodeBufferToWAV(samples []float32) ([]byte, error) {
	ws := newMemWriteSeeker()
	enc := wav.NewEncoder(ws, 44100, 16, 1, 1)

	intBuf := make([]int, len(samples))
	for i, sample := range samples {
		val := int(sample * 32767.0)
		if val > 32767 {
			val = 32767
		}
		if val < -32768 {
			val = -32768
		}
		intBuf[i] = val
	}

	if err := enc.Write(&audio.IntBuffer{Data: intBuf, Format: &audio.Format{SampleRate: 44100, NumChannels: 1}}); err != nil {
		return nil, err
	}

	if err := enc.Close(); err != nil {
		return nil, err
	}

	return ws.Bytes(), nil
}
