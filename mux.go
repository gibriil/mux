// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mux

import (
	"errors"
	"sync"
)

// Mux is a general route to handler multiplexer.
// It matches the Route of each Routable from Source against a list of registered routes and calls the handler for the route that matches the Routable's Route.
type Mux struct {
	mu     sync.RWMutex
	routes map[Route]Handler
}

// Handle registers the handler for the given Route in Mux.
func (m *Mux) Handle(route Route, handler Handler) {
	m.register(route, handler)
}

// HandleFunc registers the handler function for the given Route in Mux.
func (m *Mux) HandleFunc(route Route, handler func(Routable) error) {
	m.register(route, HandlerFunc(handler))
}

// register registers a handler for a route.
func (m *Mux) register(route Route, handler Handler) {
	routeKey := resolveRoute(route)
	if err, bad := routeKey.(error); bad && errors.Is(err, ErrNoRoute) {
		panic(err)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.routes[routeKey]; exists {
		panic("route already has a registered Handler")
	}
	m.routes[route] = handler
}

// Process resolves the route from Routable and executes the matching Handler
func (m *Mux) Process(r Routable) error {
	routeKey := resolveRoute(Route(r))
	if err, bad := routeKey.(error); bad && errors.Is(err, ErrNoRoute) {
		return err
	}

	m.mu.RLock()
	h, exists := m.routes[routeKey]
	m.mu.RUnlock()

	if !exists {
		return ErrNoHandler
	}

	return h.Process(r)
}
