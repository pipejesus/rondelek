package experiments

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/pipejesus/rondelek/ui"
)

const EXPERIMENT_SHAPAED_WINDOW = false
const EXPERIMENT_FULLSCREEN_WINDOW = false

var dragging bool
var dragStart rl.Vector2
var pads []*ui.Pad

func init() {
	dragging = false
}

func SetupFullScreenWindow() {
	if !EXPERIMENT_FULLSCREEN_WINDOW {
		return
	}

	rl.ToggleFullscreen()
}

func SetupShapedWindow(excludeFromDrag []*ui.Pad) {
	if !EXPERIMENT_SHAPAED_WINDOW {
		return
	}

	rl.SetConfigFlags(rl.FlagWindowUndecorated | rl.FlagWindowTransparent)
	pads = excludeFromDrag
}

func HandleWindowDragging() {
	if !EXPERIMENT_SHAPAED_WINDOW {
		return
	}

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

}

func isPointerOverUI() bool {
	mouse := rl.GetMousePosition()
	for _, pad := range pads {
		if rl.CheckCollisionPointRec(mouse, pad.Rect) {
			return true
		}
	}
	return false
}
