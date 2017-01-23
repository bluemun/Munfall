// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package input ordergenerator.go Defines how we handle our input.
package input

import (
	"github.com/bluemun/engine"
	"github.com/bluemun/engine/logic"
)

// OrderGenerator issues orders based on input.
type OrderGenerator interface {
	IssueOrder(order, value string)
	HandleKey(code int)
	HandleMouseMove(x, y float32)
	HandleMouseButton(button int)
}

// ScriptableOrderGenerator bla.
type ScriptableOrderGenerator struct {
	world      *logic.World
	keyScripts map[int]*engine.Order
}

// IssueOrder bla.
func (s *ScriptableOrderGenerator) IssueOrder(order, value string) {

}

// HandleKey bla.
func (s *ScriptableOrderGenerator) HandleKey(code int) {
	value, exists := s.keyScripts[code]
	if exists {
		s.world.ResolveOrder(value)
	}
}

// HandleMouseMove bla.
func (s *ScriptableOrderGenerator) HandleMouseMove(x, y float32) {
}

// HandleMouseButton bla.
func (s *ScriptableOrderGenerator) HandleMouseButton(button int) {
}

// AddKeyScript bla.
func (s *ScriptableOrderGenerator) AddKeyScript(code int, order, value string) {
	s.keyScripts[code] = &engine.Order{Order: order, Value: value}
}
