// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package logic world.go Defines our world type that runs the game.
package logic

// World container that manages the game world.
type World interface {
	AddFrameEndTask(f func())
	CreateActor(traits ...func() Trait) Actor
	RemoveActor(a Actor)
	TraitDictionary() TraitDictionary
	Tick(deltaUnit float32)
}

type world struct {
	actors          map[uint]Actor
	traitDictionary *traitDictionary
	nextActorID     uint
	endtasks        []func()
}

// CreateWorld creates and initializes the World.
func CreateWorld() World {
	world := &world{actors: make(map[uint]Actor, 10), endtasks: nil}
	world.traitDictionary = createTraitDictionary(world)
	return (World)(world)
}

func (w *world) CreateActor(traits ...func() Trait) Actor {
	a := &actor{actorID: w.nextActorID, world: w}
	w.actors[w.nextActorID] = a
	w.nextActorID++

	for _, trait := range traits {
		w.traitDictionary.addTrait(a, trait())
	}

	notify := w.traitDictionary.GetTraitsImplementing(a, (*TraitNotifyAdded)(nil))
	for _, trait := range notify {
		trait.(TraitNotifyAdded).NotifyAdded(a)
	}

	return (Actor)(a)
}

func (w *world) RemoveActor(a Actor) {
	if a == nil {
		panic("Trying to remove nil as an Actor!")
	}

	notify := w.traitDictionary.GetTraitsImplementing(a, (*TraitNotifyRemoved)(nil))
	w.traitDictionary.removeActor(a)
	delete(w.actors, a.GetActorID())
	for _, trait := range notify {
		trait.(TraitNotifyRemoved).NotifyRemoved(a)
	}
}

func (w *world) AddFrameEndTask(f func()) {
	w.endtasks = append(w.endtasks, f)
}

func (w *world) TraitDictionary() TraitDictionary {
	return (TraitDictionary)(w.traitDictionary)
}

// Tick ticks all traits on the traitmanager that implement the Tick interface.
func (w *world) Tick(deltaUnit float32) {
	tickers := w.traitDictionary.GetAllTraitsImplementing((*TraitTick)(nil))
	for _, ticker := range tickers {
		ticker.(TraitTick).Tick(deltaUnit)
	}

	for _, task := range w.endtasks {
		task()
	}

	w.endtasks = nil
}
