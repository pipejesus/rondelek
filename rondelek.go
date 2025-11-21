package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/pipejesus/rondelek/sampler"
)

type App struct {
	Sampler *sampler.Sampler
}

var app *App

func init() {
	app = &App{
		Sampler: sampler.NewSampler(),
	}

	app.Sampler.Init()
}

func main() {
	defer app.Sampler.Quit()

	rl.InitWindow(800, 600, "Rondelek")
	rl.SetTargetFPS(60)
	defer rl.CloseWindow()

	pad := NewPad(rl.Rectangle{X: 10.0, Y: 10.0, Width: 100.0, Height: 100.0})
	pad.RegisterTransition(PadStatusIdle, PadStatusPressed, func(p *Pad, from, to PadStatus) {
		fmt.Println("Recording sound!")
		app.Sampler.Record()
	})
	pad.RegisterTransition(PadStatusPressed, PadStatusIdle, func(p *Pad, from, to PadStatus) {
		fmt.Println("Stopping sound!")
		app.Sampler.Stop()
		app.Sampler.FreshSample()
	})

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		pad.Update()

		rl.EndDrawing()
	}
}
