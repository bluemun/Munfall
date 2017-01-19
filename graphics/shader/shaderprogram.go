// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package shader shaderprogram.go Defines a shader type that
// will be used by opengl.
package shader

import (
	"fmt"
	"strings"

	"github.com/bluemun/engine"
	"github.com/go-gl/gl/v3.3-core/gl"
)

// Shader used to interact with an opengl shader.
type Shader interface {
	Use()
	GetAttributeLocation(name string) uint32
	BindFragDataLocation(loc string)
}

type shader uint32

// CreateShader used to create a shader from the given source.
func CreateShader(vertexShaderSource, fragmentShaderSource string) Shader {
	var program uint32

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		engine.Logger.Critical("Failed to compile vertex shader: ", err)
		return nil
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		engine.Logger.Critical("Failed to compile fragment shader: ", err)
		return nil
	}

	program = gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		engine.Logger.Critical("Failed to link shader program: ", err)
		return nil
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	var s = (shader)(program)
	return &s
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func (s *shader) Use() {
	gl.UseProgram((uint32)(*s))
	engine.CheckGLError()
}

func (s *shader) GetAttributeLocation(name string) uint32 {
	defer engine.CheckGLError()
	return uint32(gl.GetAttribLocation((uint32)(*s), gl.Str(name+"\x00")))
}

func (s *shader) BindFragDataLocation(loc string) {
	gl.BindFragDataLocation((uint32)(*s), 0, gl.Str(loc+"\x00"))
}
