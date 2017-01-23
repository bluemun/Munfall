// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package traits render.go Defines interfaces that define traits that are used for rendering.
package traits

import (
	"github.com/bluemun/engine"
)

// TraitRender2D called by TraitRenderManager to get renderables
// for rendering (2D implementation).
type TraitRender2D interface {
	Render2D() []engine.Renderable
}
