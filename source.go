// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mux

import "context"

// Source is the interface for supplying the Processor with input.
// Next() is called to receive the Routable for matching Route() to a Handler
type Source interface {
	Next(context.Context) (Routable, error)
}
