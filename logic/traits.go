// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package logic interfaces.go Defines a collection of interfaces
// that allow traits to be called at certain points in time.
package logic

// Trait put inside actors.
type Trait interface {
}

type trait struct {
	traitRegister string
}

// TraitTicker is a trait that gets called every time the world ticks.
type TraitTicker interface {
	Tick(deltaUnit float32)
}

// TraitAddedNotifier is a trait that gets notified when it is added to the world.
type TraitAddedNotifier interface {
	NotifyAdded(owner *Actor)
}

// TraitRemovedNotifier is a trait that gets notified when it is removed from the world.
type TraitRemovedNotifier interface {
	NotifyRemoved(owner *Actor)
}
