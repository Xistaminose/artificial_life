package simulation

import "github.com/hajimehoshi/ebiten/v2"

const cellSize = 5

type Simulation interface {
	Update() error
	Draw(screen *ebiten.Image)
	Layout(outsideWidth, outsideHeight int) (int, int)
}
