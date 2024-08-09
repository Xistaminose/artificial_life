package simulation

import (
	"image/color"
	"math/rand"

	"github.com/aquilax/go-perlin"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Biome struct {
	name          string
	terrainColors []terrainColor
}

var (
	PlainsBiome = Biome{
		name: "Plains",
		terrainColors: []terrainColor{
			{color: color.RGBA{0, 0, 139, 255}, threshold: 0.2},      // Deep Water
			{color: color.RGBA{0, 0, 255, 255}, threshold: 0.4},      // Water
			{color: color.RGBA{238, 214, 175, 255}, threshold: 0.45}, // Sand
			{color: color.RGBA{34, 139, 34, 255}, threshold: 0.6},    // Grass
			{color: color.RGBA{0, 100, 0, 255}, threshold: 0.7},      // High Grass
			{color: color.RGBA{34, 139, 34, 255}, threshold: 0.8},    // Forest
			{color: color.RGBA{139, 69, 19, 255}, threshold: 0.9},    // Mountain
			{color: color.RGBA{255, 250, 250, 255}, threshold: 1.0},  // Snow
		},
	}

	DesertBiome = Biome{
		name: "Desert",
		terrainColors: []terrainColor{
			{color: color.RGBA{255, 235, 205, 255}, threshold: 0.3},  // Light Sand
			{color: color.RGBA{238, 214, 175, 255}, threshold: 0.6},  // Sand
			{color: color.RGBA{210, 180, 140, 255}, threshold: 0.75}, // Dunes
			{color: color.RGBA{160, 82, 45, 255}, threshold: 0.9},    // Rocky Sand
			{color: color.RGBA{139, 69, 19, 255}, threshold: 1.0},    // Rocky Outcrops
		},
	}

	TundraBiome = Biome{
		name: "Tundra",
		terrainColors: []terrainColor{
			{color: color.RGBA{0, 0, 139, 255}, threshold: 0.2},      // Deep Water
			{color: color.RGBA{0, 0, 255, 255}, threshold: 0.4},      // Water
			{color: color.RGBA{240, 255, 240, 255}, threshold: 0.5},  // Ice
			{color: color.RGBA{224, 255, 255, 255}, threshold: 0.6},  // Snowy Grass
			{color: color.RGBA{176, 196, 222, 255}, threshold: 0.75}, // Frozen Tundra
			{color: color.RGBA{255, 250, 250, 255}, threshold: 1.0},  // Snow
		},
	}

	MountainousBiome = Biome{
		name: "Mountainous",
		terrainColors: []terrainColor{
			{color: color.RGBA{0, 0, 139, 255}, threshold: 0.2},     // Deep Water
			{color: color.RGBA{0, 0, 255, 255}, threshold: 0.3},     // Water
			{color: color.RGBA{160, 82, 45, 255}, threshold: 0.5},   // Rocky Terrain
			{color: color.RGBA{139, 69, 19, 255}, threshold: 0.7},   // Mountain Base
			{color: color.RGBA{105, 105, 105, 255}, threshold: 0.8}, // Mountain
			{color: color.RGBA{169, 169, 169, 255}, threshold: 0.9}, // High Mountain
			{color: color.RGBA{255, 250, 250, 255}, threshold: 1.0}, // Snow Capped Peaks
		},
	}

	ForestBiome = Biome{
		name: "Forest",
		terrainColors: []terrainColor{
			{color: color.RGBA{0, 100, 0, 255}, threshold: 0.2},     // Deep Forest
			{color: color.RGBA{34, 139, 34, 255}, threshold: 0.5},   // Dense Forest
			{color: color.RGBA{107, 142, 35, 255}, threshold: 0.7},  // Light Forest
			{color: color.RGBA{85, 107, 47, 255}, threshold: 0.8},   // Forest Edge
			{color: color.RGBA{154, 205, 50, 255}, threshold: 0.9},  // Grassland
			{color: color.RGBA{238, 214, 175, 255}, threshold: 1.0}, // Forest Path
		},
	}

	WorldOfWarcraftBiome = Biome{
		name: "World of Warcraft",
		terrainColors: []terrainColor{
			{color: color.RGBA{72, 61, 139, 255}, threshold: 0.2},    // Deep Sea
			{color: color.RGBA{65, 105, 225, 255}, threshold: 0.4},   // Ocean
			{color: color.RGBA{255, 222, 173, 255}, threshold: 0.45}, // Coastal Sand
			{color: color.RGBA{34, 139, 34, 255}, threshold: 0.6},    // Grasslands
			{color: color.RGBA{139, 69, 19, 255}, threshold: 0.7},    // Barrens
			{color: color.RGBA{210, 105, 30, 255}, threshold: 0.8},   // Savanna
			{color: color.RGBA{112, 128, 144, 255}, threshold: 0.9},  // Storm Peaks
			{color: color.RGBA{255, 250, 250, 255}, threshold: 1.0},  // Snow Peaks
		},
	}
)

type Terrain struct {
	BaseSimulation
	grid         *Grid
	noise        *perlin.Perlin
	seed         int64
	freq         float64
	alpha        float64
	beta         float64
	octaves      int
	time         float64
	biomes       []Biome
	currentBiome int
}

type terrainColor struct {
	color     color.Color
	threshold float64
}

func GetBiomes() []Biome {
	return []Biome{
		PlainsBiome,
		DesertBiome,
		TundraBiome,
		MountainousBiome,
		ForestBiome,
		WorldOfWarcraftBiome,
	}
}

func NewTerrain(width, height int, seed int64, biomes []Biome) *Terrain {
	grid := NewGrid(width, height)
	p := perlin.NewPerlin(2, 2, 3, seed)

	t := &Terrain{
		grid:         grid,
		noise:        p,
		seed:         seed,
		freq:         0.05,
		alpha:        2,
		beta:         2,
		octaves:      3,
		time:         0,
		biomes:       biomes,
		currentBiome: 0,
	}

	t.generateTerrain()
	return t
}
func (t *Terrain) generateTerrain() {
	biome := t.biomes[t.currentBiome]
	for y := range t.grid.cells {
		for x := range t.grid.cells[y] {
			noiseValue := t.noise.Noise3D(float64(x)*t.freq, float64(y)*t.freq, t.time)
			noiseValue = (noiseValue + 1) / 2 // Normalize to 0-1

			for _, terrain := range biome.terrainColors {
				if noiseValue <= terrain.threshold {
					t.grid.cells[y][x] = terrain.color
					break
				}
			}
		}
	}
}

func (t *Terrain) Update() error {
	t.UpdatePauseState()
	if t.IsPaused() {
		return nil
	}

	t.time += 0.01 // Adjust this value to control the speed of the movement

	// Check for right-click to change the biome
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		t.changeBiome()
	}

	// Check for mouse wheel scroll to adjust terrain thresholds
	t.handleMouseWheel()

	t.generateTerrain()
	return nil
}

