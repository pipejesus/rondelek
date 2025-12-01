package main

import (
	"fmt"

	ui "github.com/pipejesus/rondelek/ui"
)

func transitionPadIdleToPressed(p *ui.Pad, from, to ui.PressStatus) {
	if p.Mode == ui.ModeRecord {
		app.Sampler.Record()
		return
	}

	if !p.HasSample() {
		return
	}

	if err := app.Sampler.Samples[p.SampleIdx].Play(); err != nil {
		fmt.Println("playback error:", err)
	}
}

func transitionPadPressedToIdle(p *ui.Pad, from, to ui.PressStatus) {
	if p.Mode == ui.ModeRecord {
		fmt.Println("Stopping sound!")
		sampleIdx := app.Sampler.StopRecording()
		p.SetSample(sampleIdx)
		app.Sampler.SaveCurrentSample()
	}
}
