package main

import (
	"strings"
)

type ViewMode byte

const (
	VIEWMODE_NORMAL ViewMode = iota
	VIEWMODE_LIGHT
)

type ViewPort struct {
	x, y, w, h int
	viewMode   ViewMode
}

func ViewPortCreate(x, y, w, h int) ViewPort {
	viewport := ViewPort{viewMode: VIEWMODE_NORMAL}
	viewport.Edit(x, y, w, h)
	return viewport
}

func (viewport *ViewPort) Edit(x, y, w, h int) {
	viewport.w = w
	viewport.h = h
	viewport.MoveTo(x, y)
}

func (viewport *ViewPort) Move(shiftX, shiftY int) {
	viewport.MoveTo(viewport.x+shiftX, viewport.y+shiftY)
}

func (viewport *ViewPort) MoveTo(x, y int) {
	viewport.x = (x + world.width) % world.width
	viewport.y = (y + world.height) % world.height
}

func (viewport *ViewPort) ViewPortXToWorldX(x int) int {
	worldViewportXRatio := float64(world.width) / float64(viewport.w)
	return int(worldViewportXRatio * float64(x))
}

func (viewport *ViewPort) GetImage() string {
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

	sliceEnd1 := min(viewport.x+viewport.w, world.width)
	sliceEnd2 := viewport.x + viewport.w - world.width

	for x := 0; x < viewport.w; x++ {
		wX := viewport.ViewPortXToWorldX(x)
		if viewport.x <= wX && wX <= sliceEnd1 {
			buf.WriteByte(64)
		} else if sliceEnd2 >= 0 && wX <= sliceEnd2 {
			buf.WriteByte(64)
		} else {
			buf.WriteByte(45)
		}
	}
	buf.WriteByte('\n')

	for y := 0; y < viewport.h; y++ {
		// for shift_x := 0; shift_x < viewport.w; shift_x++ {
		// 	buf.WriteByte(world.Get(viewport.x+shift_x, viewport.y+shift_y))
		// }
		shiftY := y * world.width

		buf.Write(datasource[viewport.x+shiftY : sliceEnd1+shiftY])

		if sliceEnd2 >= 0 {
			buf.Write(datasource[shiftY : sliceEnd2+shiftY])
		}
		// if shift_y < viewport.h-1 {
		buf.WriteByte('\n')
		// }
	}

	return buf.String()
}
