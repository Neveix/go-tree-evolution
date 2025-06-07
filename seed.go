package main

import "fmt"

type Seed struct {
	x, y   int
	genome *Genome
	energy int
}

var seeds []*Seed

func SeedCreate(x, y int, genome *Genome, energy int) *Seed {
	seed := Seed{x, y, genome, energy}
	world.Set(x, y, '*')

	seeds = append(seeds, &seed)

	return &seed
}

func (seed *Seed) Destroy() {
	seeds = Remove(seeds, seed)
}

func PrintSeedCount() {
	fmt.Printf("Семечек: %d\n", len(seeds))
}
