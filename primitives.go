// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package munfall Defines primitives that will be used all around the engine,
// they are here to prevent circle imports.
package munfall

// Order wraps an order that gets passed around by the order generator.
type Order struct {
	Order string
	Value interface{}
}

// Mesh type used to hold rendering data.
type Mesh struct {
	Points    []float32
	Triangles []uint32
}

// WPos world position.
type WPos struct {
	X, Y, Z float32
}

// Vector returns the vector between the two provided WPos
func (w *WPos) Vector(other *WPos) *Vector3f {
	return &Vector3f{
		X: w.X - other.X,
		Y: w.Y - other.Y,
		Z: w.Z - other.Z,
	}
}

// Vector3f holds 3 float values.
type Vector3f struct {
	X, Y, Z float32
}

// Vector2f holds 2 float values.
type Vector2f struct {
	X, Y float32
}
