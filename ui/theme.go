package ui

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var Theme = struct {
	// Base colors
	Background  rl.Color
	Panel       rl.Color
	PanelStroke rl.Color
	TextDark    rl.Color

	// Accent colors
	AccentOrange rl.Color
	AccentTeal   rl.Color

	// Pad colors
	PadFillPlay    rl.Color
	PadFillRecord  rl.Color
	PadLabelPlay   rl.Color
	PadLabelRecord rl.Color
	PadShadow      rl.Color

	// LED colors
	LEDEmpty rl.Color
	LEDFull  rl.Color

	// Shadow
	Shadow rl.Color

	// Layout constants
	InnerPadding float32
	RoundMd      float32
}{
	// Base colors
	Background:  rl.Color{R: 227, G: 226, B: 221, A: 255},
	Panel:       rl.Color{R: 246, G: 243, B: 235, A: 255},
	PanelStroke: rl.Color{R: 89, G: 88, B: 80, A: 100},
	TextDark:    rl.Color{R: 46, G: 49, B: 52, A: 255},

	// Accent colors
	AccentOrange: rl.Color{R: 255, G: 142, B: 57, A: 255},
	AccentTeal:   rl.Color{R: 63, G: 110, B: 121, A: 255},

	// Pad colors
	PadFillPlay:    rl.Color{R: 250, G: 248, B: 240, A: 255},
	PadFillRecord:  rl.Color{R: 255, G: 142, B: 57, A: 205},
	PadLabelPlay:   rl.Color{R: 46, G: 49, B: 52, A: 255},
	PadLabelRecord: rl.White,
	PadShadow:      rl.Color{R: 200, G: 196, B: 186, A: 255},

	// LED colors
	LEDEmpty: rl.Color{R: 90, G: 94, B: 95, A: 255},
	LEDFull:  rl.Color{R: 250, G: 50, B: 0, A: 255},

	// Shadow
	Shadow: rl.Color{R: 160, G: 156, B: 150, A: 110},

	// Layout constants
	InnerPadding: 16,
	RoundMd:      0.032,
}
