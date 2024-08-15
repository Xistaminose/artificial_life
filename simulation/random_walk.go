package simulation

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type randomWalker struct {
	BaseSimulation
	grid *Grid
}

func NewrandomWalker(width, height int) *randomWalker {
	fmt.Println("Creating random walker")
	grid := NewGrid(width, height)

	rw := &randomWalker{grid: grid}

	return rw
}

func (rw *randomWalker) update() error {
	return nil
}

func (rw *randomWalker) Draw(screen *ebiten.Image) {
	rw.grid.Draw(screen, cellSize)

	if rw.IsPaused() {
		ebitenutil.DebugPrintAt(screen, "Paused", screen.Bounds().Dx()/2-30, screen.Bounds().Dy()/2)
	}
}

func (rw *randomWalker) Layout(int, int) (int, int) {
	return rw.grid.width * cellSize, rw.grid.height * cellSize
}
