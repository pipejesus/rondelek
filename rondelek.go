package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
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
	for _, padConf := range app.Conf.Pads {
		if padConf.Type != "sample-pad" {
			continue
		}

		startCol := padConf.PadPosition.Col
		endCol := startCol + padConf.PadSize.Width - 1
		startRow := padConf.PadPosition.Row
		endRow := startRow + padConf.PadSize.Height - 1

		pad := ui.NewPad(app.Grid.Rectangle(startCol, endCol, startRow, endRow))
		pad.RegisterTransition(ui.PadStatusIdle, ui.PadStatusPressed, transitionPadIdleToPressed)
		pad.RegisterTransition(ui.PadStatusPressed, ui.PadStatusIdle, transitionPadPressedToIdle)

		app.Pads = append(app.Pads, pad)
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
