// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package graphics camera.go Defines a camera type
// used to control the opengl viewport.
package graphics

import (
	"github.com/go-gl/gl/v3.3-core/gl"
)

// Camera type used to hold the position of the camera.
type Camera struct {
	x, y float32
}

// Activate activates the camera on the current shader.
func (*Camera) Activate() {
	gl.UniformMatrix4fv(0, 0, false, nil)
}
