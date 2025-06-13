package core

import (
	"math/rand"
)

type World struct {
	data          []byte
	light         []byte
	Width         int
	Height        int
	updateLightAt map[int]int
}

func WorldCreate(w, h int) *World {
	data := make([]byte, w*h)
	light := make([]byte, w*h)

	world := World{
		data:          data,
		light:         light,
		Width:         w,
		Height:        h,
		updateLightAt: map[int]int{},
	}

	world.InitWorldData()
	world.InitWorldLight()

	return &world
}

func (world *World) InitWorldData() {
	w, h := world.Width, world.Height

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
	for y := 0; y < world.Height; y++ {
		for x := 0; x < world.Width; x++ {
			index = y*world.Width + x
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
	if y >= world.Height {
		return true
	}
	return false
}

func (world *World) NormalizeCoords(x, y int) (int, int) {
	if x < 0 {
		x = x + world.Width
	}
	x %= world.Width
	return x, y
}

func (world *World) Get(x, y int) byte {
	// oldX, oldY := x, y
	x, y = world.NormalizeCoords(x, y)
	// fmt.Printf("Get(...) старые x,y = %d,%d, новые x,y=%d,%d\n", oldX, oldY, x, y)
	return world.data[y*world.Width+x]
}

func (world *World) GetLight(x, y int) byte {
	return world.light[y*world.Width+x]
}

func (world *World) Set(x, y int, value byte) {
	x, y = world.NormalizeCoords(x, y)
	world.data[y*world.Width+x] = value
}

func (world *World) IndexToXY(index int) (x, y int) {
	x = index % world.Width
	y = index / world.Width
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
		x, y := rand.Intn(world.Width), rand.Intn(world.Height-1)
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
		light = world.light[world.Width*(y-1)+x]
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

		world.light[world.Width*y+x] = light
		y++
		if y > world.Height-1 {
			break
		}
	}
}
