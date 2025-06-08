package main

import "fmt"

var printGenome = true
var activeGenome *Genome = GenomeCreate()

var playgroundWorldWidth, playgroundWorldHeight = 40, 20

func playgroundReloadTree() {
	w, h := playgroundWorldWidth, playgroundWorldHeight
	world = WorldCreate(w, h)
	viewport = ViewPortCreate(0, 0, w, h)
	seeds = []*Seed{}
	trees = []*Tree{}
	genome := activeGenome
	TreeCreate(w/2, h-2, genome, 260, 'a')
}

func playgroundLoop() {
	playgroundReloadTree()
	for {
		drawWorldAndStuff()
		playgroundLogicLoop()

		text := input("Команды: save/load/q(quit)/seeds/g(genome)/x(mutate)")
		switch text {
		case "save":
			activeGenome.Save()
			input("Сохранено. Нажмите для продолжения")
		case "load":
			activeGenome.Load()
			playgroundReloadTree()
			input("Загружено. Нажмите для продолжения")
		case "g":
			printGenome = !printGenome
		case "a":
			trees[0].genome.Mutate()
		case "q":
			return
		case "seeds":
			PrintSeedCount()
		}
	}
}

func drawWorldAndStuff() {
	outputText := viewport.GetImage()
	if printGenome {
		outputText2 := trees[0].genome.ToString(true)
		outputText = MergeLines(outputText, outputText2, " ")
	}
	fmt.Print(outputText)
}

func playgroundLogicLoop() {
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
			// TreeCreate(seed.x, seed.y, seed.genome, seed.energy, randomSym())
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

	fmt.Printf("%d\n", len(trees))

	for _, tree := range trees {
		if tree.energy > 0 {
			if tree.age < maxAge {
				newTrees = append(newTrees, tree)
				tree.age = tree.age + 1
			} else {
				// tree.Die()
				tree.Destroy()
			}
		} else {
			tree.Destroy()
		}
	}

	trees := newTrees

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

	if len(trees) == 0 || len(trees[0].plant) == 0 {
		playgroundReloadTree()
	}

}
