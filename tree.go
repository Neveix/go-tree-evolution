package main

import (
	"fmt"
	"math/rand"
)

func randomSym() byte {
	return 97 + byte(rand.Intn(26))
}

type Plant struct {
	x, y int
	sym  byte
	gen  [4]byte
	old  bool
}

func PlantCreate(x, y int, sym byte, gen [4]byte) Plant {
	plant := Plant{x, y, sym, gen, false}

	world.Set(x, y, sym-32)
	world.UpdateLightAt(x, y)

	return plant
}

type Log struct {
	x, y int
	sym  byte
}

func LogCreate(x, y int, sym byte) Log {
	log := Log{x, y, sym}

	world.Set(x, y, sym)
	world.UpdateLightAt(x, y)

	return log
}

type Tree struct {
	plant     []Plant
	log       []Log
	genome    *Genome
	energy    int
	sym       byte
	age       int
	destroyed bool
}

var trees = []*Tree{}

func TreeCreate(x, y int, genome *Genome, energy int, sym byte) *Tree {
	if energy == 0 {
		energy = 260
	}
	plant := PlantCreate(x, y, sym, genome.GetGen(0))
	t := Tree{
		plant:  []Plant{plant},
		genome: genome,
		energy: energy, sym: sym,
		age: 0, destroyed: false,
	}

	trees = append(trees, &t)

	return &t
}

func (tree *Tree) Grow() {
	newPlants := []Plant{}
	// Создаём новые ростки и превращаем старые в древесину
	for _, plant := range tree.plant {
		// fmt.Printf("Обнаружен росток в %d,%d с геном %d,%d,%d,%d:\n",
		// plant.x, plant.y, plant.gen[0], plant.gen[1], plant.gen[2], plant.gen[3])

		for i := 0; i < len(sides); i += 2 {
			newGenN := plant.gen[i/2]
			// fmt.Printf("сторона %d пытается появиться %c \n", i/2, ByteTo32mal(newGenN))
			if newGenN > 15 {
				// fmt.Printf("Неподходящий newGenN %d \n", newGenN)
				continue
			}
			dx, dy := sides[i], sides[i+1]
			// fmt.Printf("сторона %d/%d; смещение: %d, %d; \n", i, len(sides), dx, dy)
			x, y := plant.x+dx, plant.y+dy
			if world.CoordinatesAreOutside(x, y) {
				continue
			}
			x, y = world.NormalizeCoords(x, y)
			if !world.CanBeOccupied(x, y) {
				// fmt.Printf("Там занято ... \n")
				continue
			}
			// fmt.Printf("Новый росток создан в %d,%d с newGenN = %d\n", x, y, newGenN)
			newPlant := PlantCreate(x, y, plant.sym, tree.genome.GetGen(newGenN))
			newPlants = append(newPlants, newPlant)
		}
		tree.log = append(tree.log, LogCreate(plant.x, plant.y, tree.sym))
	}

	tree.plant = newPlants
}

func (tree *Tree) Destroy(cleanPlants bool) {
	if cleanPlants {
		for _, plant := range tree.plant {
			world.Set(plant.x, plant.y, ' ')
			world.UpdateLightAt(plant.x, plant.y)
		}
	}

	for _, log := range tree.log {
		world.Set(log.x, log.y, ' ')
		world.UpdateLightAt(log.x, log.y)
	}

	// trees = Remove(trees, tree)
	tree.destroyed = true
}

func (tree *Tree) Die() {
	plantCount := len(tree.plant)
	if debug {
		fmt.Printf("У дерева с символом %c %d отростков\n", tree.sym, plantCount)
		if tree.destroyed {
			input("Уничтоженное дерево пытается ещё раз умереть!")
		}
	}
	if plantCount > 0 {
		energyPie := tree.energy / plantCount
		_ = energyPie
		for _, plant := range tree.plant {
			var genome *Genome
			if rand.Intn(2) == 0 {
				genome = tree.genome.Mutated()
			} else {
				genome = tree.genome
			}
			SeedCreate(plant.x, plant.y, genome, energyPie)
		}
	}
	tree.Destroy(false)

}

func (tree *Tree) recieveEnergy() {

}

func (tree *Tree) lostEnergy() {

}
