// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mux

import (
	"context"
	"errors"
	"sync"
)

// Mux is a general route to handler multiplexer.
// It matches the Route of each Routable from Source against a list of registered routes and calls the handler for the route that matches the Routable's Route.
type Mux struct {
	mu     sync.RWMutex
	routes map[Route]Handler
}

// NewMux creates a new Mux.
func NewMux() *Mux {
	return &Mux{
		routes: map[Route]Handler{},
	}
}

// Handle registers the handler for the given Route in Mux.
func (m *Mux) Handle(route Route, handler Handler) {
	m.register(route, handler)
}

// HandleFunc registers the handler function for the given Route in Mux.
func (m *Mux) HandleFunc(route Route, handler func(context.Context, Routable, Handler) error) {
	m.register(route, HandlerFunc(handler))
}

// register registers a handler for a route.
func (m *Mux) register(route Route, handler Handler) {
	routeKey, err := resolveRoute(route)
	if errors.Is(err, ErrNoRoute) {
		panic(err)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.routes[routeKey]; exists {
		panic("route already has a registered Handler")
	}
	m.routes[routeKey] = handler
}

// Process resolves the route from Routable and executes the matching Handler
func (m *Mux) Process(ctx context.Context, r Routable, tx Handler) error {
	routeKey, err := resolveRoute(Route(r))
	if errors.Is(err, ErrNoRoute) {
		return err
	}

	m.mu.RLock()
	h, exists := m.routes[routeKey]
	m.mu.RUnlock()

	if !exists {
		return ErrNoHandler
	}

	return h.Process(ctx, r, tx)
}
