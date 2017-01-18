// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package graphics renderer.go Defines a renderer in the graphics package.
package graphics

import (
	"github.com/go-gl/gl/v3.3-core/gl"
)

// Renderer used by the graphics library to draw.
type Renderer struct {
	offset       int32
	vertexArray  uint32
	vertexBuffer uint32
	indexBuffer  uint32
	s            *shader
}

const int32size int32 = 4

// CreateRenderer used to create a renderer object correctly.
func CreateRenderer(s *shader) *Renderer {
	r := new(Renderer)
	r.s = s
	do(func() {
		r.s.use()
		gl.GenVertexArrays(1, &r.vertexArray)
		checkGLError()
		gl.BindVertexArray(r.vertexArray)
		checkGLError()

		gl.GenBuffers(1, &r.vertexBuffer)
		checkGLError()
		gl.BindBuffer(gl.ARRAY_BUFFER, r.vertexBuffer)
		checkGLError()
		gl.BufferData(gl.ARRAY_BUFFER, (int)(2000*4*3*int32size), nil, gl.DYNAMIC_DRAW)
		checkGLError()

		vertAttrib := r.s.getAttributeLocation("vertex")
		logger.Info("Vertex attribute location: ", vertAttrib)
		gl.EnableVertexAttribArray(vertAttrib)
		checkGLError()
		gl.VertexAttribPointer(vertAttrib, 3, gl.INT, false, 3*int32size, gl.PtrOffset(0))
		checkGLError()

		gl.GenBuffers(1, &r.indexBuffer)
		checkGLError()
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, r.indexBuffer)
		checkGLError()

		var indices [12000]uint32
		var j uint32
		for i := 0; i < 12000; i += 6 {
			indices[i] = j
			indices[i+1] = j + 1
			indices[i+2] = j + 2
			indices[i+3] = j + 1
			indices[i+4] = j + 2
			indices[i+5] = j + 3
			j += 4
		}

		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, (int)(12000*int32size), gl.Ptr(&indices[0]), gl.DYNAMIC_DRAW)
		checkGLError()

		gl.BindFragDataLocation((uint32)(*r.s), 0, gl.Str("outputColor\x00"))
		checkGLError()
	})

	return r
}

// Begin starts the rendering procedure.
func (r *Renderer) Begin() {
	r.offset = 0
	do(func() {
		r.s.use()
		gl.BindBuffer(gl.ARRAY_BUFFER, r.vertexBuffer)
		checkGLError()
	})
}

// DrawRectangle draws a rectangle using the given values, x and y point to the top-left corner.
func (r *Renderer) DrawRectangle(x, y, w, h int32) {
	array := [12]int32{
		x, y, 0,
		x + w, y, 0,
		x, y + h, 0,
		x + w, y + h, 0,
	}
	do(func() {
		gl.BufferSubData(gl.ARRAY_BUFFER, (int)(r.offset*12*int32size), (int)(12*int32size), gl.Ptr(&array[0]))
		checkGLError()
	})
	r.offset++
	if r.offset == 2000 {
		r.Flush()
		r.offset = 0
	}
}

// Flush flushes all the draw calls that have been called on this renderer to the window
func (r *Renderer) Flush() {
	do(func() {
		gl.BindVertexArray(r.vertexArray)
		checkGLError()
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, r.indexBuffer)
		checkGLError()
		gl.DrawElements(gl.TRIANGLES, r.offset*6, gl.UNSIGNED_INT, nil)
		checkGLError()
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
		checkGLError()
	})
}

// End ends the rendering procedure.
func (r *Renderer) End() {
	do(func() {
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		checkGLError()
	})
}