func (t *Terrain) handleMouseWheel() {
	x, y := ebiten.CursorPosition()
	gridX := x / cellSize
	gridY := y / cellSize

	if gridX >= 0 && gridX < t.grid.width && gridY >= 0 && gridY < t.grid.height {
		biome := &t.biomes[t.currentBiome]
		noiseValue := t.noise.Noise3D(float64(gridX)*t.freq, float64(gridY)*t.freq, t.time)
		noiseValue = (noiseValue + 1) / 2 // Normalize to 0-1

		for i := range biome.terrainColors {
			if noiseValue <= biome.terrainColors[i].threshold {
				// Get the mouse wheel offsets
				_, yoff := ebiten.Wheel()

				// Adjust threshold based on mouse wheel input
				if yoff > 0 {
					biome.terrainColors[i].threshold += 0.01
					if biome.terrainColors[i].threshold > 1.0 {
						biome.terrainColors[i].threshold = 1.0
					}
				} else if yoff < 0 {
					biome.terrainColors[i].threshold -= 0.01
					if biome.terrainColors[i].threshold < 0.0 {
						biome.terrainColors[i].threshold = 0.0
					}
				}
				break
			}
		}
	}
}

func (t *Terrain) changeBiome() {
	t.currentBiome = (t.currentBiome + 1) % len(t.biomes)
	t.seed = rand.Int63() // Randomize seed for new biome
	t.noise = perlin.NewPerlin(2, 2, 3, t.seed)
}

func (t *Terrain) Draw(screen *ebiten.Image) {
	for y := 0; y < t.grid.height; y++ {
		for x := 0; x < t.grid.width; x++ {
			ebitenutil.DrawRect(screen, float64(x*cellSize), float64(y*cellSize), cellSize, cellSize, t.grid.cells[y][x])
		}
	}

	// Display the name of the current biome
	ebitenutil.DebugPrintAt(screen, "Biome: "+t.biomes[t.currentBiome].name, 10, 10)

	if t.IsPaused() {
		ebitenutil.DebugPrintAt(screen, "Paused", screen.Bounds().Dx()/2-30, screen.Bounds().Dy()/2)
	}
}

func (t *Terrain) Layout(outsideWidth, outsideHeight int) (int, int) {
	return t.grid.width * cellSize, t.grid.height * cellSize
}
