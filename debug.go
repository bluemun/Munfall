// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package munfall debug.go Defines our logging.
package munfall

import (
	"os"
	"runtime/debug"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/op/go-logging"
)

// Logger used to log information.
var Logger = logging.MustGetLogger("graphics")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func init() {
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	logging.NewBackendFormatter(backend1, format)

	// Only errors and more severe messages should be sent to backend1
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.INFO, "")

	// Set the backends to be used.
	logging.SetBackend(backend1Leveled)
}

// CheckGLError used to check and log any errors that happen in opengl calls.
func CheckGLError() {
	if Logger.IsEnabledFor(logging.DEBUG) {
		if e := gl.GetError(); e != gl.NO_ERROR {
			debug.PrintStack()
			Logger.Critical("OpenGL error: ", e)
		}
	}
}
