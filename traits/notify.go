// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package traits notify.go Defines interfaces that define traits that get notified.
package traits

import (
	"github.com/bluemun/engine"
)

// TraitAddedNotifier is a trait that gets notified when it is added to the world.
type TraitAddedNotifier interface {
	NotifyAdded(owner engine.Actor)
}

// TraitRemovedNotifier is a trait that gets notified when it is removed from the world.
type TraitRemovedNotifier interface {
	NotifyRemoved(owner engine.Actor)
}

// MoveNotifier is called when an Actor is moved on the map.
type MoveNotifier interface {
	NotifyMove(old, new *engine.WPos)
}
