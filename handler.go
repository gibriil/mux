// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mux

import "context"

// A Handler responds to a Routable.
type Handler interface {
	Process(context.Context, Routable, Handler) error
}

// The HandlerFunc type is an adapter to allow the use of ordinary functions as Routable handlers. If f is a function with the appropriate signature, HandlerFunc(f) is a Handler that calls f.
type HandlerFunc func(context.Context, Routable, Handler) error

// Process calls f(ctx, r, tx).
func (f HandlerFunc) Process(ctx context.Context, r Routable, tx Handler) error {
	return f(ctx, r, tx)
}
