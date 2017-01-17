package graphics

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Window struct {
	inner *glfw.Window
	r     *renderer
}

func CreateWindow() *Window {
	glfw.Init()
	window := new(Window)
	window.inner, _ = glfw.CreateWindow(800, 600, "Test", nil, nil)
	window.inner.MakeContextCurrent()
	gl.Init()
	window.r = CreateRenderer()

	return window
}

func (w *Window) GetRenderer() *renderer {
	return w.r
}

func (w *Window) SwapBuffers() {
	w.inner.SwapBuffers()
}
