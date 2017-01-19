// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package engine actormanager.go Defines our actor manager.
package engine

// Actor is the interface by which the ActorManager talks to
type Actor interface {
	Pos() (float32, float32)
	Lookup(t Trait) []Trait
}

// Trait put inside actors.
type Trait interface{}

// TraitTick is a trait that gets called every time the world ticks.
type TraitTick interface {
	Trait
	Tick()
}

// ActorManager holds actors for easy lookup.
type ActorManager struct {
	actors []Actor
}

// AddActor holds actors for easy lookup.
func (a *ActorManager) AddActor(act Actor) {
	a.actors = append(a.actors, act)
}

// Tick ticks all traits on stored actors that implement the Tick interface.
func (a *ActorManager) Tick(deltaUnit float32) {
	var t TraitTick
	tickers := a.Lookup(t)
	for _, ticker := range tickers {
		ticker.(TraitTick).Tick()
	}
}

// Lookup returns every actor that implements the given interface
func (a *ActorManager) Lookup(t interface{}) []interface{} {
	var out []interface{}
	for _, actor := range a.actors {
		traits := actor.Lookup(t)
		for _, trait := range traits {
			out = append(out, trait)
		}
	}
	return out
}
