// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package game game.go Defines the struct used to connect all the engine components together.
package game

import (
	"runtime"
	"time"

	"github.com/bluemun/munfall"
	"github.com/bluemun/munfall/graphics"
	"github.com/bluemun/munfall/graphics/render"
	"github.com/bluemun/munfall/input"
	"github.com/bluemun/munfall/logic"
	"github.com/bluemun/munfall/worldmap"
	"github.com/go-gl/glfw/v3.2/glfw"
)

var mainHasRun = false

// Game type used to gold all the components needed to run a game.
type Game struct {
	actorRegistry  *logic.ActorRegistry
	Camera         *render.Camera
	orderGenerator input.OrderGenerator
	renderer       render.RendersTraits
	window         *graphics.Window
	world          munfall.World
	worldMap       worldmap.WorldMap
}

// Initialize initializes the game.
func (g *Game) Initialize(wm worldmap.WorldMap) {
	if !mainHasRun {
		mainHasRun = true
		go func() {
			runtime.LockOSThread()
			munfall.Loop()
		}()
	}

	g.actorRegistry = logic.CreateActorRegistry()
	g.window = graphics.CreateWindow()
	g.Camera = &render.Camera{}
	g.Camera.Activate()

	g.world = logic.CreateWorld()
	g.worldMap = wm

	// TODO: Change this once we got more renderers.
	g.renderer = render.CreateRendersTraits2D(g.world)
}

// Start starts the game loop.
func (g *Game) Start(framerate int64) {
	ticker := time.NewTicker(time.Second / (time.Duration)(framerate))

	for {
		select {
		case <-ticker.C:
			if g.window.Closed() {
				ticker.Stop()
				close(munfall.Mainfunc)
				return
			}

			g.world.Tick(1.0 / (float32)(framerate))
			g.window.PollEvents()
			if g.orderGenerator != nil {
				for _, order := range g.orderGenerator.GetOrders() {
					g.world.ResolveOrder(order)
				}
			}

			g.window.Clear()
			g.renderer.Render()
			g.window.SwapBuffers()
		}
	}
}

// SetOrderGenerator sets the current active order generator for the game.
func (g *Game) SetOrderGenerator(og input.OrderGenerator) {
	g.orderGenerator = og
	g.window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, code int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press || action == glfw.Release {
			g.orderGenerator.HandleKey(code, action == glfw.Press)
		}
	})
}

// ActorRegistry returns the inner actor registry for the game..
func (g *Game) ActorRegistry() *logic.ActorRegistry {
	return g.actorRegistry
}

// World returns the underlying world.
func (g *Game) World() munfall.World {
	return g.world
}
