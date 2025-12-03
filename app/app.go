package app

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/pipejesus/rondelek/config"
	"github.com/pipejesus/rondelek/experiments"
	"github.com/pipejesus/rondelek/sampler"
	"github.com/pipejesus/rondelek/ui"
)

type App struct {
	config   *config.Config
	grid     *ui.Grid
	renderer *ui.Renderer
	sampler  *sampler.Sampler
	pads     []*ui.Pad
}

func New() (*App, error) {
	conf := config.New()
	if err := conf.Load(); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	grid := ui.NewGrid(
		conf.Window.Width, conf.Window.Height,
		conf.Layout.Columns, conf.Layout.Rows,
		conf.Layout.PaddingX, conf.Layout.PaddingY,
	)

	renderer := ui.NewRenderer(grid)

	samp := sampler.NewSampler()
	samp.Init()

	app := &App{
		config:   conf,
		grid:     grid,
		renderer: renderer,
		sampler:  samp,
		pads:     make([]*ui.Pad, 0),
	}

	return app, nil
}

func (a *App) Run() error {
	defer a.shutdown()

	a.createMainPads()
	a.createFunctionPads()

	experiments.SetupShapedWindow(a.pads)
	rl.InitWindow(int32(a.config.Window.Width), int32(a.config.Window.Height), "Rondelek TWST-1")

	defer rl.CloseWindow()

	experiments.SetupFullScreenWindow()

	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		a.update()
		a.draw()
	}

	return nil
}

func (a *App) update() {
	experiments.HandleWindowDragging()
	for _, pad := range a.pads {
		pad.Update()
	}

}

func (a *App) draw() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.White)

	a.renderer.DrawAll()
	for _, pad := range a.pads {
		pad.Draw()
	}

	rl.EndDrawing()
}

func (a *App) shutdown() {
	a.sampler.Quit()
}

func (a *App) createMainPads() {
	for _, padConf := range a.config.Pads {
		if padConf.Type != "sample-pad" {
			continue
		}

		startCol := padConf.PadPosition.Col
		endCol := startCol + padConf.PadSize.Width - 1
		startRow := padConf.PadPosition.Row
		endRow := startRow + padConf.PadSize.Height - 1

		pad := ui.NewPad(
			a.grid.Rectangle(startCol, endCol, startRow, endRow),
			padConf.Key,
			padConf.Label,
		)

		pad.RegisterTransition(ui.PadStatusIdle, ui.PadStatusPressed, a.onPadPressed)
		pad.RegisterTransition(ui.PadStatusPressed, ui.PadStatusIdle, a.onPadReleased)

		a.pads = append(a.pads, pad)
	}
}

func (a *App) createFunctionPads() {
	pad := ui.NewPad(a.grid.Rectangle(22, 23, 9, 9), rl.KeySpace, "*")
	pad.RegisterTransition(ui.PadStatusIdle, ui.PadStatusPressed, func(p *ui.Pad, from, to ui.PressStatus) {
		for _, pad := range a.pads {
			pad.ToggleMode()
		}
	})
	a.pads = append(a.pads, pad)
}

func (a *App) onPadPressed(p *ui.Pad, from, to ui.PressStatus) {
	if p.Mode == ui.ModeRecord {
		a.sampler.Record()
		return
	}

	if !p.HasSample() {
		return
	}

	if err := a.sampler.PlaySample(p.SampleIdx); err != nil {
		fmt.Println("Błąd odtwarzania sampla :", err)
	}
}

func (a *App) onPadReleased(p *ui.Pad, from, to ui.PressStatus) {
	if p.Mode == ui.ModeRecord {
		sampleIdx := a.sampler.StopRecording()
		p.SetSample(sampleIdx)
		a.sampler.SaveCurrentSample()
	}
}

func (a *App) Pads() []*ui.Pad {
	return a.pads
}

func (a *App) Grid() *ui.Grid {
	return a.grid
}
