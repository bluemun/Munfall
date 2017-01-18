// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package graphics Defines a Window type used to create and control
// a glfw window with an opengl context.
package graphics

import (
	"fmt"
	"runtime"
	"time"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/op/go-logging"
)

func init() {
	runtime.LockOSThread()
}

// Window stores all the values needed to create and work with a window.
type Window struct {
	inner *glfw.Window
	r     *Renderer
	s     *shader
}

// CreateWindow creates a window correctly.
func CreateWindow() *Window {
	var window *Window
	do(func() {
		var err error
		if err = glfw.Init(); err != nil {
			logger.Panic("Failed to initialize GLFW.")
		}

		window = new(Window)
		window.inner, err = glfw.CreateWindow(800, 600, "Test", nil, nil)
		if err != nil {
			logger.Panic("Failed to create GLFW window: ", err)
		}

		window.inner.MakeContextCurrent()
		if err = gl.Init(); err != nil {
			logger.Panic("Failed to initialize OpenGL: ", err)
		}

		version := gl.GoStr(gl.GetString(gl.VERSION))
		fmt.Println("OpenGL version", version)

		gl.ClearColor(0, 1, 1, 1)

		window.s = createShader(vertexShader, fragmentShader)
	})

	window.r = CreateRenderer(window.s)

	return window
}

// Loop used by opengl to do its calls, needs to be called by the main thread.
func Loop() {
	for f := range mainfunc {
		f()
	}
}

var mainfunc = make(chan func())

func do(f func()) {
	done := make(chan bool, 1)
	mainfunc <- func() {
		if logger.IsEnabledFor(logging.DEBUG) {
			timer := time.NewTicker(time.Second * 10)
			defer timer.Stop()
			go func() {
				<-timer.C
				logger.Critical("Main thread took more then 10 seconds to run a single function.")
			}()
		}

		f()
		done <- true
	}
	<-done
}

// Clear clears the screen.
func (w *Window) Clear() {
	do(func() {
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

// GetRenderer gets the default renderer for this window.
func (w *Window) GetRenderer() *Renderer {
	return w.r
}

// SwapBuffers swaps the window buffers, should be called at the end of every render iteration.
func (w *Window) SwapBuffers() {
	do(func() {
		w.inner.SwapBuffers()
	})
}

// PollEvents polls the window events, should be called at the end of every update iteration.
func (w *Window) PollEvents() {
	glfw.PollEvents()
}

var vertexShader = `
#version 130
in highp vec3 vertex;
void main() {
    gl_Position = vec4(vertex, 1);
}
` + "\x00"

var fragmentShader = `
#version 130
out highp vec4 outputColor;
void main() {
    outputColor = vec4(1, 0, 1, 1);
}
` + "\x00"
