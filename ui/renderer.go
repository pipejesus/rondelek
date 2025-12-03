package ui

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// A new struct for rendering passive UI elements
// subject to change, hah, but still better than what was before.
type Renderer struct {
	grid *Grid
}

func NewRenderer(grid *Grid) *Renderer {
	return &Renderer{
		grid: grid,
	}
}

func (r *Renderer) DrawCase() {
	panelRect := r.grid.Rectangle(1, r.grid.GetColumns(), 1, r.grid.GetRows())

	shadowRect := panelRect
	shadowRect.X += 16
	shadowRect.Y += 16
	rl.DrawRectangleRounded(shadowRect, Theme.RoundMd, 1, Theme.Shadow)
	rl.DrawRectangleRounded(panelRect, Theme.RoundMd, 1, Theme.Panel)
	rl.DrawRectangleRoundedLines(panelRect, Theme.RoundMd, 1, Theme.PanelStroke)
}

func (r *Renderer) DrawScreen() {
	screen := r.grid.Rectangle(1, r.grid.GetColumns(), 2, 6)
	rl.DrawRectangleRec(screen, rl.Black)
}

func (r *Renderer) DrawAll() {
	r.DrawCase()
	r.DrawScreen()
}

func (r *Renderer) Grid() *Grid {
	return r.grid
}
