// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mux

const maxDepth = 3

// Route is used for registering handlers.
//
// Route must satisfy Routable or be comparable
type Route any

// Routable is required for route resolution
//
// Route() must return a comparable before maxDepth
type Routable interface {
	Route() Route
}

// resolveRoute loops through Routable's Route() until a comparable is found that does not satisfy Routable or until maxDepth is reached.
func resolveRoute(r Route) Route {
	route := r
	if _, ok := route.(Routable); !ok {
		return route
	}

	for depth := 0; depth < maxDepth; depth++ {
		if next, ok := route.(Routable); ok {
			route = next.Route()
			continue
		}

		return route
	}

	return ErrNoRoute
}
