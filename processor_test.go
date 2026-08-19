// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mux_test

import (
	"errors"
	"io"
	"log"
	"testing"

	"github.com/gibriil/mux"
)

func TestProcessorProcess(t *testing.T) {
	events := []mux.Routable{
		event{Name: "one"},
		event{Name: "two"},
		event{Name: "three"},
	}

	source := &eventBus{
		events: events,
	}

	var received []mux.Routable

	handler := mux.HandlerFunc(func(r mux.Routable) error {
		received = append(received, r)
		return nil
	})

	p := &mux.Processor{
		Source:   source,
		Handler:  handler,
		ErrorLog: log.New(io.Discard, "", 0),
	}

	err := p.Process()

	if !errors.Is(err, io.EOF) {
		t.Fatalf("Process() error = %v, want io.EOF", err)
	}

	if len(received) != len(events) {
		t.Fatalf("received %d events, want %d", len(received), len(events))
	}

	for i := range events {
		if received[i].Route() != events[i].Route() {
			t.Errorf("received[%d] = %v, want %v", i, received[i], events[i])
		}
	}
}

func TestProcessorContinuesAfterHandlerError(t *testing.T) {
	events := []mux.Routable{
		event{Name: "one"},
		event{Name: "two"},
		event{Name: "three"},
	}

	source := &eventBus{events: events}

	var received []mux.Routable

	handler := mux.HandlerFunc(func(r mux.Routable) error {
		received = append(received, r)

		if r.Route() == "two" {
			return errors.New("handler failed")
		}

		return nil
	})

	p := &mux.Processor{
		Source:   source,
		Handler:  handler,
		ErrorLog: log.New(io.Discard, "", 0),
	}

	err := p.Process()

	if !errors.Is(err, io.EOF) {
		t.Fatalf("Process() error = %v, want io.EOF", err)
	}

	if len(received) != 3 {
		t.Fatalf("received %d events, want 3", len(received))
	}
}
