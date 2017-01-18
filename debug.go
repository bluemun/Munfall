// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package graphics debug.go Defines our logging.
package graphics

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/op/go-logging"
	"os"
	"runtime/debug"
)

var logger = logging.MustGetLogger("graphics")
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

func checkGLError() {
	if logger.IsEnabledFor(logging.DEBUG) {
		if e := gl.GetError(); e != gl.NO_ERROR {
			debug.PrintStack()
			logger.Critical("OpenGL error: ", e)
		}
	}
}
