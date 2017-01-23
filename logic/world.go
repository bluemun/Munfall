// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package logic world.go Defines our world type that runs the game.
package logic

import (
	"github.com/bluemun/engine"
	"github.com/bluemun/engine/traits"
)

// World container that manages the game world.
type world struct {
	actors          map[uint]*actor
	traitDictionary *traitDictionary
	nextActorID     uint
	endtasks        []func()
}

// CreateWorld creates and initializes the World.
func CreateWorld() engine.World {
	world := &world{actors: make(map[uint]*actor, 10), endtasks: nil}
	world.traitDictionary = createTraitDictionary(world)
	return (engine.World)(world)
}

// AddFrameEndTask adds a task that will be run at the end of the current tick.
func (w *world) AddFrameEndTask(f func()) {
	w.endtasks = append(w.endtasks, f)
}

// CreateActor creates an actor in the world and initializes it with the given traits.
func (w *world) CreateActor(ptraits ...func() engine.Trait) engine.Actor {
	a := &actor{actorID: w.nextActorID, world: w}
	w.actors[w.nextActorID] = a
	w.nextActorID++

	for _, trait := range ptraits {
		w.traitDictionary.addTrait(a, trait())
	}

	notify := w.traitDictionary.GetTraitsImplementing(a, (*traits.TraitAddedNotifier)(nil))
	for _, trait := range notify {
		trait.(traits.TraitAddedNotifier).NotifyAdded((engine.Actor)(a))
	}

	return a
}

func (w *world) GetTrait(a engine.Actor, i engine.Trait) engine.Trait {
	return w.traitDictionary.GetTrait(a.(*actor), i)
}

func (w *world) GetTraitsImplementing(a engine.Actor, i engine.Trait) []engine.Trait {
	return w.traitDictionary.GetTraitsImplementing(a.(*actor), i)
}

func (w *world) GetAllTraitsImplementing(i engine.Trait) []engine.Trait {
	return w.traitDictionary.GetAllTraitsImplementing(i)
}

// RemoveActor removes the given actor from the world.
func (w *world) RemoveActor(a engine.Actor) {
	if a == nil {
		panic("Trying to remove nil as an Actor!")
	}

	notify := w.traitDictionary.GetTraitsImplementing(a.(*actor), (*traits.TraitRemovedNotifier)(nil))
	w.traitDictionary.removeActor(a.(*actor))
	delete(w.actors, a.GetActorID())
	for _, trait := range notify {
		trait.(traits.TraitRemovedNotifier).NotifyRemoved(a)
	}
}

// ResolveOrder bla.
func (w *world) ResolveOrder(order *engine.Order) {
	//w.traitDictionary.GetAllTraitsImplementing()
}

// Tick ticks all traits on the traitmanager that implement the Tick interface.
func (w *world) Tick(deltaUnit float32) {
	tickers := w.traitDictionary.GetAllTraitsImplementing((*traits.TraitTicker)(nil))
	for _, ticker := range tickers {
		ticker.(traits.TraitTicker).Tick(deltaUnit)
	}

	for _, task := range w.endtasks {
		task()
	}

	w.endtasks = nil
}
