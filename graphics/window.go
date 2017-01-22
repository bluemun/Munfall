// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package graphics Defines a Window type used to create and control
// a glfw window with an opengl context.
package graphics

import (
	"fmt"

	"github.com/bluemun/engine"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

// Window stores all the values needed to create and work with a window.
type Window struct {
	inner *glfw.Window
}

// CreateWindow creates a window correctly.
func CreateWindow() *Window {
	var window *Window
	engine.Do(func() {
		var err error
		if err = glfw.Init(); err != nil {
			engine.Logger.Panic("Failed to initialize GLFW.")
		}

		window = &Window{}
		window.inner, err = glfw.CreateWindow(800, 600, "Test", nil, nil)
		if err != nil {
			engine.Logger.Panic("Failed to create GLFW window: ", err)
		}

		window.inner.MakeContextCurrent()
		if err = gl.Init(); err != nil {
			engine.Logger.Panic("Failed to initialize OpenGL: ", err)
		}

		version := gl.GoStr(gl.GetString(gl.VERSION))
		fmt.Println("OpenGL version", version)

		gl.ClearColor(0, 0, 0, 1)
	})

	return window
}

// Clear clears the screen.
func (w *Window) Clear() {
	engine.Do(func() {
		gl.Clear(gl.COLOR_BUFFER_BIT)
	})
}

// Closed returns if this window has been closed.
func (w *Window) Closed() bool {
	return w.inner.ShouldClose()
}

// IsKeyPressed checks if the given key code is pressed.
func (w *Window) IsKeyPressed(code int) bool {
	return false
}

// SwapBuffers swaps the window buffers, should be called at the end of every render iteration.
func (w *Window) SwapBuffers() {
	engine.Do(func() {
		w.inner.SwapBuffers()
	})
}

// PollEvents polls the window events, should be called at the end of every update iteration.
func (w *Window) PollEvents() {
	glfw.PollEvents()
}
