package mux

import (
	"context"
	"errors"
	"log"
)

var (
	ErrNoHandler = errors.New("No handler registered for Route")
)

type Multiplexer struct {
	Source   Source
	Handler  Handler
	ErrorLog *log.Logger
}

func NewMux() *Mux {
	return &Mux{
		routes: map[Route]Handler{},
	}
}

type Source interface {
	Next(context.Context) (Routable, error)
}

func (m *Multiplexer) Process() error {
	return m.ProcessWithContext(context.Background())
}

func (m *Multiplexer) ProcessWithContext(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		route, err := m.Source.Next(ctx)
		if err != nil {
			m.logError(err)
			continue
		}

		if err := m.Handler.Process(route); err != nil {
			m.logError(err)
			continue
		}
	}
}

func (m *Multiplexer) logError(err error) {
	if m.ErrorLog != nil {
		m.ErrorLog.Println(err)
	} else {
		log.Printf("Route Handler produced error: %v", err)
	}
}
