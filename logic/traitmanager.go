// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package logic traitmanager.go Defines a storage for
// traits supporting easy lookup and add/removal.
package logic

import (
	"reflect"

	"github.com/bluemun/engine"
)

// TraitManager holds traits for easy lookup.
type TraitManager interface {
	AddTrait(tr Trait)
	RemoveTrait(tr Trait)
	Lookup(t reflect.Type) []interface{}
}

// Implementation of the TraitManager interface
type traitManager struct {
	traits map[Trait]bool
	world  *World
}

// createTraitManager creates and initializes the traitManager.
func createTraitManager(w *World) *traitManager {
	tm := new(traitManager)
	tm.traits = make(map[Trait]bool)
	tm.world = w
	return tm
}

// AddTrait adds a trait to the manager.
func (tm *traitManager) AddTrait(tr Trait) {
	tm.traits[tr] = true
	engine.Logger.Info(tm.traits)
	notify, found := tr.(TraitNotifyAdded)
	if found {
		notify.NotifyAdded(tm.world)
	}
}

// RemoveTrait removes a trait from the manager.
func (tm *traitManager) RemoveTrait(tr Trait) {
	delete(tm.traits, tr)
	notify, found := tr.(TraitNotifyRemoved)
	if found {
		notify.NotifyRemoved(tm.world)
	}
}

// Lookup returns every trait that implements the given interface
func (tm *traitManager) Lookup(t reflect.Type) []interface{} {
	var out []interface{}
	for trait := range tm.traits {
		if reflect.TypeOf(trait).Implements(t) {
			out = append(out, trait)
		}
	}
	return out
}
