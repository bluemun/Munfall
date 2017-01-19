// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package logic world.go Defines our world type that runs the game.
package logic

import (
	"reflect"
)

// World container that manages the game world.
type World struct {
	Traitmanager TraitManager
}

// CreateWorld creates and initializes the World.
func CreateWorld() *World {
	world := new(World)
	world.Traitmanager = (TraitManager)(createTraitManager(world))
	return world
}

// Tick ticks all traits on the traitmanager that implement the Tick interface.
func (w *World) Tick(deltaUnit float32) {
	var t1 = reflect.TypeOf((*TraitTick)(nil)).Elem()
	tickers := w.Traitmanager.Lookup(t1)
	for _, ticker := range tickers {
		ticker.(TraitTick).Tick()
	}
}
