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
	vertexOffset, indexOffset, vertexBufferSize, indexBufferSize int
	vertexArray, vertexBuffer, indexBuffer                       uint32
	s                                                            *shader
}

const int32Size = 4
const float32Size = 4
const vertexSize = 3 * float32Size

// CreateRenderer used to create a renderer object correctly.
func CreateRenderer(s *shader, vertexBufferSize, indexBufferSize int) *Renderer {
	r := new(Renderer)
	r.s = s
	r.vertexBufferSize = vertexBufferSize
	r.indexBufferSize = indexBufferSize
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
		gl.BufferData(gl.ARRAY_BUFFER, (int)(10000*vertexSize), nil, gl.DYNAMIC_DRAW)
		checkGLError()

		vertAttrib := r.s.getAttributeLocation("vertex")
		logger.Info("Vertex attribute location: ", vertAttrib)
		gl.EnableVertexAttribArray(vertAttrib)
		checkGLError()
		gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, vertexSize, gl.PtrOffset(0))
		checkGLError()

		gl.GenBuffers(1, &r.indexBuffer)
		checkGLError()
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, r.indexBuffer)
		checkGLError()

		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, (int)(10000*int32Size), nil, gl.DYNAMIC_DRAW)
		checkGLError()

		gl.BindFragDataLocation((uint32)(*r.s), 0, gl.Str("outputColor\x00"))
		checkGLError()
	})

	return r
}

// Begin starts the rendering procedure.
func (r *Renderer) Begin() {
	r.indexOffset, r.vertexOffset = 0, 0
	do(func() {
		r.s.use()
		gl.BindBuffer(gl.ARRAY_BUFFER, r.vertexBuffer)
		checkGLError()
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, r.indexBuffer)
		checkGLError()
	})
}

// DrawRectangle draws a rectangle using the given values, x and y point to the top-left corner.
func (r *Renderer) DrawRectangle(x, y, w, h float32) {
	vertices := [12]float32{
		x, y, 0,
		x + w, y, 0,
		x, y + h, 0,
		x + w, y + h, 0,
	}
	indices := [6]uint32{
		0 + uint32(r.indexOffset), 1 + uint32(r.indexOffset), 2 + uint32(r.indexOffset),
		1 + uint32(r.indexOffset), 2 + uint32(r.indexOffset), 3 + uint32(r.indexOffset),
	}

	if r.vertexOffset+12 >= r.vertexBufferSize || r.indexOffset+6 >= r.indexBufferSize {
		r.Flush()
		r.indexOffset, r.vertexOffset = 0, 0
	}

	do(func() {
		gl.BufferSubData(gl.ARRAY_BUFFER, (r.vertexOffset)*float32Size, 12*float32Size, gl.Ptr(&vertices[0]))
		checkGLError()
		gl.BufferSubData(gl.ELEMENT_ARRAY_BUFFER, (r.indexOffset)*int32Size, 6*int32Size, gl.Ptr(&indices[0]))
		checkGLError()
	})
	r.vertexOffset += 12
	r.indexOffset += 6
}

// Submit adds the given Renderable to this draw call.
func (r *Renderer) Submit(ra Renderable) {
	mesh := ra.mesh()
	x, y := ra.pos()

	var vertices []float32
	for i := 0; i < len(mesh.vertices); i += 3 {
		vertices[i] = mesh.vertices[i] + x
		vertices[i+1] = mesh.vertices[i+1] + y
		vertices[i+2] = mesh.vertices[i+2]
	}
	var indices []uint32
	for i := 0; i < len(mesh.indices); i++ {
		indices[i] = uint32(r.indexOffset) + mesh.indices[i]
	}

	if r.vertexOffset+len(vertices) >= r.vertexBufferSize || r.indexOffset+len(indices) >= r.indexBufferSize {
		r.Flush()
		r.indexOffset, r.vertexOffset = 0, 0
	}

	do(func() {
		gl.BufferSubData(gl.ARRAY_BUFFER, (r.vertexOffset+len(vertices))*float32Size, len(vertices)*float32Size, gl.Ptr(vertices))
		checkGLError()
		gl.BufferSubData(gl.ELEMENT_ARRAY_BUFFER, (r.indexOffset+len(indices))*int32Size, len(indices)*int32Size, gl.Ptr(indices))
		checkGLError()
	})
	r.vertexOffset += len(vertices)
	r.indexOffset += len(indices)
}

// Flush flushes all the draw calls that have been called on this renderer to the window
func (r *Renderer) Flush() {
	do(func() {
		gl.BindVertexArray(r.vertexArray)
		checkGLError()
		gl.DrawElements(gl.TRIANGLES, int32(r.indexOffset), gl.UNSIGNED_INT, nil)
		checkGLError()
		gl.BindVertexArray(0)
		checkGLError()
	})
}

// End ends the rendering procedure.
func (r *Renderer) End() {
	do(func() {
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
		checkGLError()
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		checkGLError()
	})
}
