// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mux_test

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/gibriil/mux"
)

type event struct {
	Name    string
	Details map[string]any
}

func (e event) Route() mux.Route {
	return e.Name
}

type eventBus struct {
	events []mux.Routable
	cursor int
}

func (e *eventBus) Next(ctx context.Context) (mux.Routable, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if e.cursor >= len(e.events) {
		return nil, io.EOF
	}

	event := e.events[e.cursor]
	e.cursor++

	return event, nil
}

func TestMuxProcess(t *testing.T) {
	m := mux.NewMux()

	var received mux.Routable

	m.Handle("test", mux.HandlerFunc(func(r mux.Routable) error {
		received = r
		return nil
	}))

	e := event{
		Name: "test",
		Details: map[string]any{
			"msg": "hello",
		},
	}

	if err := m.Process(e); err != nil {
		t.Fatalf("Process() error = %v", err)
	}

	if received.(event).Name != e.Name {
		t.Fatalf("handler received %v, want %v", received, e)
	}
}

func TestMuxProcessNoHandler(t *testing.T) {
	m := mux.NewMux()

	err := m.Process(event{
		Name: "missing",
	})

	if !errors.Is(err, mux.ErrNoHandler) {
		t.Fatalf("Process() error = %v, want ErrNoHandler", err)
	}
}

func TestHandlerFunc(t *testing.T) {
	called := false

	handler := mux.HandlerFunc(func(r mux.Routable) error {
		called = true
		return nil
	})

	if err := handler.Process(event{}); err != nil {
		t.Fatalf("Process() error = %v", err)
	}

	if !called {
		t.Fatal("handler was not called")
	}
}

func TestMuxHandleDuplicate(t *testing.T) {
	m := mux.NewMux()

	handler := mux.HandlerFunc(func(mux.Routable) error {
		return nil
	})

	m.Handle("test", handler)

	defer func() {
		if recover() == nil {
			t.Fatal("expected duplicate route registration to panic")
		}
	}()

	m.Handle("test", handler)
}
