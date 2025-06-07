package main

import "strings"

type ViewMode byte

const (
	VIEWMODE_NORMAL ViewMode = iota
	VIEWMODE_LIGHT
)

type ViewPort struct {
	x, y, w, h int
	world      *World
	viewMode   ViewMode
}

func ViewPortCreate(x, y, w, h int, world *World) ViewPort {
	viewport := ViewPort{world: world, viewMode: VIEWMODE_LIGHT}
	viewport.Edit(x, y, w, h)
	return viewport
}

func (viewport *ViewPort) Edit(x, y, w, h int) {
	viewport.x = x % viewport.world.width
	viewport.y = y % viewport.world.height
	viewport.w = w
	viewport.h = h
}

func (viewport *ViewPort) GetImage() string {

	world := viewport.world

	var buf strings.Builder
	buf.Grow(world.width*world.height + world.height)

	var datasource []byte

	if viewport.viewMode == VIEWMODE_NORMAL {
		datasource = world.data
	} else {
		datasource = make([]byte, len(world.light))
		copy(datasource, world.light)
		displaySym := []byte{32, 126, 61, 35}
		for index, val := range datasource {
			datasource[index] = displaySym[val]
		}
	}

	for y := 0; y < viewport.h; y++ {
		// for shift_x := 0; shift_x < viewport.w; shift_x++ {
		// 	buf.WriteByte(world.Get(viewport.x+shift_x, viewport.y+shift_y))
		// }
		shiftY := y * viewport.w
		sliceEnd1 := min(viewport.x+viewport.w, world.width)
		buf.Write(datasource[viewport.x+shiftY : sliceEnd1+shiftY])

		sliceEnd2 := viewport.x + viewport.w - world.width
		if sliceEnd2 >= 0 {
			buf.Write(datasource[shiftY : sliceEnd2+shiftY])
		}
		// if shift_y < viewport.h-1 {
		buf.WriteByte('\n')
		// }
	}

	return buf.String()
}
