// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package gridworldmap worldmap.go Defines the map that actors live and move on.
package gridworldmap

import (
	"sort"

	"github.com/bluemun/munfall"
	"github.com/bluemun/munfall/traits"
)

type worldMap2DGrid struct {
	world           munfall.World
	cWidth, cHeight float32
	width, height   uint
	grid            []*cell2DRectGrid
}

// CreateGridWorldMap creates a 2D implementation of a munfall.
func CreateGridWorldMap(width, height uint, cellWidth, cellHeight float32) munfall.WorldMap {
	grid := make([]*cell2DRectGrid, width*height)

	var init uint
	for y := init; y < height; y++ {
		for x := init; x < width; x++ {
			grid[x+y*width] = &cell2DRectGrid{
				space: make([]munfall.Space, 0),
				pos:   &munfall.MPos{X: x, Y: y},
			}
		}
	}

	for y := init; y < height; y++ {
		for x := init; x < width; x++ {
			if x != 0 {
				grid[x+y*width].lc = grid[(x-1)+y*width]
			}
			if x != width-1 {
				grid[x+y*width].rc = grid[(x+1)+y*width]
			}
			if y != 0 {
				grid[x+y*width].tc = grid[x+(y-1)*width]
			}
			if y != height-1 {
				grid[x+y*width].bc = grid[x+(y+1)*width]
			}
		}
	}

	return &worldMap2DGrid{
		width:   width,
		height:  height,
		cWidth:  cellWidth,
		cHeight: cellHeight,
		grid:    grid,
	}
}

func (wm *worldMap2DGrid) Width() float32 {
	return float32(wm.width) * wm.cWidth
}

func (wm *worldMap2DGrid) Height() float32 {
	return float32(wm.height) * wm.cHeight
}

func (wm *worldMap2DGrid) Initialize(world munfall.World) {
	wm.world = world
}

func (wm *worldMap2DGrid) CellAt(pos *munfall.MPos) munfall.Cell {
	return wm.grid[pos.X+pos.Y*wm.width]
}

func (wm *worldMap2DGrid) CreatePath(positions []*munfall.WPos) munfall.Path {
	path := &path2DGrid{
		m:    wm,
		cell: wm.CellAt(wm.ConvertToMPos(positions[0])).(*cell2DRectGrid),
	}
	path.first = path
	for _, pos := range positions[1:] {
		cell := wm.CellAt(wm.ConvertToMPos(pos)).(*cell2DRectGrid)
		npath := &path2DGrid{
			m:     wm,
			cell:  cell,
			first: path.first,
		}
		path.next = npath
		path = npath
	}

	iter := path.First().(*path2DGrid)
	for {
		iter.last = path
		if iter.next == nil {
			break
		}
	}

	return path.First()
}

func (wm *worldMap2DGrid) GetPath(a munfall.Actor, p1, p2 *munfall.WPos) munfall.Path {
	path := &path2DGrid{
		m:    wm,
		cell: wm.CellAt(wm.ConvertToMPos(p1)).(*cell2DRectGrid),
	}
	path.first = path

	lastCell := wm.CellAt(wm.ConvertToMPos(p2)).(*cell2DRectGrid)

	ts := wm.world.GetTraitsImplementing(a, (*traits.OccupySpace)(nil))
	offset := p2.Subtract(p1)
	intersects := false
Outer:
	for _, trait := range ts {
		os := trait.(traits.OccupySpace)
		if os.OutOfBounds(offset) {
			path.last = path
			return path
		}

		for _, osspace := range os.Space() {
			cell := wm.CellAt(wm.ConvertToMPos(osspace.Offset().Add(offset)))
			for _, space := range cell.Space() {
				spacetrait := space.Trait().(traits.OccupySpace)
				if spacetrait.Owner().ActorID() != os.Owner().ActorID() && os.Intersects(spacetrait, offset) {
					intersects = true
					break Outer
				}
			}
		}
	}

	if intersects {
		path.last = path
	} else {
		path.last = &path2DGrid{
			m:     wm,
			cell:  lastCell,
			first: path,
		}
		path.last.last = path.last
		path.next = path.last
	}

	return path
}

func (wm *worldMap2DGrid) InsideMapWPos(pos *munfall.WPos) bool {
	return pos.X < float32(wm.width)*wm.cWidth && pos.X >= 0 &&
		pos.Y < float32(wm.height)*wm.cHeight && pos.Y >= 0
}

func (wm *worldMap2DGrid) InsideMapMPos(pos *munfall.MPos) bool {
	return pos.X < wm.width && pos.X >= 0 &&
		pos.Y < wm.height && pos.Y >= 0
}

func (wm *worldMap2DGrid) ConvertToWPos(m *munfall.MPos) *munfall.WPos {
	return &munfall.WPos{X: wm.cWidth * float32(m.X), Y: wm.cHeight * float32(m.Y)}
}

