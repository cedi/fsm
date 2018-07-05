package fsm

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventCompare(t *testing.T) {
	ev1 := Event{Name: "change", Data: "random_payload1"}
	ev2 := Event{Name: "change", Data: "random_payload2"}
	ev3 := Event{Name: "anotherevnt", Data: "random_payload3"}

	assert.True(t, ev1.Compare(ev2))
	assert.False(t, ev2.Compare(ev3))
}

func TestFiniteState(t *testing.T) {
	finite := NewFiniteState()
	finite2 := NewFiniteState()

	assert.True(t, finite.String() == Finite)
	assert.True(t, finite.Compare(finite2))

	newState, evnt, reason := finite.Run(nil, Event{})
	assert.True(t, newState.Compare(finite))
	assert.True(t, reason == Finite)
	assert.True(t, evnt.Compare(Event{}))
}

// idle state type
type fooState struct {
	fsm            *FSM
	state          string
	newStateReason string
}

func FooState() string {
	return "Idle"
}

func (s fooState) Run(lastS State, lastE Event) (State, Event, string) {
	for {
		event := <-s.fsm.InEvents

		switch event.Name {
		case "connect":
			return newFooState(s.fsm), event, "Client connected"

		case "disconnect":
			return NewFiniteState(), event, "Client disconnected"

		default:
			continue
		}
	}
}

func (s fooState) String() string {
	return s.state
}

func (s fooState) Compare(to State) bool {
	return s.String() == to.String()
}

// Create a new Idle-State
func newFooState(fsm *FSM) *fooState {
	return &fooState{
		fsm:   fsm,
		state: FooState(),
	}
}

func TestStateTransition(t *testing.T) {
	f := NewNonLoggingFSM(nil)
	f.SetIdleState(newFooState(f))

	go func(f *FSM) {
		f.SendEvent("connect", nil)
		f.SendEvent("disconnect", nil)
	}(f)

	f.Run()

	fmt.Printf("Current state: %v", f.State())
	assert.True(t, f.State().String() == Finite)
}
