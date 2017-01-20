// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package render traits.go Defines traits that can be used to render objects.
package render

import (
	"reflect"

	"github.com/bluemun/engine/logic"
)

// TraitRender2D called by TraitRenderManager to get renderables
// for rendering (2D implementation).
type TraitRender2D interface {
	logic.Trait
	Render2D() []Renderable
}

// RendersTraits defines a collection of objects that can be used in conjunction
// with a world object to render traits.
type RendersTraits interface {
	Render()
}

type renderTraits2d struct {
	world    *logic.World
	renderer Renderer
}

// CreateRendersTraits2D creates a 2D implementation of RendersTraits.
func CreateRendersTraits2D(w *logic.World) RendersTraits {
	rm := new(renderTraits2d)
	rm.world = w
	rm.renderer = CreateRenderer2D(10000, 10000)
	return rm
}

func (r *renderTraits2d) Render() {
	var t1 = reflect.TypeOf((*TraitRender2D)(nil)).Elem()
	traits := r.world.Traitmanager.Lookup(t1)
	r.renderer.Begin()
	for _, trait := range traits {
		for _, renderable := range trait.(TraitRender2D).Render2D() {
			r.renderer.Submit(renderable)
		}
	}
	r.renderer.Flush()
	r.renderer.End()
}
