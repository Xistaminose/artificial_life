package simulation

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Schelling struct {
	BaseSimulation
	grid           *Grid
	threshold      float64
	emptyColor     color.Color
	groupColors    []color.Color
	stableCounter  int     // Counts how many updates the grid has been stable
	maxStableCount int     // Maximum number of stable updates before increasing the threshold
	stableLimit    float64 // The percentage of satisfied agents to consider the grid stable
}

func NewSchelling(width, height int, threshold float64, groupColors []color.Color) *Schelling {
	grid := NewGrid(width, height)
	sim := &Schelling{
		grid:           grid,
		threshold:      threshold,
		emptyColor:     color.Black,
		groupColors:    groupColors,
		maxStableCount: 5, // Number of stable updates before increasing the threshold
		stableLimit:    1, // Consider grid stable if 95% or more agents are satisfied
	}

	sim.randomizeGrid(0.05)
	return sim
}

func (s *Schelling) randomizeGrid(emptyRatio float64) {
	for y := range s.grid.cells {
		for x := range s.grid.cells[y] {
			if rand.Float64() < emptyRatio {
				s.grid.cells[y][x] = s.emptyColor
			} else {
				s.grid.cells[y][x] = s.groupColors[rand.Intn(len(s.groupColors))]
			}
		}
	}
}

func (s *Schelling) isSatisfied(x, y int) bool {
	currentColor := s.grid.cells[y][x]
	if currentColor == s.emptyColor {
		return true
	}

	likeNeighbors := 0
	totalNeighbors := 0

	for _, d := range s.grid.directions() {
		nx, ny := x+d.x, y+d.y
		if nx >= 0 && ny >= 0 && nx < s.grid.width && ny < s.grid.height {
			neighborColor := s.grid.cells[ny][nx]
			if neighborColor != s.emptyColor {
				totalNeighbors++
				if neighborColor == currentColor {
					likeNeighbors++
				}
			}
		}
	}

	if totalNeighbors == 0 {
		return true
	}

	return float64(likeNeighbors)/float64(totalNeighbors) >= s.threshold
}

func (s *Schelling) moveAgent(x, y int) {
	for {
		nx, ny := rand.Intn(s.grid.width), rand.Intn(s.grid.height)
		if s.grid.cells[ny][nx] == s.emptyColor {
			s.grid.cells[ny][nx] = s.grid.cells[y][x]
			s.grid.cells[y][x] = s.emptyColor
			break
		}
	}
}

func (s *Schelling) Update() error {
	s.UpdatePauseState()
	if s.IsPaused() {
		return nil
	}

	totalAgents := 0
	satisfiedAgents := 0

	for y := range s.grid.cells {
		for x := range s.grid.cells[y] {
			if s.grid.cells[y][x] != s.emptyColor {
				totalAgents++
				if s.isSatisfied(x, y) {
					satisfiedAgents++
				} else {
					s.moveAgent(x, y)
				}
			}
		}
	}

	// Calculate the satisfaction ratio
	satisfactionRatio := float64(satisfiedAgents) / float64(totalAgents)

	// Check if the grid is stable
	if satisfactionRatio >= s.stableLimit {
		s.stableCounter++
	} else {
		s.stableCounter = 0
	}

	// If the grid has been stable for the maxStableCount, increase the threshold
	if s.stableCounter >= s.maxStableCount {
		s.threshold += 0.1
		s.stableCounter = 0 // Reset the stable counter
	}

	return nil
}

func (s *Schelling) Draw(screen *ebiten.Image) {
	for y := 0; y < s.grid.height; y++ {
		for x := 0; x < s.grid.width; x++ {
			ebitenutil.DrawRect(screen, float64(x*cellSize), float64(y*cellSize), cellSize, cellSize, s.grid.cells[y][x])
		}
	}

	// Draw the current threshold value on the screen
	ebitenutil.DebugPrint(screen, "Threshold: "+fmt.Sprintf("%.2f", s.threshold))
}

func (s *Schelling) Layout(outsideWidth, outsideHeight int) (int, int) {
	return s.grid.width * cellSize, s.grid.height * cellSize
}
