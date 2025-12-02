package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	_ "github.com/octoper/go-ray"
	"github.com/pipejesus/rondelek/experiments"
	"github.com/pipejesus/rondelek/providers"
	"github.com/pipejesus/rondelek/sampler"
	ui "github.com/pipejesus/rondelek/ui"
)

type App struct {
	Pads []*ui.Pad
	Conf *Config
}

var app *App

func init() {
	conf := NewConfig()
	conf.Load()

	app = &App{
		Pads: []*ui.Pad{},
		Conf: conf,
	}

	sampler := sampler.NewSampler()
	sampler.Init()

	providers.SetContainer(&providers.Container{
		Grid: ui.NewGrid(
			conf.Window.Width, conf.Window.Height,
			conf.Layout.Columns, conf.Layout.Rows,
			conf.Layout.PaddingX, conf.Layout.PaddingY,
		),
		Sampler: sampler,
	})

}

func main() {
	sampler := providers.GetContainer().Sampler

	defer sampler.Quit()

	createMainPads()
	createFunctionPads()

	experiments.SetupShapedWindow(app.Pads)

	rl.InitWindow(int32(app.Conf.Window.Width), int32(app.Conf.Window.Height), "Rondelek TWST-1")
	experiments.SetupFullScreenWindow()
	rl.InitAudioDevice()
	rl.SetTargetFPS(60)

	defer rl.CloseAudioDevice()
	defer rl.CloseWindow()

	for !rl.WindowShouldClose() {
		experiments.HandleWindowDragging()

		rl.BeginDrawing()
		rl.ClearBackground(rl.White)

		ui.DrawCase()
		ui.DrawScreen()

		for _, pad := range app.Pads {
			pad.Update()
		}

		rl.EndDrawing()
	}
}

func createMainPads() {
	grid := providers.GetContainer().Grid

	for _, padConf := range app.Conf.Pads {
		if padConf.Type != "sample-pad" {
			continue
		}

		startCol := padConf.PadPosition.Col
		endCol := startCol + padConf.PadSize.Width - 1
		startRow := padConf.PadPosition.Row
		endRow := startRow + padConf.PadSize.Height - 1

		pad := ui.NewPad(grid.Rectangle(startCol, endCol, startRow, endRow), padConf.Key, padConf.Label)
		pad.RegisterTransition(ui.PadStatusIdle, ui.PadStatusPressed, transitionPadIdleToPressed)
		pad.RegisterTransition(ui.PadStatusPressed, ui.PadStatusIdle, transitionPadPressedToIdle)

		app.Pads = append(app.Pads, pad)
	}
}

func createFunctionPads() {
	grid := providers.GetContainer().Grid

	pad := ui.NewPad(grid.Rectangle(22, 23, 9, 9), rl.KeySpace, "*")
	pad.RegisterTransition(ui.PadStatusIdle, ui.PadStatusPressed, func(p *ui.Pad, from, to ui.PressStatus) {
		for _, pad := range app.Pads {
			pad.ToggleMode()
		}
	})

	app.Pads = append(app.Pads, pad)
}
