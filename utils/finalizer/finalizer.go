package finalizer

type Finalizer struct {
	handlers []func()
}

func NewFinalizer() *Finalizer {
	return &Finalizer{handlers: []func(){}}
}

func (finalizer *Finalizer) AddFunc(handler func()) {
	finalizer.handlers = append(finalizer.handlers, handler)
}

func (finalizer *Finalizer) Execute() {
	for _, handler := range finalizer.handlers {
		defer handler()
	}
}
