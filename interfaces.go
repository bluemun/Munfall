// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package engine interfaces.go Defines interfaces used to prevent circle imports.
package engine

// Trait defines the interface used by every Trait that lives on an Actor.
type Trait interface {
	Initialize(World, Actor, map[string]interface{})
	Owner() Actor
}

// World defines the interface for the world struct.
type World interface {
	AddFrameEndTask(f func())
	GetTrait(a Actor, i interface{}) Trait
	GetTraitsImplementing(a Actor, i interface{}) []Trait
	GetAllTraitsImplementing(i interface{}) []Trait
	RemoveActor(a Actor)
	ResolveOrder(order *Order)
	Tick(deltaUnit float32)
}

// WorldTrait defines the interface used by every Trait that lives on a World.
type WorldTrait interface {
	Initialize(World, map[string]interface{})
}

// Actor defines the interface for the actor struct.
type Actor interface {
	Pos() *WPos
	World() World
	GetActorID() uint
}

// Renderable interface used to pass data to a renderer.
type Renderable interface {
	Mesh() *Mesh
	Pos() (float32, float32)
	Color() uint32
}
