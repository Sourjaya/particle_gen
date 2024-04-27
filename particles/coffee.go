package particles

import (
	"log"
	"math"
	"math/rand"
	"time"
)

type Coffee struct {
	ParticleSystem
}

func ascii(row, col int, counts [][]int) string {
	count := counts[row][col]
	if count < 2 {
		return " "
	}
	if count < 6 {
		return "."
	}
	if count < 9 {
		return ":"
	}
	if count < 12 {
		return "{"
	}
	return "}"
}

func reset(p *Particle, params *ParticleParams) {
	p.Lifetime = int64(math.Floor(float64(params.MaxLife) * rand.Float64()))
	p.Speed = params.MaxSpeed * rand.Float64()

	maxX := math.Floor(float64(params.X) / 2)
	// fmt.Print(maxX)
	x := math.Max(-maxX, math.Min(rand.NormFloat64()*params.XScale, maxX))
	p.X = x + maxX
	p.Y = 0
}

func nextPos(particle *Particle, deltaMS int64) {
	particle.Lifetime -= deltaMS
	if particle.Lifetime <= 0 {
		return
	}
	particle.Y += particle.Speed * (float64(deltaMS) / 1000.0)
}

func mutate(row, col int, counts [][]int) {
	if row == 0 || col == 0 || col == len(counts[0])-1 || row == len(counts)-1 {
		return
	}
	dirs := [][]int{
		{-1, -1},
		{1, 0},
		{-1, 1},

		{0, -1},
		{0, 1},

		{1, 0},
		{1, 1},
		{1, -1},
	}
	count := 0
	for _, dir := range dirs {
		count = counts[row+dir[0]][col+dir[1]]
	}
	if count > 3 {
		counts[row][col] = 0
	}
}

func NewCoffee(width, height int, scale float64) Coffee {
	if width%2 == 0 {
		log.Fatal("width of the particle system must be odd")
	}
	startTime := time.Now().UnixMilli()
	ascii := func(row, col int, counts [][]int) string {
		mutate(row, col, counts)
		count := counts[row][col]
		if count == 0 {
			return " "
		}
		if count < 2 {
			return "▒"
		}
		if count < 4 {
			return "▓"
		}
		/*if count == 1 {
			return "."
		}*/

		direction := row + int(((time.Now().UnixMilli()-startTime)/2000)%2)

		if direction%2 == 0 {
			return "█"
		}

		return "▓"
	}
	params := ParticleParams{
		MaxLife:       7000,
		MaxSpeed:      1.5,
		ParticleCount: 700,

		reset:        reset,
		ascii:        ascii,
		NextPosition: nextPos,

		XScale: scale,
		X:      width,
		Y:      height,
	}
	return Coffee{
		ParticleSystem: NewParticleSystem(params),
	}
}
