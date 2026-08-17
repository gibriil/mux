// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mux

import "errors"

type Mux struct {
	routes map[Route]Handler
}

func (m *Mux) Handle(route Route, handler Handler) {
	m.register(route, handler)
}

func (m *Mux) HandleFunc(route Route, handler func(Routable) error) {
	m.register(route, HandlerFunc(handler))
}

func (m *Mux) register(route Route, handler Handler) {
	routeKey := resolveRoute(route)
	if err, bad := routeKey.(error); bad && errors.Is(err, ErrNoRoute) {
		panic(err)
	}
	if _, exists := m.routes[routeKey]; exists {
		panic("route already has a registered Handler")
	}
	m.routes[route] = handler
}

func (m *Mux) Process(r Routable) error {
	routeKey := resolveRoute(Route(r))
	if err, bad := routeKey.(error); bad && errors.Is(err, ErrNoRoute) {
		return err
	}
	h, exists := m.routes[routeKey]

	if !exists {
		return ErrNoHandler
	}

	return h.Process(r)
}
