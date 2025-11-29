package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/octoper/go-ray"
	_ "github.com/octoper/go-ray"
	"github.com/pipejesus/rondelek/sampler"
	ui "github.com/pipejesus/rondelek/ui"
)

type App struct {
	Sampler *sampler.Sampler
	Pads    []*ui.Pad
	Grid    *ui.Grid
	Conf    *Config
}

var app *App

func init() {
	conf := NewConfig()
	conf.Load()

	app = &App{
		Sampler: sampler.NewSampler(),
		Pads:    []*ui.Pad{},
		Grid: ui.NewGrid(
			conf.Window.Width, conf.Window.Height,
			conf.Layout.Columns, conf.Layout.Rows,
			conf.Layout.PaddingX, conf.Layout.PaddingY,
		),
		Conf: conf,
	}

	app.Sampler.Init()
}

func main() {
	defer app.Sampler.Quit()

	createMainPads()
	createFunctionPads()

	rl.InitWindow(int32(app.Conf.Window.Width), int32(app.Conf.Window.Height), "Rondelek TWST-1")
	rl.InitAudioDevice()
	rl.SetTargetFPS(60)
	defer rl.CloseAudioDevice()
	defer rl.CloseWindow()

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		app.Grid.DrawDebug()

		for _, pad := range app.Pads {
			pad.Update()
		}

		rl.EndDrawing()
	}
}

func createMainPads() {

	const (
		realStartRow = 6
		padsPerAxis  = 4
		padSizeX     = 3
		padSizeY     = 4
		padGap       = 1
	)

	for rowIdx := range padsPerAxis {
		startRow := realStartRow + rowIdx*(padSizeY+padGap)
		endRow := startRow + padSizeY - 1

		for colIdx := range padsPerAxis {
			startCol := 1 + colIdx*(padSizeX+padGap)
			endCol := startCol + padSizeX - 1

			positionAsString := fmt.Sprintf("%d Col start: %d Row start: %d", colIdx, startCol, startRow)
			ray.Ray(positionAsString)

			pad := ui.NewPad(app.Grid.Rectangle(startCol, endCol, startRow, endRow))
			pad.RegisterTransition(ui.PadStatusIdle, ui.PadStatusPressed, func(p *ui.Pad, from, to ui.PressStatus) {
				if p.Mode == ui.ModeRecord {
					fmt.Println("Recording sound!")
					app.Sampler.Record()
					return
				}

				if !p.HasSample() {
					fmt.Println("No sample assigned to this pad!")
					return
				}

				if err := app.Sampler.Samples[p.SampleIdx].Play(); err != nil {
					fmt.Println("playback error:", err)
				}
			})
			pad.RegisterTransition(ui.PadStatusPressed, ui.PadStatusIdle, func(p *ui.Pad, from, to ui.PressStatus) {
				if p.Mode == ui.ModeRecord {
					fmt.Println("Stopping sound!")
					sampleIdx := app.Sampler.StopRecording()
					p.SetSample(sampleIdx)
					app.Sampler.FreshSample()
				}
			})

			app.Pads = append(app.Pads, pad)
		}
	}

}

func createFunctionPads() {
	pad := ui.NewPad(app.Grid.Rectangle(17, 19, 6, 9))
	pad.RegisterTransition(ui.PadStatusIdle, ui.PadStatusPressed, func(p *ui.Pad, from, to ui.PressStatus) {
		for _, pad := range app.Pads {
			pad.ToggleMode()
		}
	})

	app.Pads = append(app.Pads, pad)
}
