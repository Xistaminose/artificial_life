package main

import (
	"image/color"
	"log"

	"artificialLife/simulation"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 2 << 9
	screenHeight = 2 << 8
	cellSize     = 5
)

func main() {
	var sim simulation.Simulation

	// Choose the simulation mode and type
	simType := "terrain" // "game_of_life", "schelling", "brians_brain", or "terrain"
	threshold := 0.1     // Satisfaction threshold for Schelling model
	frameRate := 10      // Frame rate for the simulation

	switch simType {
	case "game_of_life":
		setColors := []color.Color{
			color.RGBA{255, 0, 0, 255}, // Red
			color.RGBA{0, 255, 0, 255}, // Green
			color.RGBA{0, 0, 255, 255}, // Blue
		}
		sim = simulation.NewGameOfLife(screenWidth/cellSize, screenHeight/cellSize, simulation.BlackWhite, setColors)
	case "schelling":
		groupColors := []color.Color{
			color.RGBA{255, 0, 0, 255}, // Red
			color.RGBA{0, 0, 255, 255}, // Blue
		}
		sim = simulation.NewSchelling(screenWidth/cellSize, screenHeight/cellSize, threshold, groupColors)
	case "brians_brain":
		sim = simulation.NewBriansBrain(screenWidth/cellSize, screenHeight/cellSize)
	case "terrain":
		sim = simulation.NewTerrain(screenWidth/cellSize, screenHeight/cellSize, 0, simulation.GetBiomes())
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetMaxTPS(frameRate)
	ebiten.SetWindowTitle("Simulation")
	if err := ebiten.RunGame(sim); err != nil {
		log.Fatal(err)
	}
}
