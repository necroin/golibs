package fsm

type FSM[Args any] struct {
	states       []*State[Args]
	currentState *State[Args]
}

func NewFSM[Args any]() *FSM[Args] {
	return &FSM[Args]{
		states:       []*State[Args]{},
		currentState: nil,
	}
}

func (fsm *FSM[Args]) AddState(action func()) *State[Args] {
	state := &State[Args]{
		action:      action,
		entryAction: nil,
		exitAction:  nil,
		transitions: []*Transition[Args]{},
	}

	fsm.states = append(fsm.states, state)

	if fsm.currentState == nil {
		fsm.currentState = state
	}

	return state
}

func (fsm *FSM[Args]) Handle(args Args) {
	if fsm.currentState != nil {
		fsm.currentState = fsm.currentState.Handle(args)
	}
}

func (fsm *FSM[Args]) Execute() {
	if fsm.currentState != nil {
		fsm.currentState.Execute()
	}
}

func (fsm *FSM[Args]) SetCurrentState(state *State[Args]) {
	fsm.currentState = state
}
