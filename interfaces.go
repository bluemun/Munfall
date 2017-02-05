// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package munfall interfaces.go Defines interfaces used to prevent circle imports.
package munfall

// Trait defines the interface used by every Trait that lives on an Actor.
type Trait interface {
	Initialize(World, Actor, map[string]interface{})
	Owner() Actor
}

// World defines the interface for the world struct.
type World interface {
	AddFrameEndTask(f func())
	Tick(deltaUnit float32)

	GetTrait(a Actor, i interface{}) Trait
	GetTraitsImplementing(a Actor, i interface{}) []Trait
	GetAllTraitsImplementing(i interface{}) []Trait

	AddToWorld(a Actor)
	RemoveFromWorld(a Actor)

	IssueGlobalOrder(order *Order)
	IssueOrder(a Actor, order *Order)

	WorldMap() WorldMap
}

// WorldMap is the interface for the world map.
type WorldMap interface {
	Initialize(World)

	InsideMapWPos(pos *WPos) bool
	InsideMapMPos(pos *MPos) bool

	CellAt(*MPos) Cell
	GetPath(a Actor, p1, p2 *MPos) Path

	ConvertToWPos(*MPos) *WPos
	ConvertToMPos(*WPos) *MPos

	Register(Actor)
	Move(Actor, Path, float32)
	Deregister(Actor)
}

// Path is an iterator that defines a path through the world map.
type Path interface {
	Cell() Cell
	MPos() *MPos
	WPos(percent float32) *WPos

	IsEnd() bool
	Next() Path

	First() Path
	Last() Path
}

// Cell is the building block of the world map.
type Cell interface {
	AdjacentCells() []Cell
	Space() []Space
}

// Space defines a space that can be used to
type Space interface {
	Trait() Trait
	Initialize(Trait)
	Offset() *WPos
	Intersects(other Space) bool
}

// Actor defines the interface for the actor struct.
type Actor interface {
	Pos() *WPos
	SetPos(pos *WPos)
	World() World
	ActorID() uint
}

// Renderable interface used to pass data to a renderer.
type Renderable interface {
	Mesh() *Mesh
	Pos() *WPos
	Color() uint32
}
