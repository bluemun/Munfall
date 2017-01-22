// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package logic actor.go Defines something
package logic

// Actor temp
type Actor interface {
	World() World
	GetActorID() uint
	Lookup(t Trait) []Trait
}

type actor struct {
	actorID uint
	world   World
}

//Pos() (float32, float32)

func (a *actor) World() World {
	return a.world
}

// Lookup gets all traits implementing given trait.
func (a *actor) GetActorID() uint {
	return a.actorID
}

// Lookup gets all traits implementing given trait.
func (a *actor) Lookup(t Trait) []Trait {
	return nil
}
