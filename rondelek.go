package main

import (
	"math"

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

	rl.SetConfigFlags(rl.FlagWindowUndecorated /*| rl.FlagWindowTransparent*/)
	rl.InitWindow(int32(app.Conf.Window.Width), int32(app.Conf.Window.Height), "Rondelek TWST-1")
	rl.InitAudioDevice()
	rl.SetTargetFPS(60)

	defer rl.CloseAudioDevice()
	defer rl.CloseWindow()

	dragging := false
	var dragStart rl.Vector2

	isPointerOverUI := func() bool {
		mouse := rl.GetMousePosition()
		for _, pad := range app.Pads {
			if rl.CheckCollisionPointRec(mouse, pad.Rect) {
				return true
			}
		}
		return false
	}

	for !rl.WindowShouldClose() {

		if rl.IsMouseButtonPressed(rl.MouseLeftButton) && !isPointerOverUI() {
			dragging = true
			dragStart = rl.GetMousePosition()

		}

		if dragging && rl.IsMouseButtonDown(rl.MouseLeftButton) {
			curr := rl.GetMousePosition()
			currWin := rl.GetWindowPosition()
			dx := curr.X - dragStart.X
			dy := curr.Y - dragStart.Y
			if math.Abs(float64(dx)) >= 1 || math.Abs(float64(dy)) >= 1 {
				rl.SetWindowPosition(int(currWin.X+dx), int(currWin.Y+dy))
			}
		}

		if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
			dragging = false
		}

		rl.BeginDrawing()
		// rl.ClearBackground(rl.Black)
		// app.Grid.DrawDebug()
		ui.DrawCase()

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

		pad := ui.NewPad(app.Grid.Rectangle(startCol, endCol, startRow, endRow), padConf.Key, padConf.Label)
		pad.RegisterTransition(ui.PadStatusIdle, ui.PadStatusPressed, transitionPadIdleToPressed)
		pad.RegisterTransition(ui.PadStatusPressed, ui.PadStatusIdle, transitionPadPressedToIdle)

		app.Pads = append(app.Pads, pad)
	}
}

func createFunctionPads() {
	pad := ui.NewPad(app.Grid.Rectangle(17, 21, 6, 9), rl.KeySpace, "PLAY/REC")
	pad.RegisterTransition(ui.PadStatusIdle, ui.PadStatusPressed, func(p *ui.Pad, from, to ui.PressStatus) {
		for _, pad := range app.Pads {
			pad.ToggleMode()
		}
	})

	app.Pads = append(app.Pads, pad)
}
