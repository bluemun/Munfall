// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package munfall Defines primitives that will be used all around the engine,
// they are here to prevent circle imports.
package munfall

import (
	"math"
)

// Order wraps an order that gets passed around by the order generator.
type Order struct {
	Order    string
	Value    interface{}
	IsGlobal bool
}

// Mesh type used to hold rendering data.
type Mesh struct {
	Points    []float32
	Triangles []uint32
}

// MPos map position.
type MPos struct {
	X, Y uint
}

// WPos world position.
type WPos struct {
	X, Y, Z float32
}

// Add returns a WPos that is the sum of the 2 given WPos.
func (w *WPos) Add(other *WPos) *WPos {
	return &WPos{
		X: w.X + other.X,
		Y: w.Y + other.Y,
		Z: w.Z + other.Z,
	}
}

// Subtract returns a WPos that is the difference of the 2 given WPos.
func (w *WPos) Subtract(other *WPos) *WPos {
	return &WPos{
		X: w.X - other.X,
		Y: w.Y - other.Y,
		Z: w.Z - other.Z,
	}
}

// Clamp clamps the WPos to the given vector values.
func (w *WPos) Clamp(low, high *WPos) *WPos {
	return &WPos{
		X: float32(math.Min(math.Max(float64(w.X), float64(low.X)), float64(high.X))),
		Y: float32(math.Min(math.Max(float64(w.Y), float64(low.Y)), float64(high.Y))),
		Z: float32(math.Min(math.Max(float64(w.Z), float64(low.Z)), float64(high.Z))),
	}
}

// Vector2f holds 2 float values.
type Vector2f struct {
	X, Y float32
}
