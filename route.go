// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mux

import "context"

const maxDepth = 3

// Route is a value used to identify a registered route.
//
// A Route may be a Routable, in which case it is resolved recursively,
// or it may be a terminal comparable value.
type Route any

// Routable is required for route resolution
//
// Route() must return a comparable before maxDepth
type Routable interface {
	Route() Route
}

// routableWithContext is an internal wrapper for a Routable to provide context to Handlers
type routableWithContext struct {
	Routable
	ctx context.Context
}

// resolveRoute loops through Routable's Route() until a comparable is found that does not satisfy Routable or until maxDepth is reached.
func resolveRoute(r Route) (Route, error) {
	for depth := 0; depth < maxDepth; depth++ {
		next, ok := r.(Routable)
		if !ok {
			return r, nil
		}

		r = next.Route()
	}

	if _, ok := r.(Routable); ok {
		return nil, ErrNoRoute
	}

	return r, nil
}
