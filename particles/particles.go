package particles

import (
	"math"
	"slices"
	"strings"
	"time"
)

type Particle struct {
	Lifetime int64
	Speed    float64

	X float64
	Y float64
}

type ParticleParams struct {
	X, Y int

	XScale float64

	ParticleCount int

	MaxLife  int64
	MaxSpeed float64

	ascii        Ascii
	reset        Reset
	NextPosition NextPosition
}

type (
	NextPosition func(particle *Particle, deltaMS int64)
	Ascii        func(row, col int, count [][]int) string
	Reset        func(particle *Particle, params *ParticleParams)
)

type ParticleSystem struct {
	ParticleParams
	particles []*Particle

	lastTime int64
}

func NewParticleSystem(params ParticleParams) ParticleSystem {
	particles := make([]*Particle, 0)
	for i := 0; i < params.ParticleCount; i++ {
		particles = append(particles, &Particle{})
	}
	return ParticleSystem{
		ParticleParams: params,
		particles:      particles,
		lastTime:       time.Now().UnixMilli(),
	}
}

func (ps *ParticleSystem) Start() {
	for _, p := range ps.particles {
		ps.reset(p, &ps.ParticleParams)
		// fmt.Printf("%+v\n", p)
	}
}

func (ps *ParticleSystem) Update() {
	now := time.Now().UnixMilli()
	delta := now - ps.lastTime
	ps.lastTime = now

	for _, p := range ps.particles {
		ps.NextPosition(p, delta)

		if p.Y >= float64(ps.Y) || p.X >= float64(ps.X) || p.Lifetime <= 0 {
			ps.reset(p, &ps.ParticleParams)
		}

	}
}

func (ps *ParticleSystem) Display() []string {
	counts := make([][]int, 0)

	for row := 0; row < ps.Y; row++ {
		count := make([]int, 0)
		for col := 0; col < ps.X; col++ {
			count = append(count, 0)
		}
		counts = append(counts, count)
	}

	for _, p := range ps.particles {
		row := int(math.Floor(p.Y))
		col := int(math.Round(p.X))

		counts[row][col]++
	}
	out := make([][]string, 0)
	for r, row := range counts {
		outRow := make([]string, 0)
		for c := range row {
			outRow = append(outRow, ps.ascii(r, c, counts))
		}
		out = append(out, outRow)
	}

	slices.Reverse(out)
	outStr := make([]string, 0)
	for _, row := range out {
		outStr = append(outStr, strings.Join(row, ""))
	}
	return outStr
}
