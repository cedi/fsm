package fsm

import (
	"testing"
)

// idle state type
type idleState struct {
	fsm            *FSM
	state          string
	newStateReason string
}

func IdleState() string {
	return "Idle"
}

// Do work
func (s idleState) Run(last State) (State, string) {
	for {
		event := <-s.fsm.Events

		switch event.Name {
		case "connect":
			return newConnectedState(s.fsm), "Client connected"

		default:
			continue
		}
	}
}

func (s idleState) String() string {
	return s.state
}

func (s idleState) Compare(to State) bool {
	return s.String() == to.String()
}

// Create a new Idle-State
func newIdleState(fsm *FSM) *idleState {
	return &idleState{
		fsm:   fsm,
		state: IdleState(),
	}
}

// connect state type
type connectedState struct {
	fsm   *FSM
	state string
}

func ConnectedState() string {
	return "Connected"
}

// Do work
func (s connectedState) Run(last State) (State, string) {
	for {
		event := <-s.fsm.Events

		switch event.Name {
		case "disconnect":
			return newIdleState(s.fsm), "Client disconnected"

		default:
			continue
		}
	}
}

func (s connectedState) String() string {
	return s.state
}

func (s connectedState) Compare(to State) bool {
	return s.String() == to.String()
}

// Create a new Idle-State
func newConnectedState(fsm *FSM) *connectedState {
	return &connectedState{
		fsm:   fsm,
		state: ConnectedState(),
	}
}

func BenchmarkStateChanges(b *testing.B) {
	f := NewNonLoggingFSM(nil)
	f.SetIdleState(newIdleState(f))

	go f.Run()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Events <- Event{Name: "connect"}
		f.Events <- Event{Name: "disconnect"}
	}
}
