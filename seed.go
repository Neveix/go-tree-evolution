package main

import "fmt"

type Seed struct {
	x, y   int
	genome *Genome
	energy int
}

const seedEnergyLight = 1000

var seeds []*Seed

func SeedCreate(x, y int, genome *Genome, energy int, doDraw bool) *Seed {

	if energy > seedEnergyLight {
		energy = seedEnergyLight
	}

	seed := Seed{x, y, genome, energy}
	if doDraw {
		world.Set(x, y, '*')
	}

	seeds = append(seeds, &seed)

	return &seed
}

func (seed *Seed) Destroy() {
	seeds = Remove(seeds, seed)
}

func PrintSeedCount() {
	fmt.Printf("Семечек: %d\n", len(seeds))
}
