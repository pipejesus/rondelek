package ui

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Grid struct {
	Width    float32 // Window width in pixels
	Height   float32 // Window height in pixels
	Columns  int     // Number of columns in the grid
	Rows     int     // Number of rows in the grid
	PaddingX float32 // Window padding
	PaddingY float32 // Window padding
}

func NewGrid(width, height float32, columns, rows int, paddingX, paddingY float32) *Grid {
	return &Grid{
		Width:    width,
		Height:   height,
		Columns:  columns,
		Rows:     rows,
		PaddingX: paddingX,
		PaddingY: paddingY,
	}
}

func (g *Grid) GetColumns() int {
	return g.Columns
}

func (g *Grid) GetRows() int {
	return g.Rows
}

func (g *Grid) DrawDebug() {
	innerWidth := g.Width - 2*g.PaddingX
	innerHeight := g.Height - 2*g.PaddingY
	cellWidth := innerWidth / float32(g.Columns)
	cellHeight := innerHeight / float32(g.Rows)

	left := g.PaddingX
	top := g.PaddingY
	right := left + innerWidth
	bottom := top + innerHeight

	const dashSize = 8.0
	const thickness = 1.0
	color := rl.Fade(rl.Green, 0.5)

	for col := 0; col <= g.Columns; col++ {
		x := left + float32(col)*cellWidth
		start := rl.Vector2{X: x, Y: top}
		end := rl.Vector2{X: x, Y: bottom}
		drawDashedLine(start, end, dashSize, thickness, color)
	}

	for row := 0; row <= g.Rows; row++ {
		y := top + float32(row)*cellHeight
		start := rl.Vector2{X: left, Y: y}
		end := rl.Vector2{X: right, Y: y}
		drawDashedLine(start, end, dashSize, thickness, color)
	}
}

func (g *Grid) Rectangle(startColumn, endColumn, startRow, endRow int) rl.Rectangle {
	startColumn = clampInt(startColumn, 1, g.Columns)
	endColumn = clampInt(endColumn, 1, g.Columns)
	startRow = clampInt(startRow, 1, g.Rows)
	endRow = clampInt(endRow, 1, g.Rows)

	if startColumn > endColumn {
		startColumn, endColumn = endColumn, startColumn
	}
	if startRow > endRow {
		startRow, endRow = endRow, startRow
	}

	innerWidth := g.Width - 2*g.PaddingX
	innerHeight := g.Height - 2*g.PaddingY
	cellWidth := innerWidth / float32(g.Columns)
	cellHeight := innerHeight / float32(g.Rows)

	startColIdx := startColumn - 1
	startRowIdx := startRow - 1

	return rl.Rectangle{
		X:      g.PaddingX + float32(startColIdx)*cellWidth,
		Y:      g.PaddingY + float32(startRowIdx)*cellHeight,
		Width:  float32(endColumn-startColumn+1) * cellWidth,
		Height: float32(endRow-startRow+1) * cellHeight,
	}
}

func clampInt(value, minVal, maxVal int) int {
	if value < minVal {
		return minVal
	}
	if value > maxVal {
		return maxVal
	}
	return value
}

func drawDashedLine(start, end rl.Vector2, dashLength, thickness float32, color rl.Color) {
	if dashLength <= 0 {
		dashLength = 1
	}
	dir := rl.Vector2Subtract(end, start)
	total := rl.Vector2Length(dir)
	if total <= 0 {
		return
	}
	dir = rl.Vector2Scale(dir, 1/total)

	current := start
	progressed := float32(0)
	drawSegment := true

	for progressed < total {
		segment := dashLength
		remaining := total - progressed
		if segment > remaining {
			segment = remaining
		}

		next := rl.Vector2Add(current, rl.Vector2Scale(dir, segment))
		if drawSegment {
			rl.DrawLineEx(current, next, thickness, color)
		}

		current = next
		progressed += segment
		drawSegment = !drawSegment
	}
}
