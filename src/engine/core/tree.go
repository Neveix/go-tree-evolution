package core

import (
	"fmt"
	"math/rand"
)

func RandomSym() byte {
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

	MainWorld.Set(x, y, sym-32)
	MainWorld.UpdateLightAt(x, y)

	return plant
}

type Log struct {
	x, y int
	sym  byte
}

func LogCreate(x, y int, sym byte) Log {
	log := Log{x, y, sym}

	MainWorld.Set(x, y, sym)
	MainWorld.UpdateLightAt(x, y)

	return log
}

type Tree struct {
	plant     []Plant
	log       []Log
	Genome    *Genome
	Energy    int
	sym       byte
	Age       int
	destroyed bool
}

var Trees = []*Tree{}

// var TreesMutex sync.Mutex

func TreeCreate(x, y int, Genome *Genome, energy int, sym byte) *Tree {
	if energy == 0 {
		energy = 260
	}
	plant := PlantCreate(x, y, sym, Genome.GetGen(0))
	t := Tree{
		plant:  []Plant{plant},
		Genome: Genome,
		Energy: energy, sym: sym,
		Age: 0, destroyed: false,
	}

	// TreesMutex.Lock()
	Trees = append(Trees, &t)
	// TreesMutex.Unlock()

	return &t
}

func (tree *Tree) Grow() {
	newPlants := make([]Plant, 0, len(tree.plant))
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
			if MainWorld.CoordinatesAreOutside(x, y) {
				continue
			}
			x, y = MainWorld.NormalizeCoords(x, y)
			if !MainWorld.CanBeOccupied(x, y) {
				// fmt.Printf("Там занято ... \n")
				continue
			}
			// fmt.Printf("Новый росток создан в %d,%d с newGenN = %d\n", x, y, newGenN)
			newPlant := PlantCreate(x, y, plant.sym, tree.Genome.GetGen(newGenN))
			newPlants = append(newPlants, newPlant)
		}
		tree.log = append(tree.log, LogCreate(plant.x, plant.y, tree.sym))
	}

	tree.plant = newPlants
}

func (tree *Tree) Destroy() {
	for _, plant := range tree.plant {

		MainWorld.Set(plant.x, plant.y, ' ')
		MainWorld.UpdateLightAt(plant.x, plant.y)
	}

	for _, log := range tree.log {
		MainWorld.Set(log.x, log.y, ' ')
		MainWorld.UpdateLightAt(log.x, log.y)
	}

	// Trees = Remove(Trees, tree)
	tree.destroyed = true
}

func (tree *Tree) Die() {
	plantCount := len(tree.plant)
	if Debug {
		fmt.Printf("У дерева с символом %c %d отростков\n", tree.sym, plantCount)
		if tree.destroyed {
			Input("Уничтоженное дерево пытается ещё раз умереть!")
		}
	}
	tree.Destroy()
	if plantCount > 0 {
		energyPie := tree.Energy / plantCount
		_ = energyPie
		for _, plant := range tree.plant {
			var Genome *Genome
			if rand.Intn(2) == 0 {
				Genome = tree.Genome.Mutated()
			} else {
				Genome = tree.Genome
			}
			SeedCreate(plant.x, plant.y, Genome, energyPie, false)
		}
	}
}

func (tree *Tree) RecieveEnergy() {
	for _, log := range tree.log {
		light1 := MainWorld.GetLight(log.x, log.y)

		tree.Energy += int(light1) * int(10+float64(MainWorld.height-log.y)*0.2)
	}
}

func (tree *Tree) LostEnergy(logEnergy int) {
	logCount := len(tree.log)
	tree.Energy -= logCount * logEnergy
}

func PrintTreeCount() {
	fmt.Printf("Tree Count: %d\n", len(Trees))
}
