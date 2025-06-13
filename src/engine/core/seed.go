package core

import "fmt"

type Seed struct {
	X, Y   int
	Genome *Genome
	Energy int
}

const seedEnergyLight = 1000

var Seeds []*Seed

func SeedCreate(x, y int, genome *Genome, energy int, doDraw bool) *Seed {

	if energy > seedEnergyLight {
		energy = seedEnergyLight
	}

	seed := Seed{x, y, genome, energy}
	if doDraw {
		MainWorld.Set(x, y, '*')
	}

	Seeds = append(Seeds, &seed)

	return &seed
}

func (seed *Seed) Destroy() {
	Seeds = Remove(Seeds, seed)
}

func PrintSeedCount() {
	fmt.Printf("Семечек: %d\n", len(Seeds))
}
