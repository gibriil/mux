package mux

type Route string

type Routable interface {
	Route() Route
}

type Handler interface {
	Process(Routable) error
}

type HandlerFunc func(Routable) error

func (f HandlerFunc) Process(r Routable) error {
	return f(r)
}

type Mux struct {
	routes map[Route]Handler
}

func (m *Mux) Handle(route Route, handler Handler) {
	m.register(route, handler)
}

func (m *Mux) HandleFunc(route Route, handler func(r Routable) error) {
	m.register(route, HandlerFunc(handler))
}

func (m *Mux) register(route Route, handler Handler) {
	if _, ok := m.routes[route]; !ok {
		panic("Route already has a registered Handler")
	}
	m.routes[route] = handler
}

func (m *Mux) Process(r Routable) error {
	h := m.routes[r.Route()]
	return h.Process(r)
}
