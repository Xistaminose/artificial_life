package simulation

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type BaseSimulation struct {
	paused       bool
	mousePressed bool
	grid         *Grid
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

func (b *BaseSimulation) Update() error {
	return nil
}

func (b *BaseSimulation) Draw(screen *ebiten.Image) {}

func (b *BaseSimulation) Layout(outsideWidth, outsideHeight int) (int, int) {
	return b.grid.width * cellSize, b.grid.height * cellSize
}
