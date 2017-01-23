// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package render camera.go Defines a camera type
// used to control the opengl viewport.
package render

import (
	"github.com/bluemun/engine/graphics/shader"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Camera type used to hold the position of the camera.
type Camera struct {
	X, Y, Width, Height float32
}

var activeCamera *Camera

// Activate activates the camera on the current shader.
func (c *Camera) Activate() {
	activeCamera = c
}

func (c *Camera) use(s *shader.Shader) {
	view := mgl32.Ortho2D(
		c.X-c.Width/2,
		c.X+c.Width/2,
		c.Y-c.Height/2,
		c.Y+c.Height/2)

	gl.UniformMatrix4fv(s.GetUniformLocation("pr"), 1, false, &view[0])
}
