package main

import (
	"fmt"
	"strconv"

	"evo1/core"
)

var skipGameTick = false
var viewport *core.ViewPort
var maxAge int = 16
var logEnergy int = 8
var viewportScrollSpeed = 25

func main() {
	mainMenu()
}

func mainMenu() {
	for {
		text := core.Input("Команды: s(start)/a(run api)/q(quit)")
		switch text {
		case "s":
			gameLoop()
		case "a":
			RunAPI()
		case "q":
			return
		}
	}
}

func createSeveralSeeds() {
	core.MainWorld.CreateSeeds(700)
}

func initGame() {
	core.MainWorld = core.WorldCreate(2048, 20*4)
	viewport = core.ViewPortCreate(0, 0, 120, core.MainWorld.Height)
	core.Seeds = []*core.Seed{}
	core.Trees = []*core.Tree{}
	createSeveralSeeds()
}

func gameLoop() {
	initGame()
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
	inputText := "Команды: save/load/q(quit)/se(seed count)/s(simulate)/crs(create core.Seeds)"
	inputText += "/l(logEnergy)/tr(tree count)"
	text := core.Input(inputText)
	switch text {
	case "save":
		core.Trees[0].Genome.Save()
		skipGameTick = true
	case "load":
		skipGameTick = true
		if len(core.Trees) == 0 {
			core.Input("Нет деревьев для загрузки")
		} else {
			core.Trees[0].Genome.Load()
		}
	case "s":
		skipGameTick = true
		stepsStr := core.Input("Сколько шагов симулировать?")
		steps, _ := strconv.Atoi(stepsStr)
		for i := 0; i < steps; i++ {
			gameLogicLoop()
		}
	case "m":
		skipGameTick = true
		viewport.ViewMode = (viewport.ViewMode + 1) % core.VIEWMODE_COUNT
	case "q":
		return 1
	case "se":
		skipGameTick = true
		core.PrintSeedCount()
	case "tr":
		skipGameTick = true
		core.PrintTreeCount()
	case "crs":
		createSeveralSeeds()
	case "a":
		viewport.Move(-viewportScrollSpeed, 0)
		skipGameTick = true
	case "d":
		viewport.Move(viewportScrollSpeed, 0)
		skipGameTick = true
	case "l":
		logEnergyStr := core.Input(fmt.Sprintf("logEnergy = %d, введите новый или q",
			logEnergy))
		newLogEnergy, opError := strconv.Atoi(logEnergyStr)
		if opError == nil {
			logEnergy = newLogEnergy
		} else {
			fmt.Println("Ошибка при вводе")
		}
		skipGameTick = true
	}
	return 0
}

func gameLogicLoop() {

	// семена стираются
	for _, seed := range core.Seeds {
		core.MainWorld.Set(seed.X, seed.Y, ' ')
		if core.Debug {
			fmt.Println("Семечко стёрлось")
		}
	}

	// если ниже семечек земля, они превращаются в отростки
	oldSeeds := make([]*core.Seed, 0, len(core.Seeds))

	for _, seed := range core.Seeds {
		if core.MainWorld.Get(seed.X, seed.Y+1) == '#' {
			if core.Debug {
				fmt.Println("Семечко коснулось земли и стало деревом")
			}
			core.TreeCreate(seed.X, seed.Y, seed.Genome, seed.Energy, core.RandomSym())
		} else {
			oldSeeds = append(oldSeeds, seed)
		}
	}

	core.Seeds = oldSeeds

	// просчитывается энергия дерева
	for _, tree := range core.Trees {
		tree.RecieveEnergy()
		tree.LostEnergy(logEnergy)
	}

	// деревья умирают от недостатка эн.
	// и умирают от старости
	newTrees := make([]*core.Tree, 0, len(core.Trees))

	for _, tree := range core.Trees {
		if tree.Energy > 0 {
			if tree.Age < maxAge {
				newTrees = append(newTrees, tree)
				tree.Age = tree.Age + 1
			} else {
				tree.Die()
			}
		} else {
			tree.Destroy()
		}
	}

	core.Trees = newTrees

	// отростки растут и прев. в древ.
	for _, tree := range core.Trees {
		if tree.Age != 1 {
			tree.Grow()
		}
	}

	oldSeeds = make([]*core.Seed, 0, len(core.Seeds))

	// семечка проверяет блок в ней, если он не воздух
	// если снизу воздух, семена двигаются вниз
	for _, seed := range core.Seeds {
		if core.MainWorld.CanBeOccupied(seed.X, seed.Y) && core.MainWorld.CanBeOccupied(seed.X, seed.Y+1) {
			seed.Y += 1
			oldSeeds = append(oldSeeds, seed)
			core.MainWorld.Set(seed.X, seed.Y, '*')
		}

	}

	core.Seeds = oldSeeds

	core.MainWorld.PerformLightUpdates()
}
