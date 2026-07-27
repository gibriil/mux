package mux

import (
	"context"
	"log"
)

type Multiplexer struct {
	Source   Source
	Handler  Handler
	ErrorLog *log.Logger
}

func (m *Multiplexer) Process() error
func (m *Multiplexer) ProcessWithContext(ctx context.Context) error
