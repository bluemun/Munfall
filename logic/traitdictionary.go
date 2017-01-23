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
type TraitDictionary struct {
	traits map[reflect.Type]map[uint][]Trait
	world  *World
}

// CreateTraitDictionary creates and initializes the traitManager.
func CreateTraitDictionary(w *World) *TraitDictionary {
	return &TraitDictionary{
		traits: make(map[reflect.Type]map[uint][]Trait),
		world:  w,
	}
}

func (td *TraitDictionary) addTrait(a *Actor, t Trait) {
	traittype := reflect.TypeOf(t)
	at, exist := td.traits[traittype]
	if !exist {
		at = make(map[uint][]Trait)
		td.traits[traittype] = at
	}

	traits, exist := at[a.GetActorID()]
	if !exist {
		at[a.GetActorID()] = []Trait{t}
	} else {
		at[a.GetActorID()] = append(traits, t)
	}
}

func (td *TraitDictionary) removeActor(a *Actor) {
	for _, at := range td.traits {
		delete(at, a.GetActorID())
	}
}

// GetTrait Gets the given trait from the actor, doesnt support inheritance and
// panics if the trait doesn't exist.
func (td *TraitDictionary) GetTrait(a *Actor, i Trait) (Trait, bool) {
	engine.Logger.Info(td.traits[reflect.TypeOf(i)])
	engine.Logger.Info(a.GetActorID())
	traits, exists := td.traits[reflect.TypeOf(i)][a.GetActorID()]
	if !exists || len(traits) != 1 {
		return nil, false
	}

	return traits[0], true
}

// GetTraitsImplementing Gets all the traits on the given actor that implement
// the given Trait interface.
func (td *TraitDictionary) GetTraitsImplementing(a *Actor, i Trait) []Trait {
	out := make([]Trait, 0, 1)
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

// GetAllTraitsImplementing Gets all the traits that are in the dictionary
// that implement the given interface.
func (td *TraitDictionary) GetAllTraitsImplementing(i Trait) []Trait {
	out := make([]Trait, 0, 1)
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
