// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package logic actor.go Defines something
package logic

// Actor temp
type Actor struct {
	actorID uint
	world   *World
}

// World returns the world that this actor currently resides in.
func (a *Actor) World() *World {
	return a.world
}

// GetActorID Gets the actor id of this actor.
func (a *Actor) GetActorID() uint {
	return a.actorID
}
