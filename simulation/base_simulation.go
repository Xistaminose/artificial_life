package simulation

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type BaseSimulation struct {
	paused       bool
	mousePressed bool
}

func (b *BaseSimulation) UpdatePauseState() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if !b.mousePressed {
			b.paused = !b.paused
		}
		b.mousePressed = true
	} else {
		b.mousePressed = false
	}
}

func (b *BaseSimulation) IsPaused() bool {
	return b.paused
}
