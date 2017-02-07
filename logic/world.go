// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package logic world.go Defines our world type that runs the game.
package logic

import (
	"github.com/bluemun/munfall"
	"github.com/bluemun/munfall/traits"
)

// World container that manages the game world.
type world struct {
	actors          map[uint]*actor
	traitDictionary *traitDictionary
	endtasks        []func()
	wm              munfall.WorldMap
}

// CreateWorld creates and initializes the World.
func CreateWorld(wm munfall.WorldMap) munfall.World {
	world := &world{actors: make(map[uint]*actor, 10), endtasks: nil, wm: wm}
	world.traitDictionary = createTraitDictionary(world)
	wm.Initialize(world)
	return (munfall.World)(world)
}

// AddFrameEndTask adds a task that will be run at the end of the current tick.
func (w *world) AddFrameEndTask(f func()) {
	w.endtasks = append(w.endtasks, f)
}

func (w *world) GetTrait(a munfall.Actor, i interface{}) munfall.Trait {
	return w.traitDictionary.GetTrait(a.(*actor), i)
}

func (w *world) GetTraitsImplementing(a munfall.Actor, i interface{}) []munfall.Trait {
	return w.traitDictionary.GetTraitsImplementing(a.(*actor), i)
}

func (w *world) GetAllTraitsImplementing(i interface{}) []munfall.Trait {
	return w.traitDictionary.GetAllTraitsImplementing(i)
}

// IssueGlobalOrder issues an order to be resolved by every TraitOrderResolver.
func (w *world) IssueGlobalOrder(order *munfall.Order) {
	order.IsGlobal = true
	resolvers := w.traitDictionary.GetAllTraitsImplementing((*traits.TraitOrderResolver)(nil))
	for _, trait := range resolvers {
		trait.(traits.TraitOrderResolver).ResolveOrder(order)
	}
}

// IssueOrder issues an order to be resolved by every TraitOrderResolver on a given Actor.
func (w *world) IssueOrder(a munfall.Actor, order *munfall.Order) {
	order.IsGlobal = false
	resolvers := w.traitDictionary.GetTraitsImplementing(a.(*actor), (*traits.TraitOrderResolver)(nil))
	for _, trait := range resolvers {
		trait.(traits.TraitOrderResolver).ResolveOrder(order)
	}
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

func (w *world) AddToWorld(a munfall.Actor) {
	actor := a.(*actor)
	w.actors[a.ActorID()] = actor

	w.wm.Register(a)
	notify := w.GetTraitsImplementing(a, (*traits.TraitAddedToWorldNotifier)(nil))
	for _, trait := range notify {
		trait.(traits.TraitAddedToWorldNotifier).NotifyAddedToWorld()
	}
}

func (w *world) RemoveFromWorld(a munfall.Actor) {
	if a == nil {
		panic("Trying to remove nil as an Actor!")
	}

	w.wm.Deregister(a)
	notify := w.traitDictionary.GetTraitsImplementing(a.(*actor), (*traits.TraitRemovedFromWorldNotifier)(nil))
	for _, trait := range notify {
		trait.(traits.TraitRemovedFromWorldNotifier).NotifyRemovedFromWorld()
	}
}

func (w *world) WorldMap() munfall.WorldMap {
	return w.wm
}
