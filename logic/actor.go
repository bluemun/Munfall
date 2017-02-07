// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package logic actor.go Defines something
package logic

import (
	"github.com/bluemun/munfall"
)

// Actor temp
type actor struct {
	actorID uint
	world   *world
	pos     *munfall.WPos
	dead    bool
}

// World returns the world that this actor currently resides in.
func (a *actor) World() munfall.World {
	return (munfall.World)(a.world)
}

// GetActorID gets the actor id of this actor.
func (a *actor) ActorID() uint {
	return a.actorID
}

func (a *actor) Pos() *munfall.WPos {
	return a.pos
}

func (a *actor) SetPos(pos *munfall.WPos) {
	a.pos = pos
}

func (a *actor) Kill() {
	a.dead = true
}

func (a *actor) IsDead() bool {
	return a.dead
}
