package simulation

import (
	"image/color"
	"math/rand"
)

type Grid struct {
	cells  [][]color.Color
	width  int
	height int
}

func NewGrid(width, height int) *Grid {
	g := &Grid{
		cells:  make([][]color.Color, height),
		width:  width,
		height: height,
	}
	for i := range g.cells {
		g.cells[i] = make([]color.Color, width)
	}
	return g
}

func randomColor() color.Color {
	return color.RGBA{
		R: uint8(rand.Intn(256)),
		G: uint8(rand.Intn(256)),
		B: uint8(rand.Intn(256)),
		A: 255,
	}
}

func (g *Grid) Randomize(randomizeColorFunc func(isAlive bool) color.Color) {
	for y := range g.cells {
		for x := range g.cells[y] {
			g.cells[y][x] = randomizeColorFunc(rand.Float64() < 0.5)
		}
	}
}

func (g *Grid) directions() []struct{ x, y int } {
	return []struct{ x, y int }{
		{-1, -1}, {0, -1}, {1, -1},
		{-1, 0}, {1, 0},
		{-1, 1}, {0, 1}, {1, 1},
	}
}

func (g *Grid) CountLiveNeighbors(x, y int) int {
	directions := []struct{ x, y int }{
		{-1, -1}, {0, -1}, {1, -1},
		{-1, 0}, {1, 0},
		{-1, 1}, {0, 1}, {1, 1},
	}

	count := 0
	for _, d := range directions {
		nx, ny := x+d.x, y+d.y
		if nx >= 0 && ny >= 0 && nx < g.width && ny < g.height && g.cells[ny][nx] != color.Black {
			count++
		}
	}
	return count
}

func (g *Grid) Update(rules func(int, color.Color) color.Color) {
	newCells := make([][]color.Color, len(g.cells))
	for i := range g.cells {
		newCells[i] = make([]color.Color, len(g.cells[i]))
		for j := range g.cells[i] {
			liveNeighbors := g.CountLiveNeighbors(j, i)
			newCells[i][j] = rules(liveNeighbors, g.cells[i][j])
		}
	}
	g.cells = newCells
}
