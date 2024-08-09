package simulation

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Mode int

const (
	BlackWhite Mode = iota
	RandomColor
	SetColor
)

type GameOfLife struct {
	BaseSimulation
	grid      *Grid
	mode      Mode
	setColors []color.Color
}

func NewGameOfLife(width, height int, mode Mode, setColors []color.Color) *GameOfLife {
	grid := NewGrid(width, height)
	game := &GameOfLife{grid: grid, mode: mode, setColors: setColors}

	grid.Randomize(game.randomizeCellColor)
	return game
}

func (g *GameOfLife) randomizeCellColor(isAlive bool) color.Color {
	if !isAlive {
		return color.Black
	}

	switch g.mode {
	case BlackWhite:
		return color.White
	case RandomColor:
		return randomColor()
	case SetColor:
		return g.setColors[rand.Intn(len(g.setColors))]
	default:
		return color.White
	}
}

func (g *GameOfLife) Update() error {
	g.UpdatePauseState()
	if g.IsPaused() {
		return nil
	}

	g.grid.Update(func(liveNeighbors int, currentColor color.Color) color.Color {
		if currentColor != color.Black && (liveNeighbors == 2 || liveNeighbors == 3) {
			return currentColor
		} else if currentColor == color.Black && liveNeighbors == 3 {
			return g.randomizeCellColor(true)
		}
		return color.Black
	})
	return nil
}

func (g *GameOfLife) Draw(screen *ebiten.Image) {
	for y := 0; y < g.grid.height; y++ {
		for x := 0; x < g.grid.width; x++ {
			ebitenutil.DrawRect(screen, float64(x*cellSize), float64(y*cellSize), cellSize, cellSize, g.grid.cells[y][x])
		}
	}
}

func (g *GameOfLife) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.grid.width * cellSize, g.grid.height * cellSize
}
