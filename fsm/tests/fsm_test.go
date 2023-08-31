package tests

import (
	"testing"

	"github.com/necroin/golibs/fsm"
)

func TestFsm(t *testing.T) {
	channel := make(chan string, 3)
	FSM := fsm.NewFSM[int]()
	fisrtState := FSM.AddState(func() {
		channel <- "State 1"
	})

	secondState := FSM.AddState(func() {
		channel <- "State 2"
	})

	fisrtState.AddTransition(func(value int) *fsm.State[int] {
		if value == 2 {
			return secondState
		}
		return nil
	})

	secondState.AddTransition(func(value int) *fsm.State[int] {
		if value == 1 {
			return fisrtState
		}
		return nil
	})

	FSM.Execute()
	FSM.Handle(2)
	FSM.Execute()
	FSM.Handle(1)
	FSM.Execute()

	message := <-channel
	if message != "State 1" {
		t.Errorf("[FSM] [Test] [Error] received message: %s", message)
	}

	message = <-channel
	if message != "State 2" {
		t.Errorf("[FSM] [Test] [Error] received message: %s", message)
	}

	message = <-channel
	if message != "State 1" {
		t.Errorf("[FSM] [Test] [Error] received message: %s", message)
	}
}
