package main

import (
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type PadStatus int

const (
	PadStatusIdle PadStatus = iota
	PadStatusPressed
)

type PadAction func(p *Pad, from, to PadStatus)

type transitionRegistry map[PadStatus]map[PadStatus][]PadAction

type Pad struct {
	Rect        rl.Rectangle
	Status      PadStatus
	transitions transitionRegistry
}

func NewPad(rect rl.Rectangle) *Pad {
	return &Pad{
		Rect:        rect,
		Status:      PadStatusIdle,
		transitions: make(transitionRegistry),
	}
}

func (p *Pad) Update() {
	isPressed := gui.Button(p.Rect, "")
	oldStatus := p.Status

	if isPressed {
		p.Status = PadStatusPressed
	} else {
		p.Status = PadStatusIdle
	}

	p.ExecuteActions(oldStatus, p.Status)
}

func (p *Pad) ExecuteActions(fromStatus PadStatus, toStatus PadStatus) {
	if fromStatus == toStatus {
		return
	}

	actions := p.actionsFor(fromStatus, toStatus)
	for _, action := range actions {
		action(p, fromStatus, toStatus)
	}
}

func (p *Pad) RegisterTransition(fromStatus, toStatus PadStatus, action PadAction) {
	p.ensureRegistry()
	if _, ok := p.transitions[fromStatus]; !ok {
		p.transitions[fromStatus] = make(map[PadStatus][]PadAction)
	}
	p.transitions[fromStatus][toStatus] = append(p.transitions[fromStatus][toStatus], action)
}

func (p *Pad) actionsFor(fromStatus, toStatus PadStatus) []PadAction {
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
