// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package logic world.go Defines our world type that runs the game.
package logic

// World container that manages the game world.
type World struct {
	actors          map[uint]*Actor
	traitDictionary *TraitDictionary
	nextActorID     uint
	endtasks        []func()
}

// CreateWorld creates and initializes the World.
func CreateWorld() *World {
	world := &World{actors: make(map[uint]*Actor, 10), endtasks: nil}
	world.traitDictionary = CreateTraitDictionary(world)
	return world
}

// CreateActor Creates an actor in the world and initializes it with the given traits.
func (w *World) CreateActor(traits ...func() Trait) *Actor {
	a := &Actor{actorID: w.nextActorID, world: w}
	w.actors[w.nextActorID] = a
	w.nextActorID++

	for _, trait := range traits {
		w.traitDictionary.addTrait(a, trait())
	}

	notify := w.traitDictionary.GetTraitsImplementing(a, (*TraitAddedNotifier)(nil))
	for _, trait := range notify {
		trait.(TraitAddedNotifier).NotifyAdded(a)
	}

	return a
}

// RemoveActor Removes the given actor from the world.
func (w *World) RemoveActor(a *Actor) {
	if a == nil {
		panic("Trying to remove nil as an Actor!")
	}

	notify := w.traitDictionary.GetTraitsImplementing(a, (*TraitRemovedNotifier)(nil))
	w.traitDictionary.removeActor(a)
	delete(w.actors, a.GetActorID())
	for _, trait := range notify {
		trait.(TraitRemovedNotifier).NotifyRemoved(a)
	}
}

// AddFrameEndTask Adds a task that will be run at the end of the current tick.
func (w *World) AddFrameEndTask(f func()) {
	w.endtasks = append(w.endtasks, f)
}

// TraitDictionary Gets the stored TraitDictionary.
func (w *World) TraitDictionary() *TraitDictionary {
	return w.traitDictionary
}

// Tick ticks all traits on the traitmanager that implement the Tick interface.
func (w *World) Tick(deltaUnit float32) {
	tickers := w.traitDictionary.GetAllTraitsImplementing((*TraitTicker)(nil))
	for _, ticker := range tickers {
		ticker.(TraitTicker).Tick(deltaUnit)
	}

	for _, task := range w.endtasks {
		task()
	}

	w.endtasks = nil
}
