// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package traits notify.go Defines interfaces that define traits that get notified.
package traits

import (
	"github.com/bluemun/munfall"
)

// TraitAddedNotifier is a trait that gets notified when it is added to the world.
type TraitAddedNotifier interface {
	munfall.Trait
	NotifyAdded(owner munfall.Actor)
}

// TraitRemovedNotifier is a trait that gets notified when it is removed from the world.
type TraitRemovedNotifier interface {
	munfall.Trait
	NotifyRemoved(owner munfall.Actor)
}

// MoveNotifier is called when an Actor is moved on the map.
type MoveNotifier interface {
	munfall.Trait
	NotifyMove(old, new *munfall.WPos)
}
