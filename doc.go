// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package mux provides a general purpose multiplexer for executing a Handler on any comparable.

A route simple needs to satisfy Routable. The Route() must return a comparable at some point. That comparable is used as the Route key for registering the Handler.

Current maxDepth is 3. Route() must return a comparable with less than maxDepth calls to Route().

# Simple String Router

Start by setting up a new mux
`mux := mux.NewMux()`

Register handlers for each route
```go

	mux.Handlefunc("Route 1", func(r Routable) error {
		.......
	})

	mux.Handlefunc("Route 2", func(r Routable) error {
		.......
	})

	mux.Handlefunc("Route 3", func(r Routable) error {
		.......
	})

```

Create a Routable
```go

	type Event struct {
		Name string
		Detail string
	}

	func (e Event) Route() mux.Route {
		return e.Name
	}

```

Send a Routable through the processor
```go

	err := mux.Process(Event{
		Name: "Rout 2",
		Detail: "This is my test event"
	})

```

# Full Processor

Like the net/http Server, mux provides a Processor. The processor can take a Source to provide continuing input for processing.

Start by setting up your Routable and your Source
```go

	type Event struct {
		Name    string
		Details map[string]any
	}

	func (e Event) Route() mux.Route {
		return e.Name
	}

	type EventBus struct {
		events []mux.Routable
		cursor int
	}

	func (e *EventBus) Next(ctx context.Context) (mux.Routable, error) {
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

```

Set up your mux and attach Handlers to routes
```go

	mux := mux.NewMux()

	mux.Handlefunc("Route 1", func(r Routable) error {
		.......
	})

	mux.Handlefunc("Route 2", func(r Routable) error {
		.......
	})

	mux.Handlefunc("Route 3", func(r Routable) error {
		.......
	})

```

Setup your processes and start the event loop
```go

	p := &mux.Processor{
		Source:   &eventBus{
			events: []mux.Routable{
				Event{Name: "Route 1"},
				Event{Name: "Route 2"},
				Event{Name: "Route 3"},
			},
		},
		Handler:  mux,
		ErrorLog: log.New(io.Discard, "", 0),
	}

	err := p.Process()

```

Once the source has stopped receiving input, expect err to receive `io.EOF`
*/
package mux
