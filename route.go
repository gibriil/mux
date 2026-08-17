// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mux

const maxDepth = 3

type Route interface{}

type Routable interface {
	Route() Route
}

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
