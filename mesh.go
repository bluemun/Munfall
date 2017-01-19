// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package graphics mesh.go Defines a mesh type
// used to descripe an object for rendering.
package graphics

// Mesh type used to hold rendering data.
type Mesh struct {
	vertices []float32
	indices  []uint32
}
