// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package engine Defines primitives that will be used all around the engine,
// they are here to prevent circle imports.
package engine

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
