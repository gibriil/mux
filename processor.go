// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mux

import (
	"context"
	"errors"
	"log"
)

var (
	ErrNoHandler = errors.New("no handler registered for Route")
	ErrNoRoute   = errors.New("no valid comparable Route")
)

// A Processor defines parameters for running a multiplexer process
type Processor struct {
	Source   Source
	Handler  Handler // Handles Routable from Source
	Emitter  Handler // Processable Handler from within a Handler
	ErrorLog *log.Logger
}

// Process reads from Source to begin handling routes
func (m *Processor) Process() error {
	return m.ProcessWithContext(context.Background())
}

// Process reads from Source to begin handling routes with a user defined context
func (m *Processor) ProcessWithContext(ctx context.Context) error {
	for {
		route, err := m.Source.Next(ctx)
		if err != nil {
			return err
		}

		if err := m.Handler.Process(ctx, route, m.Emitter); err != nil {
			m.logError(err)
			continue
		}
	}
}

// logError logs errors to the processor's ErrorLog
func (m *Processor) logError(err error) {
	if m.ErrorLog != nil {
		m.ErrorLog.Println(err)
	} else {
		log.Printf("Route Handler produced error: %v", err)
	}
}
