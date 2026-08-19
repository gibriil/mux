// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mux_test

import (
	"context"
	"errors"
	"io"
	"log"
	"testing"

	"github.com/gibriil/mux"
)

type blockingSource struct{}

func (blockingSource) Next(ctx context.Context) (mux.Routable, error) {
	<-ctx.Done()
	return nil, ctx.Err()
}

func TestProcessorContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	p := &mux.Processor{
		Source:   blockingSource{},
		Handler:  mux.HandlerFunc(func(mux.Routable) error { return nil }),
		ErrorLog: log.New(io.Discard, "", 0),
	}

	err := p.ProcessWithContext(ctx)

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("ProcessWithContext() error = %v, want context.Canceled", err)
	}
}
