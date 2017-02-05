// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package logic actorregistry.go Defines a container for actor definitions so they can be easily added to the world.
package logic

import (
	"reflect"

	"github.com/bluemun/munfall"
	"github.com/bluemun/munfall/traits"
)

// ActorDefinition holds all the information needed to make an Actor.
type ActorDefinition struct {
	Name       string
	parameters map[string]map[string]interface{}
}

// CreateActorDefinition creates a new actor definition, initialzed with the given name.
func CreateActorDefinition(name string) *ActorDefinition {
	return &ActorDefinition{Name: name, parameters: make(map[string]map[string]interface{})}
}

// AddTrait adds a trait that will be constructed and added to the actor
// when this ActorDefinition is used to create an actor.
func (ad *ActorDefinition) AddTrait(traitName string) {
	_, exists := ad.parameters[traitName]
	if !exists {
		ad.parameters[traitName] = make(map[string]interface{}, 1)
	}
}

// AddParameter adds a parameter for the given trait that is provided to it every
// time this ActorDefinition is used to create an actor.
func (ad *ActorDefinition) AddParameter(traitName, parameterName string, value interface{}) {
	traitParams, exists := ad.parameters[traitName]
	if !exists {
		munfall.Logger.Panic("Tried adding a parameter to trait", traitName,
			"on ActorDefinition", ad.Name, ": ActorDefinition does not define", traitName)
	}

	traitParams[parameterName] = value
}

// TraitCreate holds all the information needed to create a trait.
type TraitCreate struct {
	Name       string
	Type       reflect.Type
	Parameters map[string]interface{}
}

// ActorRegistry contains definitions for actors.
type ActorRegistry struct {
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
func (ar *ActorRegistry) CreateActor(name string, runtimeParameters map[string]interface{}, w munfall.World) {
	world := w.(*world)
	params := ar.builders[name]

	a := &actor{actorID: world.nextActorID, world: world}
	world.actors[world.nextActorID] = a
	world.nextActorID++

	for name, parameter := range params.parameters {
		obj := reflect.New(ar.definitions[name])

		//munfall.Logger.Info(obj)
		trait := obj.Interface().(munfall.Trait)

		if runtimeParameters == nil {
			trait.Initialize(w, a, parameter)
		} else {
			np := make(map[string]interface{}, len(runtimeParameters)+len(parameter))
			for key, value := range parameter {
				np[key] = value
			}

			for key, value := range runtimeParameters {
				np[key] = value
			}

			trait.Initialize(w, a, np)
		}

		world.traitDictionary.addTrait(a, trait)
	}

	notify := world.GetTraitsImplementing(a, (*traits.TraitAddedNotifier)(nil))
	for _, trait := range notify {
		trait.(traits.TraitAddedNotifier).NotifyAdded((munfall.Actor)(a))
	}
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
