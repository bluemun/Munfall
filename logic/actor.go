// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package logic actor.go Defines something
package logic

import (
	"github.com/bluemun/engine"
)

// Actor temp
type actor struct {
	actorID uint
	world   *world
}

// World returns the world that this actor currently resides in.
func (a *actor) World() engine.World {
	return (engine.World)(a.world)
}

// GetActorID gets the actor id of this actor.
func (a *actor) GetActorID() uint {
	return a.actorID
}
