// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package input ordergenerator.go Defines how we handle our input.
package input

import (
	"github.com/bluemun/engine"
)

// OrderGenerator issues orders based on input.
type OrderGenerator interface {
	GetOrders() []*engine.Order
	HandleKey(code int, state bool)
	HandleMouseMove(x, y float32)
	HandleMouseButton(button int)
}

// ScriptableOrderGenerator an order generator implementation that allows you to
// add hotkeys as orders that will be generated.
type ScriptableOrderGenerator struct {
	keyScripts map[int]map[bool]*engine.Order
	orders     []*engine.Order
}

// CreateScriptableOrderGenerator creates a ScriptableOrderGenerator and initializes it.
func CreateScriptableOrderGenerator() *ScriptableOrderGenerator {
	return &ScriptableOrderGenerator{keyScripts: make(map[int]map[bool]*engine.Order, 0)}
}

// GetOrders bla.
func (s *ScriptableOrderGenerator) GetOrders() []*engine.Order {
	x := s.orders
	s.orders = nil
	return x
}

// HandleKey handles key presses.
func (s *ScriptableOrderGenerator) HandleKey(code int, state bool) {
	codeMap, exists := s.keyScripts[code]
	if exists {
		value, exists := codeMap[state]
		if exists {
			newValue := &engine.Order{Order: value.Order, Value: value.Value}
			s.orders = append(s.orders, newValue)
		}
	}
}

// HandleMouseMove handles mouse movement.
func (s *ScriptableOrderGenerator) HandleMouseMove(x, y float32) {
}

// HandleMouseButton handles mouse buttons.
func (s *ScriptableOrderGenerator) HandleMouseButton(button int) {
}

// AddKeyScript adds a key as an order that this generator sends.
func (s *ScriptableOrderGenerator) AddKeyScript(code int, down bool, order string, value interface{}) {
	codeMap, exists := s.keyScripts[code]
	if !exists {
		codeMap = make(map[bool]*engine.Order, 0)
		s.keyScripts[code] = codeMap
	}

	codeMap[down] = &engine.Order{Order: order, Value: value}
}
