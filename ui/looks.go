package ui

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/pipejesus/rondelek/providers"
)

var (
	bgColor      = rl.Color{R: 227, G: 226, B: 221, A: 255}
	panelColor   = rl.Color{R: 246, G: 243, B: 235, A: 255}
	panelStroke  = rl.Color{R: 89, G: 88, B: 80, A: 255}
	accentOrange = rl.Color{R: 255, G: 142, B: 57, A: 255}
	accentTeal   = rl.Color{R: 63, G: 110, B: 121, A: 255}
	textDark     = rl.Color{R: 46, G: 49, B: 52, A: 255}

	padFillPlay    = rl.Color{R: 250, G: 248, B: 240, A: 255}
	padFillRecord  = rl.Color{R: 255, G: 142, B: 57, A: 205}
	padLabelPlay   = textDark
	padLabelRecord = rl.White
	padShadowColor = rl.Color{R: 200, G: 196, B: 186, A: 255}

	ledColorEmpty = rl.Color{R: 90, G: 94, B: 95, A: 255}
	ledColorFull  = rl.Color{R: 250, G: 50, B: 0, A: 255}
	shadowColor   = rl.Color{R: 160, G: 156, B: 150, A: 110}
	innerPadding  = float32(16)
	roundMd       = float32(0.032)
)

func DrawCase() {
	grid := providers.GetContainer().Grid
	panelRect := grid.Rectangle(1, grid.GetColumns(), 1, grid.GetRows())
	shadowRect := panelRect
	shadowRect.X += 16
	shadowRect.Y += 16
	rl.DrawRectangleRounded(shadowRect, roundMd, 1, shadowColor)

	rl.DrawRectangleRounded(panelRect, roundMd, 1, panelColor)
	rl.DrawRectangleRoundedLines(panelRect, roundMd, 1, panelStroke)

}

func DrawScreen() {
	grid := providers.GetContainer().Grid
	screen := grid.Rectangle(1, grid.GetColumns(), 2, 6)
	rl.DrawRectangleRec(screen, rl.Black)
}
