// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package engine interfaces.go Defines interfaces used to prevent circle imports.
package engine

// Trait defines the interface that is used for marking Traits.
type Trait interface {
}

// World defines the interface for the world struct.
type World interface {
	AddFrameEndTask(f func())
	CreateActor(traits ...func() Trait) Actor
	GetTrait(a Actor, i Trait) Trait
	GetTraitsImplementing(a Actor, i Trait) []Trait
	GetAllTraitsImplementing(i Trait) []Trait
	RemoveActor(a Actor)
	ResolveOrder(order *Order)
	Tick(deltaUnit float32)
}

// Actor defines the interface for the actor struct.
type Actor interface {
	World() World
	GetActorID() uint
}

// Renderable interface used to pass data to a renderer.
type Renderable interface {
	Mesh() *Mesh
	Pos() (float32, float32)
	Color() uint32
}