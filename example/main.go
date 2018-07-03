package main

import (
	"time"

	"github.com/cedi/fsm"
	log "github.com/sirupsen/logrus"
)

// idle state type
type idleState struct {
	fsm            *fsm.FSM
	state          string
	newStateReason string
}

func IdleState() string {
	return "Idle"
}

// Do work
func (s idleState) Run(lastS fsm.State, lastE fsm.Event) (fsm.State, fsm.Event, string) {
	for {
		event := <-s.fsm.InEvents

		switch event.Name {
		case "connect":
			return newConnectedState(s.fsm), event, "Client connected"

		default:
			log.WithField("event", event).Warn("unknown event")
			continue
		}
	}
}

func (s idleState) String() string {
	return s.state
}

func (s idleState) Compare(to fsm.State) bool {
	return s.String() == to.String()
}

// Create a new Idle-State
func newIdleState(fsm *fsm.FSM) *idleState {
	return &idleState{
		fsm:   fsm,
		state: IdleState(),
	}
}

// connect state type
type connectedState struct {
	fsm   *fsm.FSM
	state string
}

func ConnectedState() string {
	return "Connected"
}

// Do work
func (s connectedState) Run(lastS fsm.State, lastE fsm.Event) (fsm.State, fsm.Event, string) {
	for {
		event := <-s.fsm.InEvents

		switch event.Name {
		case "disconnect":
			return fsm.NewFiniteState(), event, "Client disconnected"

		default:
			log.WithField("event", event).Warn("unknown event")
			continue
		}
	}
}

func (s connectedState) String() string {
	return s.state
}

func (s connectedState) Compare(to fsm.State) bool {
	return s.String() == to.String()
}

// Create a new Idle-State
func newConnectedState(fsm *fsm.FSM) *connectedState {
	return &connectedState{
		fsm:   fsm,
		state: ConnectedState(),
	}
}

func main() {
	rules := fsm.NewFsmWhitelist()
	rules.AddTransition(IdleState(), ConnectedState())
	rules.AddTransition(ConnectedState(), IdleState())

	f := fsm.NewLoggingFSM(rules)
	f.SetIdleState(newIdleState(f))

	go func(f *fsm.FSM) {
		for {
			f.InEvents <- fsm.Event{Name: "connect", Data: nil}
			f.InEvents <- fsm.Event{Name: "foo", Data: nil}
			time.Sleep(2 * time.Second)

			f.InEvents <- fsm.Event{Name: "bar", Data: nil}
			time.Sleep(2 * time.Second)

			f.InEvents <- fsm.Event{Name: "disconnect", Data: nil}
		}
	}(f)

	f.Run()
}
