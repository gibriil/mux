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
	Handler  Handler
	ErrorLog *log.Logger
}

// New
func NewMux() *Mux {
	return &Mux{
		routes: map[Route]Handler{},
	}
}

func (m *Processor) Process() error {
	return m.ProcessWithContext(context.Background())
}

func (m *Processor) ProcessWithContext(ctx context.Context) error {
	for {
		route, err := m.Source.Next(ctx)
		if err != nil {
			return err
		}

		if err := m.Handler.Process(route); err != nil {
			m.logError(err)
			continue
		}
	}
}

func (m *Processor) logError(err error) {
	if m.ErrorLog != nil {
		m.ErrorLog.Println(err)
	} else {
		log.Printf("Route Handler produced error: %v", err)
	}
}
