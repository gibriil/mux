package mux

type mux struct {
	routes map[string]Handler
}

func (m *mux) Handle(route string, handler Handler)
func (m *mux) HandleFunc(route string, handler HandlerFunc)
func (m *mux) Process(r Routable) error

type Routable interface {
	Route() string
}

type Handler interface {
	Process(Routable) error
}

type HandlerFunc func(Routable) error

type Source interface {
	Next() (Routable, error)
}
