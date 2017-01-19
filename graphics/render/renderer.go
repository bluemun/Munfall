// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package render renderer.go Defines a renderer in the graphics package.
package render

import (
	"github.com/bluemun/engine"
	"github.com/bluemun/engine/graphics/shader"
	"github.com/go-gl/gl/v3.3-core/gl"
)

// Renderable interface used to pass data to a renderer.
type Renderable interface {
	mesh() *Mesh
	pos() (float32, float32)
}

// Renderer interface used to talk to renderers.
type Renderer interface {
	Begin()
	DrawRectangle(x, y, w, h float32)
	Submit(ra Renderable)
	Flush()
	End()
}

// Renderer used by the graphics library to draw.
type renderer2d struct {
	vertexOffset, indexOffset, vertexBufferSize, indexBufferSize int
	vertexArray, vertexBuffer, indexBuffer                       uint32
	s                                                            shader.Shader
}

const int32Size = 4
const float32Size = 4
const vertexSize = 3 * float32Size

const vertexShader = `
#version 130
in highp vec3 vertex;
void main() {
    gl_Position = vec4(vertex, 1);
}
` + "\x00"

const fragmentShader = `
#version 130
out highp vec4 outputColor;
void main() {
    outputColor = vec4(1, 0, 1, 1);
}
` + "\x00"

// CreateRenderer2D used to create a renderer2d object correctly.
func CreateRenderer2D(vertexBufferSize, indexBufferSize int) Renderer {
	r := new(renderer2d)
	r.vertexBufferSize = vertexBufferSize
	r.indexBufferSize = indexBufferSize
	engine.Do(func() {
		r.s = shader.CreateShader(vertexShader, fragmentShader)
		r.s.Use()
		gl.GenVertexArrays(1, &r.vertexArray)
		engine.CheckGLError()
		gl.BindVertexArray(r.vertexArray)
		engine.CheckGLError()

		gl.GenBuffers(1, &r.vertexBuffer)
		engine.CheckGLError()
		gl.BindBuffer(gl.ARRAY_BUFFER, r.vertexBuffer)
		engine.CheckGLError()
		gl.BufferData(gl.ARRAY_BUFFER, (int)(10000*vertexSize), nil, gl.DYNAMIC_DRAW)
		engine.CheckGLError()

		vertAttrib := r.s.GetAttributeLocation("vertex")
		engine.Logger.Info("Vertex attribute location: ", vertAttrib)
		gl.EnableVertexAttribArray(vertAttrib)
		engine.CheckGLError()
		gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, vertexSize, gl.PtrOffset(0))
		engine.CheckGLError()

		gl.GenBuffers(1, &r.indexBuffer)
		engine.CheckGLError()
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, r.indexBuffer)
		engine.CheckGLError()

		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, (int)(10000*int32Size), nil, gl.DYNAMIC_DRAW)
		engine.CheckGLError()

		r.s.BindFragDataLocation("outputColor")
		engine.CheckGLError()
	})

	return r
}

// Begin starts the rendering procedure.
func (r *renderer2d) Begin() {
	r.indexOffset, r.vertexOffset = 0, 0
	engine.Do(func() {
		r.s.Use()
		gl.BindBuffer(gl.ARRAY_BUFFER, r.vertexBuffer)
		engine.CheckGLError()
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, r.indexBuffer)
		engine.CheckGLError()
	})
}

// DrawRectangle draws a rectangle using the given values, x and y point to the top-left corner.
func (r *renderer2d) DrawRectangle(x, y, w, h float32) {
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

	engine.Do(func() {
		gl.BufferSubData(gl.ARRAY_BUFFER, (r.vertexOffset)*float32Size, 12*float32Size, gl.Ptr(&vertices[0]))
		engine.CheckGLError()
		gl.BufferSubData(gl.ELEMENT_ARRAY_BUFFER, (r.indexOffset)*int32Size, 6*int32Size, gl.Ptr(&indices[0]))
		engine.CheckGLError()
	})
	r.vertexOffset += 12
	r.indexOffset += 6
}

// Submit adds the given Renderable to this draw call.
func (r *renderer2d) Submit(ra Renderable) {
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

	engine.Do(func() {
		gl.BufferSubData(gl.ARRAY_BUFFER, (r.vertexOffset+len(vertices))*float32Size, len(vertices)*float32Size, gl.Ptr(vertices))
		engine.CheckGLError()
		gl.BufferSubData(gl.ELEMENT_ARRAY_BUFFER, (r.indexOffset+len(indices))*int32Size, len(indices)*int32Size, gl.Ptr(indices))
		engine.CheckGLError()
	})
	r.vertexOffset += len(vertices)
	r.indexOffset += len(indices)
}

// Flush flushes all the draw calls that have been called on this renderer to the window
func (r *renderer2d) Flush() {
	engine.Do(func() {
		gl.BindVertexArray(r.vertexArray)
		engine.CheckGLError()
		gl.DrawElements(gl.TRIANGLES, int32(r.indexOffset), gl.UNSIGNED_INT, nil)
		engine.CheckGLError()
		gl.BindVertexArray(0)
		engine.CheckGLError()
	})
}

// End ends the rendering procedure.
func (r *renderer2d) End() {
	engine.Do(func() {
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
		engine.CheckGLError()
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		engine.CheckGLError()
	})
}
