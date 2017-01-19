// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package render mesh.go Defines a mesh type
// used to descripe an object for rendering.
package render

// Mesh type used to hold rendering data.
type Mesh struct {
	points    []float32
	triangles []uint32
}

// ToColor used to compress a color down to a uint32.
func ToColor(r, g, b, a uint32) uint32 {
	return (r&0xff)<<24 | (g&0xff)<<16 | (b&0xff)<<8 | (a & 0xff)
}

// FromColor used to decompress a color up to 4 uint32.
func FromColor(c uint32) (uint32, uint32, uint32, uint32) {
	return (c >> 24) & 0xff, (c >> 16) & 0xff, (c >> 8) & 0xff, c & 0xff
}