func (wm *worldMap2DGrid) ConvertToMPos(w *munfall.WPos) *munfall.MPos {
	low := munfall.WPos{}
	high := munfall.WPos{X: wm.cWidth * float32(wm.width), Y: wm.cHeight * float32(wm.height)}
	realw := w.Clamp(&low, &high)
	npos := &munfall.MPos{X: uint(realw.X / wm.cWidth), Y: uint(realw.Y / wm.cHeight)}
	return npos
}

func (wm *worldMap2DGrid) Register(a munfall.Actor) {
	for _, trait := range wm.world.GetTraitsImplementing(a, (*traits.OccupySpace)(nil)) {
		os := trait.(traits.OccupySpace)
		for _, space := range os.Space() {
			cell := wm.CellAt(wm.ConvertToMPos(space.Offset())).(*cell2DRectGrid)
			cell.AddSpace(space)
		}
	}
}

func (wm *worldMap2DGrid) Move(a munfall.Actor, p munfall.Path, percent float32) {
	path, exists := p.(*path2DGrid)
	if !exists {
		munfall.Logger.Panic("Tried using", p, "on a GridWorldMap, it requires a *path2DGrid type.")
	}

	for _, trait := range wm.world.GetTraitsImplementing(a, (*traits.OccupySpace)(nil)) {
		os := trait.(traits.OccupySpace)
		for _, space := range os.Space() {
			cell := wm.CellAt(wm.ConvertToMPos(space.Offset())).(*cell2DRectGrid)
			cell.RemoveSpace(space)
		}
	}

	a.SetPos(path.WPos(percent))

	for _, trait := range wm.world.GetTraitsImplementing(a, (*traits.OccupySpace)(nil)) {
		os := trait.(traits.OccupySpace)
		for _, space := range os.Space() {
			cell := wm.CellAt(wm.ConvertToMPos(space.Offset())).(*cell2DRectGrid)
			cell.AddSpace(space)
		}
	}
}

func (wm *worldMap2DGrid) Deregister(a munfall.Actor) {
	for _, trait := range wm.world.GetTraitsImplementing(a, (*traits.OccupySpace)(nil)) {
		os := trait.(traits.OccupySpace)
		for _, space := range os.Space() {
			cell := wm.CellAt(wm.ConvertToMPos(space.Offset())).(*cell2DRectGrid)
			cell.RemoveSpace(space)
		}
	}
}

type path2DGrid struct {
	m     *worldMap2DGrid
	first *path2DGrid
	last  *path2DGrid

	cell *cell2DRectGrid
	next *path2DGrid
}

func (p *path2DGrid) Cell() munfall.Cell {
	return p.cell
}

func (p *path2DGrid) MPos() *munfall.MPos {
	return p.cell.pos
}

func (p *path2DGrid) WPos(percent float32) *munfall.WPos {
	start := p.m.ConvertToWPos(p.cell.pos)
	if p.next == nil {
		return start
	}

	next := p.m.ConvertToWPos(p.next.cell.pos)
	return start.Subtract(next)
}

func (p *path2DGrid) IsEnd() bool {
	return p.next == nil
}

func (p *path2DGrid) Next() munfall.Path {
	return p.next
}

func (p *path2DGrid) First() munfall.Path {
	return p.first
}

func (p *path2DGrid) Last() munfall.Path {
	return p.last
}

type cell2DRectGrid struct {
	pos            *munfall.MPos
	lc, rc, tc, bc *cell2DRectGrid
	space          []munfall.Space
}

func (c *cell2DRectGrid) AdjacentCells() []munfall.Cell {
	cells := make([]munfall.Cell, 4)
	cells[0] = c.lc
	cells[0] = c.rc
	cells[0] = c.tc
	cells[0] = c.bc
	return cells
}

func (c *cell2DRectGrid) Space() []munfall.Space {
	return c.space
}

func (c *cell2DRectGrid) AddSpace(s munfall.Space) {
	length := len(c.space)
	id := s.Trait().Owner().ActorID()

	if length == 0 {
		c.space = make([]munfall.Space, 1)
		c.space[0] = s
		return
	}

	i := sort.Search(length, func(i int) bool {
		return id == c.space[i].Trait().Owner().ActorID()
	})

	if i == length {
		slice := make([]munfall.Space, 1)
		slice[0] = s
		c.space = append(slice, c.space...)
	} else if i == length-1 {
		c.space = append(c.space, s)
	} else {
		c.space = append(c.space, nil)
		copy(c.space[i+1:], c.space[i:])
		c.space[i] = s
	}
}

func (c *cell2DRectGrid) RemoveSpace(s munfall.Space) {
	length := len(c.space)
	id := s.Trait().Owner().ActorID()
	index := sort.Search(length, func(arg2 int) bool {
		return id == c.space[arg2].Trait().Owner().ActorID()
	})

	if index == -1 {
		return
	} else if index == 0 {
		c.space = c.space[1:]
	} else if index+1 == length {
		c.space = c.space[:index]
	} else {
		c.space = append(c.space[:index-1], c.space[index+1:]...)
	}
}
