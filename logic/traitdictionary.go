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

// TraitDictionary holds traits for easy lookup.
type traitDictionary struct {
	traits map[reflect.Type]map[uint][]engine.Trait
	world  *world
}

// CreateTraitDictionary creates and initializes the traitManager.
func createTraitDictionary(w *world) *traitDictionary {
	return &traitDictionary{
		traits: make(map[reflect.Type]map[uint][]engine.Trait),
		world:  w,
	}
}

func (td *traitDictionary) addTrait(a *actor, t engine.Trait) {
	traittype := reflect.TypeOf(t)
	at, exist := td.traits[traittype]
	if !exist {
		at = make(map[uint][]engine.Trait)
		td.traits[traittype] = at
	}

	traits, exist := at[a.GetActorID()]
	if !exist {
		at[a.GetActorID()] = []engine.Trait{t}
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
func (td *traitDictionary) GetTrait(a *actor, i engine.Trait) engine.Trait {
	t := reflect.TypeOf(i)
	engine.Logger.Info(td.traits[t])
	engine.Logger.Info(a.GetActorID())
	traits, exists := td.traits[t][a.GetActorID()]
	if !exists || len(traits) != 1 {
		engine.Logger.Panic("Trait", t, "doesnt exist on actor", a.GetActorID())
	}

	return traits[0]
}

// GetTraitsImplementing gets all the traits on the given actor that implement
// the given Trait interface.
func (td *traitDictionary) GetTraitsImplementing(a *actor, i engine.Trait) []engine.Trait {
	out := make([]engine.Trait, 0, 1)
	requiredType := reflect.TypeOf(i).Elem()
	for traitType, actorMap := range td.traits {
		engine.Logger.Debug(traitType, "->", requiredType, "=", traitType.Implements(requiredType), "_", a.GetActorID(), ":", actorMap)
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
func (td *traitDictionary) GetAllTraitsImplementing(i engine.Trait) []engine.Trait {
	out := make([]engine.Trait, 0, 1)
	requiredType := reflect.TypeOf(i).Elem()
	engine.Logger.Debug("Check for traits implementing", requiredType)
	for traitType, actorMap := range td.traits {
		engine.Logger.Debug(traitType, "->", requiredType, "=", traitType.Implements(requiredType), "_", actorMap)
		if traitType.Implements(requiredType) {
			for _, y := range actorMap {
				out = append(out, y...)
			}
		}
	}

	engine.Logger.Debug(requiredType, ":=", out)
	return out
}
