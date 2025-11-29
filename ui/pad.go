package ui

import (
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type PressStatus int
type Mode int

const (
	PadStatusIdle PressStatus = iota
	PadStatusPressed
)

const (
	ModePlay Mode = iota
	ModeRecord
)

type PadAction func(p *Pad, from, to PressStatus)

type transitionRegistry map[PressStatus]map[PressStatus][]PadAction

type Pad struct {
	Rect        rl.Rectangle
	Status      PressStatus
	Mode        Mode
	transitions transitionRegistry
	SampleIdx   int
	Key         int32
}

func NewPad(rect rl.Rectangle, key int32) *Pad {
	return &Pad{
		Rect:        rect,
		Status:      PadStatusIdle,
		Mode:        ModePlay,
		transitions: make(transitionRegistry),
		SampleIdx:   -1,
		Key:         key,
	}
}

func (p *Pad) Draw() any {
	if p.Mode == ModePlay {
		return gui.Button(p.Rect, "")
	}

	return gui.Button(p.Rect, "REC")
}

func (p *Pad) SetSample(samplerIdx int) {
	p.SampleIdx = samplerIdx
}

func (p *Pad) HasSample() bool {
	return p.SampleIdx >= 0
}

func (p *Pad) ToggleMode() {
	if p.Mode == ModePlay {
		p.Mode = ModeRecord
	} else {
		p.Mode = ModePlay
	}
}

func (p *Pad) Update() {
	p.Draw()

	wasPressed := p.Status == PadStatusPressed
	isInside := rl.CheckCollisionPointRec(rl.GetMousePosition(), p.Rect)

	isHeldByMouse := isInside && rl.IsMouseButtonDown(rl.MouseButtonLeft)

	isHeldByKeyboard := rl.IsKeyDown(p.Key)
	isHeld := isHeldByMouse || isHeldByKeyboard

	if isHeld {
		p.Status = PadStatusPressed
	} else {
		p.Status = PadStatusIdle
	}

	if wasPressed != isHeld {
		from := PadStatusIdle
		to := PadStatusPressed
		if !isHeld {
			from, to = PadStatusPressed, PadStatusIdle
		}
		p.ExecuteActions(from, to)
	}
}

func (p *Pad) ExecuteActions(fromStatus PressStatus, toStatus PressStatus) {
	if fromStatus == toStatus {
		return
	}

	actions := p.actionsFor(fromStatus, toStatus)
	for _, action := range actions {
		action(p, fromStatus, toStatus)
	}
}

func (p *Pad) RegisterTransition(fromStatus, toStatus PressStatus, action PadAction) {
	p.ensureRegistry()
	if _, ok := p.transitions[fromStatus]; !ok {
		p.transitions[fromStatus] = make(map[PressStatus][]PadAction)
	}
	p.transitions[fromStatus][toStatus] = append(p.transitions[fromStatus][toStatus], action)
}

func (p *Pad) actionsFor(fromStatus, toStatus PressStatus) []PadAction {
	if p.transitions == nil {
		return nil
	}
	if toMap, ok := p.transitions[fromStatus]; ok {
		return toMap[toStatus]
	}
	return nil
}

func (p *Pad) ensureRegistry() {
	if p.transitions == nil {
		p.transitions = make(transitionRegistry)
	}
}
