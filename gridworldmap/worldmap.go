// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package gridworldmap worldmap.go Defines the map that actors live and move on.
package gridworldmap

import (
	"github.com/bluemun/munfall"
	"github.com/bluemun/munfall/worldmap"
)

type worldMap2DGrid struct {
	cWidth, cHeight float32
	width, height   uint
	grid            []*cell2DRectGrid
}

// CreateGridWorldMap creates a 2D implementation of a WorldMap.
func CreateGridWorldMap(width, height uint, cellWidth, cellHeight float32) worldmap.WorldMap {
	grid := make([]*cell2DRectGrid, 0, width*height)

	var init uint
	for y := init; y < height; y++ {
		for x := init; x < width; x++ {
			grid[x+y*width] = &cell2DRectGrid{pos: &worldmap.MPos{X: x, Y: y}}
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

func (wm *worldMap2DGrid) CellAt(pos *worldmap.MPos) worldmap.Cell {
	return wm.grid[pos.X+pos.Y*wm.width]
}

func (wm *worldMap2DGrid) GetPath(p1, p2 *worldmap.MPos) worldmap.Path {
	return nil
}

func (wm *worldMap2DGrid) ConvertToWPos(m *worldmap.MPos) *munfall.WPos {
	return &munfall.WPos{X: wm.cWidth * float32(m.X), Y: wm.cHeight * float32(m.Y)}
}

func (wm *worldMap2DGrid) ConvertToMPos(w *munfall.WPos) *worldmap.MPos {
	return &worldmap.MPos{X: uint(w.X / wm.cWidth), Y: uint(w.Y / wm.cHeight)}
}

type path2DGrid struct {
	m     *worldMap2DGrid
	first *path2DGrid
	last  *path2DGrid

	cell *cell2DRectGrid
	next *path2DGrid
}

func (p *path2DGrid) Cell() worldmap.Cell {
	return p.cell
}

func (p *path2DGrid) MPos() *worldmap.MPos {
	return p.cell.pos
}

func (p *path2DGrid) WPos(percent float32) *munfall.WPos {
	start := p.m.ConvertToWPos(p.cell.pos)
	next := p.m.ConvertToWPos(p.next.cell.pos)
	movement := start.Vector(next)

	// TODO: Implement a heightmap somewhere so we can give a real z value here.
	return &munfall.WPos{
		X: start.X + movement.X*percent,
		Y: start.Y * movement.Y * percent,
		Z: 0,
	}
}

func (p *path2DGrid) IsEnd() bool {
	return p.next == nil
}

func (p *path2DGrid) Next() worldmap.Path {
	return p.next
}

func (p *path2DGrid) First() worldmap.Path {
	return p.first
}

func (p *path2DGrid) Last() worldmap.Path {
	return p.last
}

type cell2DRectGrid struct {
	pos            *worldmap.MPos
	lc, rc, tc, bc *cell2DRectGrid
}

func (c *cell2DRectGrid) AdjacentCells() []worldmap.Cell {
	cells := make([]worldmap.Cell, 4)
	cells[0] = c.lc
	cells[0] = c.rc
	cells[0] = c.tc
	cells[0] = c.bc
	return cells
}
