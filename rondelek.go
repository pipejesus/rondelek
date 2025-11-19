package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(800, 600, "Rondelek")
	rl.SetTargetFPS(60)
	defer rl.CloseWindow()

	pad := NewPad(rl.Rectangle{X: 10.0, Y: 10.0, Width: 100.0, Height: 100.0})
	pad.RegisterTransition(PadStatusIdle, PadStatusPressed, func(p *Pad, from, to PadStatus) {
		fmt.Println("Playing sound!")
	})

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		pad.Update()

		rl.EndDrawing()
	}
}
