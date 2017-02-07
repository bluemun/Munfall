// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package traits worldmap.go Defines interfaces that define traits
// that are used for worldmap interaction.
package traits

import (
	"github.com/bluemun/munfall"
)

// SpaceCell defines a Space that is a full cell.
type SpaceCell struct {
	trait munfall.Trait

	LocalOffset *munfall.WPos
}

// Initialize called when this space is created.
func (s *SpaceCell) Initialize(trait munfall.Trait) {
	s.trait = trait
}

// Trait returns the trait the owns this Space
func (s *SpaceCell) Trait() munfall.Trait {
	return s.trait
}

// Offset returns the centered offset of this space relative to the actor position.
func (s *SpaceCell) Offset() *munfall.WPos {
	return s.LocalOffset.Add(s.trait.Owner().Pos())
}

// Intersects returns if the two spaces intersect.
func (s *SpaceCell) Intersects(other munfall.Space, offset *munfall.WPos) bool {
	if offset == nil {
		offset = &munfall.WPos{}
	}

	mpos1 := s.trait.Owner().World().WorldMap().ConvertToMPos(s.Offset().Add(offset))
	mpos2 := s.trait.Owner().World().WorldMap().ConvertToMPos(other.Offset())
	return *mpos1 == *mpos2
}

// OccupySpace defines a trait that occupies space and should collide with
// other actors that have an OccupySpace trait.
type OccupySpace interface {
	munfall.Trait
	Intersects(OccupySpace, *munfall.WPos) bool
	Space() []munfall.Space
	OutOfBounds(*munfall.WPos) bool
}
