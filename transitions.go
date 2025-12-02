package main

import (
	"fmt"

	"github.com/pipejesus/rondelek/providers"
	ui "github.com/pipejesus/rondelek/ui"
)

func transitionPadIdleToPressed(p *ui.Pad, from, to ui.PressStatus) {
	s := providers.GetContainer().Sampler
	if p.Mode == ui.ModeRecord {
		s.Record()
		return
	}

	if !p.HasSample() {
		return
	}

	if err := s.PlaySample(p.SampleIdx); err != nil {
		fmt.Println("playback error:", err)
	}
}

func transitionPadPressedToIdle(p *ui.Pad, from, to ui.PressStatus) {
	s := providers.GetContainer().Sampler
	if p.Mode == ui.ModeRecord {
		fmt.Println("Stopping sound!")
		sampleIdx := s.StopRecording()
		p.SetSample(sampleIdx)
		s.SaveCurrentSample()
	}
}
