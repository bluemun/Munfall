// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package graphics Defines a Window type used to create and control
// a glfw window with an opengl context.
package graphics

import (
	"fmt"

	"github.com/bluemun/munfall"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

// Window stores all the values needed to create and work with a window.
type Window struct {
	inner *glfw.Window
}

// CreateWindow creates a window correctly.
func CreateWindow() *Window {
	window := &Window{}

	munfall.Do(func() {
		var err error
		if err = glfw.Init(); err != nil {
			munfall.Logger.Panic("Failed to initialize GLFW.")
		}

		window.inner, err = glfw.CreateWindow(800, 600, "Test", nil, nil)
		if err != nil {
			munfall.Logger.Panic("Failed to create GLFW window: ", err)
		}

		window.inner.MakeContextCurrent()
		if err = gl.Init(); err != nil {
			munfall.Logger.Panic("Failed to initialize OpenGL: ", err)
		}
		version := gl.GoStr(gl.GetString(gl.VERSION))
		fmt.Println("OpenGL version", version)

		gl.ClearColor(0, 0, 0, 1)
	})

	return window
}

// SetKeyCallback sets the current key callback on the window.
func (w *Window) SetKeyCallback(f func(w *glfw.Window, key glfw.Key, code int, action glfw.Action, mods glfw.ModifierKey)) {
	w.inner.SetKeyCallback(f)
}

// Clear clears the screen.
func (w *Window) Clear() {
	munfall.Do(func() {
		gl.Clear(gl.COLOR_BUFFER_BIT)
	})
}

// Closed returns if this window has been closed.
func (w *Window) Closed() bool {
	return w.inner.ShouldClose()
}

// PollEvents polls the window events, should be called at the end of every update iteration.
func (w *Window) PollEvents() {
	glfw.PollEvents()
}

// SwapBuffers swaps the window buffers, should be called at the end of every render iteration.
func (w *Window) SwapBuffers() {
	munfall.Do(func() {
		w.inner.SwapBuffers()
	})
}
