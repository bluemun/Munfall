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
type TraitDictionary interface {
	GetTrait(a Actor, i Trait) (Trait, bool)
	GetTraitsImplementing(a Actor, i Trait) []Trait
	GetAllTraitsImplementing(i Trait) []Trait
}

type traitDictionary struct {
	traits map[reflect.Type]map[uint][]Trait
	world  *world
}

// createTraitDictionary creates and initializes the traitManager.
func createTraitDictionary(w *world) *traitDictionary {
	return &traitDictionary{
		traits: make(map[reflect.Type]map[uint][]Trait),
		world:  w,
	}
}

func (td *traitDictionary) addTrait(a Actor, t Trait) {
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

func (td *traitDictionary) removeActor(a Actor) {
	for _, at := range td.traits {
		delete(at, a.GetActorID())
	}
}

func (td *traitDictionary) GetTrait(a Actor, i Trait) (Trait, bool) {
	engine.Logger.Info(td.traits[reflect.TypeOf(i)])
	engine.Logger.Info(a.GetActorID())
	traits, exists := td.traits[reflect.TypeOf(i)][a.GetActorID()]
	if !exists || len(traits) != 1 {
		return nil, false
	}

	return traits[0], true
}

func (td *traitDictionary) GetTraitsImplementing(a Actor, i Trait) []Trait {
	out := make([]Trait, 0, 1)
	requiredType := reflect.TypeOf(i).Elem()
	for traitType, actorMap := range td.traits {
		traits, exists := actorMap[a.GetActorID()]
		engine.Logger.Debug(traitType, "->", requiredType, "=", traitType.ConvertibleTo(requiredType))
		if exists && traitType.ConvertibleTo(requiredType) {
			for _, trait := range traits {
				out = append(out, trait)
			}
		}
	}

	return out
}

func (td *traitDictionary) GetAllTraitsImplementing(i Trait) []Trait {
	out := make([]Trait, 0, 1)
	requiredType := reflect.TypeOf(i).Elem()
	for traitType, actorMap := range td.traits {
		engine.Logger.Debug(traitType, "->", requiredType, "=", traitType.Implements(requiredType))
		if traitType.Implements(requiredType) {
			for _, y := range actorMap {
				out = append(out, y...)
			}
		}
	}

	return out
}
