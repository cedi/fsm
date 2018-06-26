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
func (s idleState) Run() (fsm.State, string) {
	for {
		event := <-s.fsm.Events

		switch event {
		case "connect":
			return newConnectedState(s.fsm), "Client connected"

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
func (s connectedState) Run() (fsm.State, string) {
	for {
		event := <-s.fsm.Events

		switch event {
		case "disconnect":
			return newIdleState(s.fsm), "Client disconnected"

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

	f := fsm.NewFSM(rules)
	f.SetIdleState(newIdleState(f))

	go func(fsm *fsm.FSM) {
		for {
			fsm.Events <- "connect"
			fsm.Events <- "foo"
			time.Sleep(2 * time.Second)
			fsm.Events <- "disconnect"
			fsm.Events <- "bar"
			time.Sleep(2 * time.Second)
		}
	}(f)

	f.Run()
}
