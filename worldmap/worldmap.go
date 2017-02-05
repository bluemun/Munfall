// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package worldmap worldmap.go Defines the interfaces and structs for the worldmap.
package worldmap

import (
	"github.com/bluemun/munfall"
)

// WorldMap is the interface for the world map.
type WorldMap interface {
	CellAt(pos *MPos) Cell
	GetPath(p1, p2 *MPos) Path
	ConvertToWPos(m *MPos) *munfall.WPos
	ConvertToMPos(w *munfall.WPos) *MPos
}

// Path is an iterator that defines a path through the world map.
type Path interface {
	Cell() Cell
	MPos() *MPos
	WPos(percent float32) *munfall.WPos

	IsEnd() bool
	Next() Path

	First() Path
	Last() Path
}

// Cell is the building block of the world map.
type Cell interface {
	AdjacentCells() []Cell
}

// MPos map position.
type MPos struct {
	X, Y uint
}
