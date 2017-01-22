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

// TraitTick is a trait that gets called every time the world ticks.
type TraitTick interface {
	Tick(deltaUnit float32)
}

// TraitNotifyAdded is a trait that gets notified when it is added to the world.
type TraitNotifyAdded interface {
	NotifyAdded(owner Actor)
}

// TraitNotifyRemoved is a trait that gets notified when it is removed from the world.
type TraitNotifyRemoved interface {
	NotifyRemoved(owner Actor)
}
