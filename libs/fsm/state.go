package fsm

type Transition[Args any] struct {
	handler func(Args) *State[Args]
}

func (t *Transition[Args]) Handle(args Args) *State[Args] {
	return t.handler(args)
}

type State[Args any] struct {
	action      func()
	entryAction *func()
	exitAction  *func()
	transitions []*Transition[Args]
}

func (state *State[Args]) AddTransition(handler func(Args) *State[Args]) {
	state.transitions = append(state.transitions, &Transition[Args]{
		handler: handler,
	})
}

func (state *State[Args]) Handle(args Args) *State[Args] {
	for _, transitions := range state.transitions {
		newState := transitions.Handle(args)
		if newState != nil {
			return newState
		}
	}
	return state
}

func (state *State[Args]) Execute() {
	state.action()
}

func (state *State[Args]) Entry() {
	if state.entryAction != nil {
		(*state.entryAction)()
	}
}

func (state *State[Args]) Exit() {
	if state.exitAction != nil {
		(*state.exitAction)()
	}
}
