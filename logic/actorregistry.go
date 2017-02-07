// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package logic actorregistry.go Defines a container for actor definitions so they can be easily added to the world.
package logic

import (
	"reflect"

	"github.com/bluemun/munfall"
)

// ActorDefinition holds all the information needed to make an Actor.
type ActorDefinition struct {
	Name   string
	traits []*TraitDefinition
}

// TraitDefinition used for constructing a trait.
type TraitDefinition struct {
	Type       string
	parameters map[string]interface{}
}

// CreateTraitDefinition creates a struct for creating traits.
func CreateTraitDefinition(Type string) *TraitDefinition {
	return &TraitDefinition{
		Type:       Type,
		parameters: make(map[string]interface{}),
	}
}

// AddParameter adds a parameter to the trait, returns the definintion it is called
// on for easy chaining.
func (td *TraitDefinition) AddParameter(parameterName string, value interface{}) *TraitDefinition {
	td.parameters[parameterName] = value
	return td
}

// CreateActorDefinition creates a new actor definition, initialzed with the given name.
func CreateActorDefinition(name string) *ActorDefinition {
	return &ActorDefinition{Name: name, traits: make([]*TraitDefinition, 0)}
}

// AddTrait adds a trait that will be constructed and added to the actor
// when this ActorDefinition is used to create an actor.
func (ad *ActorDefinition) AddTrait(def *TraitDefinition) {
	ad.traits = append(ad.traits, def)
}

// TraitCreate holds all the information needed to create a trait.
type TraitCreate struct {
	Name       string
	Type       reflect.Type
	Parameters map[string]interface{}
}

// ActorRegistry contains definitions for actors.
type ActorRegistry struct {
	nextID      uint
	definitions map[string]reflect.Type
	builders    map[string]*ActorDefinition
}

// CreateActorRegistry creates and initializes an ActorRegistry.
func CreateActorRegistry() *ActorRegistry {
	ar := &ActorRegistry{
		definitions: make(map[string]reflect.Type),
		builders:    make(map[string]*ActorDefinition),
	}

	return ar
}

// CreateActor creates an actor in the given world by using the trait parameters
// registered to the given name and the provided runtime parameters.
func (ar *ActorRegistry) CreateActor(name string, runtimeParameters map[string]interface{}, w munfall.World, addToWorld bool) munfall.Actor {
	world := w.(*world)
	params := ar.builders[name]

	a := &actor{actorID: ar.nextID, pos: &munfall.WPos{}, world: world}
	ar.nextID++

	for _, traitdef := range params.traits {
		obj := reflect.New(ar.definitions[traitdef.Type])
		trait := obj.Interface().(munfall.Trait)

		if runtimeParameters == nil {
			trait.Initialize(w, a, traitdef.parameters)
		} else {
			np := make(map[string]interface{}, len(runtimeParameters)+len(traitdef.parameters))
			for key, value := range traitdef.parameters {
				np[key] = value
			}

			for key, value := range runtimeParameters {
				np[key] = value
			}

			trait.Initialize(w, a, np)
		}

		world.traitDictionary.addTrait(a, trait)
	}

	if addToWorld {
		world.AddToWorld(a)
	}

	return a
}

// DisposeActor disposes of all the traits from the world.
func (ar *ActorRegistry) DisposeActor(a munfall.Actor, w munfall.World) {
	world := w.(*world)
	world.RemoveFromWorld(a)
}

// RegisterTrait adds a trait type as a candidate for creation, panics if it already exists.
func (ar *ActorRegistry) RegisterTrait(name string, t interface{}) {
	_, exists := ar.definitions[name]
	if exists {
		munfall.Logger.Panic("Trait:", name, "already exists in the trait registry.")
	}

	ar.definitions[name] = reflect.TypeOf(t).Elem()
}

// RegisterActor adds trait parameters to the registered actor name, used by the
// CreateActor method to
func (ar *ActorRegistry) RegisterActor(definition *ActorDefinition) {
	_, exists := ar.builders[definition.Name]
	if exists {
		munfall.Logger.Panic("An actor with the name", definition.Name, "has already been registered.")
	}

	ar.builders[definition.Name] = definition
}
