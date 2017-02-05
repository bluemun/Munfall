// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package logic traitmanager.go Defines a storage for
// traits supporting easy lookup and add/removal.
package logic

import (
	"reflect"

	"github.com/bluemun/munfall"
)

// TraitDictionary holds traits for easy lookup.
type traitDictionary struct {
	traits map[reflect.Type]map[uint][]munfall.Trait
	world  *world
}

// CreateTraitDictionary creates and initializes the traitManager.
func createTraitDictionary(w *world) *traitDictionary {
	return &traitDictionary{
		traits: make(map[reflect.Type]map[uint][]munfall.Trait),
		world:  w,
	}
}

func (td *traitDictionary) addTrait(a *actor, t munfall.Trait) {
	_, exists := t.(munfall.Trait)
	if !exists {
		munfall.Logger.Panic(t, "does not implement the Trait interface.")
	}

	traittype := reflect.TypeOf(t)
	at, exist := td.traits[traittype]
	if !exist {
		at = make(map[uint][]munfall.Trait)
		td.traits[traittype] = at
	}

	traits, exist := at[a.GetActorID()]
	if !exist {
		at[a.GetActorID()] = []munfall.Trait{t}
	} else {
		at[a.GetActorID()] = append(traits, t)
	}
}

func (td *traitDictionary) removeActor(a *actor) {
	for _, at := range td.traits {
		delete(at, a.GetActorID())
	}
}

// GetTrait gets the given trait from the actor, doesnt support inheritance and
// panics if the trait doesn't exist.
func (td *traitDictionary) GetTrait(a *actor, i interface{}) munfall.Trait {
	t := reflect.TypeOf(i)
	munfall.Logger.Info(td.traits[t])
	munfall.Logger.Info(a.GetActorID())
	traits, exists := td.traits[t][a.GetActorID()]
	if !exists || len(traits) != 1 {
		munfall.Logger.Panic("Trait", t, "doesnt exist on actor", a.GetActorID())
	}

	return traits[0]
}

// GetTraitsImplementing gets all the traits on the given actor that implement
// the given Trait interface.
func (td *traitDictionary) GetTraitsImplementing(a *actor, i interface{}) []munfall.Trait {
	out := make([]munfall.Trait, 0, 1)
	requiredType := reflect.TypeOf(i).Elem()
	for traitType, actorMap := range td.traits {
		munfall.Logger.Debug(traitType, "->", requiredType, "=", traitType.Implements(requiredType), "_", a.GetActorID(), ":", actorMap)
		traits, exists := actorMap[a.GetActorID()]
		if exists && traitType.Implements(requiredType) {
			for _, trait := range traits {
				out = append(out, trait)
			}
		}
	}

	return out
}

// GetAllTraitsImplementing gets all the traits that are in the dictionary
// that implement the given interface.
func (td *traitDictionary) GetAllTraitsImplementing(i interface{}) []munfall.Trait {
	out := make([]munfall.Trait, 0, 1)
	requiredType := reflect.TypeOf(i).Elem()
	munfall.Logger.Debug("Check for traits implementing", requiredType)
	for traitType, actorMap := range td.traits {
		munfall.Logger.Debug(traitType, "->", requiredType, "=", traitType.Implements(requiredType), "_", actorMap)
		if traitType.Implements(requiredType) {
			for _, y := range actorMap {
				out = append(out, y...)
			}
		}
	}

	munfall.Logger.Debug(requiredType, ":=", out)
	return out
}
