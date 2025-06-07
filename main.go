package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

var maxAge int = 16

var debug = false
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

func gameLoop() {
	world = WorldCreate(120, 20)
	viewport = ViewPortCreate(0, 0, 120, 20, world)
	seeds = []*Seed{}
	trees = []*Tree{}
	world.CreateSeeds(80)
	for {
		gameLogicLoop()
		fmt.Print(viewport.GetImage())
		text := input("Команды: save/load/q(quit)/seeds/s(simulate)")
		switch text {
		case "save":
			trees[0].genome.Save()
			input("Сохранено. Нажмите для продолжения")
		case "load":
			if len(trees) == 0 {
				input("Нет деревьев для загрузки")
			} else {
				trees[0].genome.Load()
				input("Загружено. Нажмите для продолжения")
			}
		case "s":
			stepsStr := input("Сколько шагов симулировать?")
			steps, _ := strconv.Atoi(stepsStr)
			for i := 0; i < steps-1; i++ {
				gameLogicLoop()
			}
		case "m":
			switch viewport.viewMode {
			case VIEWMODE_NORMAL:
				viewport.viewMode = VIEWMODE_LIGHT
			case VIEWMODE_LIGHT:
				viewport.viewMode = VIEWMODE_NORMAL
			}
		case "q":
			return
		case "seeds":
			PrintSeedCount()
		}
	}
}

func gameLogicLoop() {
	// print ("before input called")

	// world.resetlight2()

	// семена стираются
	for _, seed := range seeds {
		world.Set(seed.x, seed.y, ' ')
		if debug {
			fmt.Println("Семечко стёрлось")
		}
	}

	// если ниже семечек земля, они превращаются в отростки
	oldSeeds := []*Seed{}

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
	seeds = []*Seed{}

	// просчитывается энергия дерева
	for _, tree := range trees {
		tree.recieveEnergy()
		tree.lostEnergy()
	}

	// деревья умирают от недостатка эн.
	// и умирают от старости
	newTrees := []*Tree{}

	for _, tree := range trees {
		if tree.energy > 0 {
			if tree.age < maxAge {
				newTrees = append(newTrees, tree)
				tree.age = tree.age + 1
			} else {
				tree.Die()
			}
		} else {
			tree.Destroy(true)
		}
	}

	trees = newTrees

	// отростки растут и прев. в древ.
	for _, tree := range trees {
		tree.Grow()
	}

	// семечка проверяет блок в ней, если он не воздух
	// если снизу воздух, семена двигаются вниз
	for _, seed := range oldSeeds {
		if debug {
			fmt.Println("Старое семечко проверяет что снизу и в нём...")
		}
		if world.CanBeOccupied(seed.x, seed.y) && world.CanBeOccupied(seed.x, seed.y+1) {
			if debug {
				fmt.Println("Всё норм")
			}
			seed.y += 1
			seeds = append(seeds, seed)
		}

	}

	// семена рендерятся
	for _, seed := range seeds {
		if debug {
			fmt.Println("Семечко рендерится")
		}
		world.Set(seed.x, seed.y, '*')
	}

}
