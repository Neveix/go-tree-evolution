package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

var maxAge int = 30
var logEnergy int = 4

var debug = false
var viewportScrollSpeed = 25
var skipGameTick = false
var world *World
var viewport ViewPort
var reader = bufio.NewReader(os.Stdin)

func main() {
	activeGenome.Load()

	mainMenu()
}

func mainMenu() {
	for {
		text := input("Команды: s(start)/g(playground)/q(quit)")
		switch text {
		case "s":
			gameLoop()
		case "g":
			playgroundLoop()
		case "q":
			return
		}
	}
}

func createSeveralSeeds() {
	world.CreateSeeds(700)
}

func gameLoop() {
	world = WorldCreate(2048, 20)
	viewport = ViewPortCreate(0, 0, 120, 20)
	seeds = []*Seed{}
	trees = []*Tree{}
	createSeveralSeeds()
	for {
		if !skipGameTick {
			gameLogicLoop()
		}
		skipGameTick = false
		fmt.Print(viewport.GetImage())
		if handleUserInput() != 0 {
			return
		}
	}
}

func handleUserInput() int {
	inputText := "Команды: save/load/q(quit)/seeds/s(simulate)/crs(create seeds)"
	inputText += "/l(logEnergy set)/getl(logEnergy get)/trees(tree count)"
	text := input(inputText)
	switch text {
	case "save":
		trees[0].genome.Save()
		skipGameTick = true
	case "load":
		skipGameTick = true
		if len(trees) == 0 {
			input("Нет деревьев для загрузки")
		} else {
			trees[0].genome.Load()
		}
	case "s":
		skipGameTick = true
		stepsStr := input("Сколько шагов симулировать?")
		steps, _ := strconv.Atoi(stepsStr)
		for i := 0; i < steps; i++ {
			gameLogicLoop()
		}
	case "m":
		skipGameTick = true
		viewport.viewMode = (viewport.viewMode + 1) % viewModeCount

	case "q":
		return 1
	case "seeds":
		skipGameTick = true
		PrintSeedCount()
	case "trees":
		skipGameTick = true
		PrintTreeCount()
	case "crs":
		createSeveralSeeds()
	case "a":
		viewport.Move(-viewportScrollSpeed, 0)
		skipGameTick = true
	case "d":
		viewport.Move(viewportScrollSpeed, 0)
		skipGameTick = true
	case "l":
		logEnergyStr := input("Введите новый logEnergy")
		newLogEnergy, opError := strconv.Atoi(logEnergyStr)
		if opError == nil {
			logEnergy = newLogEnergy
		} else {
			fmt.Println("Ошибка при вводе")
		}
		skipGameTick = true
	case "getl":
		skipGameTick = true
		input(fmt.Sprintf("logEnergy = %d\n", logEnergy))
	}
	return 0
}

func gameLogicLoop() {

	// семена стираются
	for _, seed := range seeds {
		world.Set(seed.x, seed.y, ' ')
		if debug {
			fmt.Println("Семечко стёрлось")
		}
	}

	// если ниже семечек земля, они превращаются в отростки
	oldSeeds := make([]*Seed, 0, len(seeds))

	for _, seed := range seeds {
		if world.Get(seed.x, seed.y+1) == '#' {
			if debug {
				fmt.Println("Семечко коснулось земли и стало деревом")
			}
			TreeCreate(seed.x, seed.y, seed.genome, seed.energy, randomSym())
		} else {
			oldSeeds = append(oldSeeds, seed)
		}
	}

	seeds = oldSeeds

	// просчитывается энергия дерева
	for _, tree := range trees {
		tree.recieveEnergy()
		tree.lostEnergy()
	}

	// деревья умирают от недостатка эн.
	// и умирают от старости
	newTrees := make([]*Tree, 0, len(trees))

	for _, tree := range trees {
		if tree.energy > 0 {
			if tree.age < maxAge {
				newTrees = append(newTrees, tree)
				tree.age = tree.age + 1
			} else {
				tree.Die()
			}
		} else {
			tree.Destroy()
		}
	}

	trees = newTrees

	// отростки растут и прев. в древ.
	for _, tree := range trees {
		if tree.age != 1 {
			tree.Grow()
		}
	}

	oldSeeds = make([]*Seed, 0, len(seeds))

	// семечка проверяет блок в ней, если он не воздух
	// если снизу воздух, семена двигаются вниз
	for _, seed := range seeds {
		if world.CanBeOccupied(seed.x, seed.y) && world.CanBeOccupied(seed.x, seed.y+1) {
			seed.y += 1
			oldSeeds = append(oldSeeds, seed)
			world.Set(seed.x, seed.y, '*')
		}

	}

	seeds = oldSeeds

	world.PerformLightUpdates()

}
