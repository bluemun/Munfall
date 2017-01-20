// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package logic actor.go Defines something
package logic

// Actor temp
type Actor interface {
	Pos() (float32, float32)
	Lookup(t Trait) []Trait
}
