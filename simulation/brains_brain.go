package simulation

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	BrainOn    = 1
	BrainDying = 2
	BrainOff   = 0
)

type BriansBrain struct {
	BaseSimulation
	grid *Grid
}

func NewBriansBrain(width, height int) *BriansBrain {
	grid := NewGrid(width, height)
	bb := &BriansBrain{grid: grid}

	grid.Randomize(func(isAlive bool) color.Color {
		if isAlive {
			return color.White
		}
		return color.Black
	})
	return bb
}

func (bb *BriansBrain) Update() error {
	bb.UpdatePauseState()
	if bb.IsPaused() {
		return nil
	}

	newCells := make([][]color.Color, len(bb.grid.cells))
	for i := range bb.grid.cells {
		newCells[i] = make([]color.Color, len(bb.grid.cells[i]))
		for j := range bb.grid.cells[i] {
			currentState := bb.getState(bb.grid.cells[i][j])
			neighborCount := bb.countNeighbors(j, i)

			switch currentState {
			case BrainOn:
				newCells[i][j] = bb.stateToColor(BrainDying)
			case BrainDying:
				newCells[i][j] = bb.stateToColor(BrainOff)
			case BrainOff:
				if neighborCount == 2 {
					newCells[i][j] = bb.stateToColor(BrainOn)
				} else {
					newCells[i][j] = bb.stateToColor(BrainOff)
				}
			}
		}
	}
	bb.grid.cells = newCells
	return nil
}

func (bb *BriansBrain) Draw(screen *ebiten.Image) {
	for y := 0; y < bb.grid.height; y++ {
		for x := 0; x < bb.grid.width; x++ {
			ebitenutil.DrawRect(screen, float64(x*cellSize), float64(y*cellSize), cellSize, cellSize, bb.grid.cells[y][x])
		}
	}

	if bb.IsPaused() {
		ebitenutil.DebugPrintAt(screen, "Paused", 480/2-30, 480/2)
	}
}

func (bb *BriansBrain) Layout(outsideWidth, outsideHeight int) (int, int) {
	return bb.grid.width * cellSize, bb.grid.height * cellSize
}

func (bb *BriansBrain) getState(c color.Color) int {
	switch c {
	case color.White:
		return BrainOn
	case color.Gray{Y: 128}: // Mid-gray color for "dying"
		return BrainDying
	default:
		return BrainOff
	}
}

func (bb *BriansBrain) stateToColor(state int) color.Color {
	switch state {
	case BrainOn:
		return color.White
	case BrainDying:
		return color.Gray{Y: 128} // Mid-gray color
	default:
		return color.Black
	}
}

func (bb *BriansBrain) countNeighbors(x, y int) int {
	count := 0
	for _, d := range bb.grid.directions() {
		nx, ny := x+d.x, y+d.y
		if nx >= 0 && ny >= 0 && nx < bb.grid.width && ny < bb.grid.height {
			if bb.getState(bb.grid.cells[ny][nx]) == BrainOn {
				count++
			}
		}
	}
	return count
}
