// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mux

// A Handler responds to an Routable.
type Handler interface {
	Process(Routable) error
}

// The HandlerFunc type is an adapter to allow the use of ordinary functions as Routable handlers. If f is a function with the appropriate signature, HandlerFunc(f) is a Handler that calls f.
type HandlerFunc func(Routable) error

// Process calls f(r).
func (f HandlerFunc) Process(r Routable) error {
	return f(r)
}
