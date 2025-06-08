package main

import (
	"math/rand"
)

type World struct {
	data          []byte
	light         []byte
	width         int
	height        int
	updateLightAt map[int]int
}

func WorldCreate(w, h int) *World {
	data := make([]byte, w*h)
	light := make([]byte, w*h)

	world := World{
		data:          data,
		light:         light,
		width:         w,
		height:        h,
		updateLightAt: map[int]int{},
	}

	world.InitWorldData()
	world.InitWorldLight()

	return &world
}

func (world *World) InitWorldData() {
	w, h := world.width, world.height

	groundYs := make([]int, w)

	for i := 0; i < w; i++ {
		maxShift := 10
		if i == 0 {
			groundYs[i] = rand.Intn(maxShift) + h - 1 - (maxShift - 1)
			continue
		}
		groundYs[i] = groundYs[i-1]
		if rand.Intn(7) != 0 {
			continue
		}
		inc := (rand.Intn(3) - 1)
		groundYs[i] += inc
		if groundYs[i] > h-1 {
			groundYs[i] = h - 1
		}
		if groundYs[i] < 0 {
			groundYs[i] = 0
		}

	}

	var index int
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			index = y*w + x
			if y < groundYs[x] {
				world.data[index] = ' '
			} else {
				world.data[index] = '#'
			}

		}
	}

}

func (world *World) InitWorldLight() {
	var index int
	for y := 0; y < world.height; y++ {
		for x := 0; x < world.width; x++ {
			index = y*world.width + x
			if world.CanBeOccupied(x, y) {
				world.light[index] = 3
			} else {
				world.light[index] = 0
			}
		}
	}
}

func (world *World) CoordinatesAreOutside(x, y int) bool {
	if y < 0 {
		return true
	}
	if y >= world.height {
		return true
	}
	return false
}

func (world *World) NormalizeCoords(x, y int) (int, int) {
	if x < 0 {
		x = x + world.width
	}
	x %= world.width
	return x, y
}

func (world *World) Get(x, y int) byte {
	// oldX, oldY := x, y
	x, y = world.NormalizeCoords(x, y)
	// fmt.Printf("Get(...) старые x,y = %d,%d, новые x,y=%d,%d\n", oldX, oldY, x, y)
	return world.data[y*world.width+x]
}

func (world *World) GetLight(x, y int) byte {
	return world.light[y*world.width+x]
}

func (world *World) Set(x, y int, value byte) {
	x, y = world.NormalizeCoords(x, y)
	world.data[y*world.width+x] = value
}

func (world *World) IndexToXY(index int) (x, y int) {
	x = index % world.width
	y = index / world.width
	return
}

func (world *World) CanBeOccupied(x, y int) bool {
	if y < 0 {
		return false
	}
	switch world.Get(x, y) {
	case ' ':
		return true
	case '*':
		return true
	default:
		return false
	}
}

func (world *World) CreateSeeds(n int) {
	Repeat(n, func(i int) {
		x, y := rand.Intn(world.width), rand.Intn(world.height-1)
		if world.Get(x, y) != ' ' {
			return
		}
		SeedCreate(x, y, GenomeCreate(), 50, true)
	})
}

// var updateLightAtMutex sync.Mutex

func (world *World) UpdateLightAt(x, y int) {
	// updateLightAtMutex.Lock()

	newValue := y + 1
	value := world.updateLightAt[x]
	if value == 0 || value > newValue {
		world.updateLightAt[x] = newValue
	}
	// updateLightAtMutex.Unlock()
}

func (world *World) PerformLightUpdates() {
	for x, y := range world.updateLightAt {
		world.PerformLightUpdate(x, y-1)
	}
	world.updateLightAt = map[int]int{}
}

func (world *World) PerformLightUpdate(x, y int) {
	var light byte
	if y == 0 {
		light = 3
	} else {
		light = world.light[world.width*(y-1)+x]
	}
	for {
		if world.CanBeOccupied(x, y) {
			if light < 3 {
				light += 1
			}
		} else {
			if light > 0 {
				light -= 1
			}
		}

		world.light[world.width*y+x] = light
		y++
		if y > world.height-1 {
			break
		}
	}
}
