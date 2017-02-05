// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package traits worldmap.go Defines interfaces that define traits
// that are used for worldmap interaction.
package traits

import (
	"math"

	"github.com/bluemun/engine"
)

// Space defines a space that can be used to
type Space interface {
	Offset() *engine.WPos
	Intersects(other Space) bool
}

// SpaceRect defines a Space that is a rectangle.
type SpaceRect struct {
	LocalOffset *engine.WPos
	HalfWidth   float32
	HalfHeight  float32
}

// Offset returns the centered offset of this space relative to the actor position.
func (sr *SpaceRect) Offset() *engine.WPos {
	return sr.LocalOffset
}

// Intersects returns if the two spaces intersect.
func (sr *SpaceRect) Intersects(other Space) bool {
	srpos := ((Space)(sr)).(engine.Trait).Owner().Pos()
	otherpos := other.(engine.Trait).Owner().Pos()
	sroff := sr.Offset()
	otheroff := other.Offset()

	srpos.X += sroff.X
	srpos.Y += sroff.Y
	srpos.Z += sroff.Z

	otherpos.X += otheroff.X
	otherpos.Y += otheroff.Y
	otherpos.Z += otheroff.Z

	// Checks intersection if other is a SpaceRect.
	osr, defined := other.(*SpaceRect)
	if defined {
		diff := srpos.Vector(otherpos)
		if math.Abs(float64(diff.X)) < float64(sr.HalfWidth+osr.HalfWidth) ||
			math.Abs(float64(diff.Y)) < float64(sr.HalfHeight+osr.HalfHeight) {
			return true
		}
	}

	return false
}

// OccupySpace defines a trait that occupies space and should collide with
// other actors that have an OccupySpace trait.
type OccupySpace interface {
	Space() []*Space
}
