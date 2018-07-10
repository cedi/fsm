package fsm

import (
	"fmt"
	"sync"

	log "github.com/cedi/logrus"
)

const (
	Finite = "Finite"
)

// event
type Event struct {
	// Name of the event
	Name string

	// Payload
	Data interface{}

	// Output Chan
	OutEvents chan Event
}

func (e Event) Compare(to Event) bool {
	return e.Name == to.Name
}

// Generic State Interface
type State interface {
	// Execute the work in this State.
	// To transition into the next state, return the new state and give a reason
	// lastS: the previous state. For the initial state: empty
	// lastE: the previous event. For the initial state: empty
	Run(lastS State, lastE Event) (State, Event, string)

	// MUST return true if "to" is the same state
	Compare(to State) bool

	// return the State-Name as String
	String() string
}

type FiniteState struct {
}

func (s FiniteState) Run(lastS State, lastE Event) (State, Event, string) {
	return s, Event{}, Finite
}

func (s FiniteState) Compare(to State) bool {
	return s.String() == to.String()
}

func (s FiniteState) String() string {
	return Finite
}

func NewFiniteState() *FiniteState {
	return &FiniteState{}
}

type FSM struct {
	errs  chan error
	state State
	rules *FsmRules
	mu    sync.RWMutex

	// event-Handling
	InEvents  chan Event
	OutEvents chan Event

	// Logging
	Logger *log.Entry
}

func NewLoggingFSM(rules *FsmRules, logger *log.Entry) *FSM {
	return NewFSM(rules, logger)
}

func NewNonLoggingFSM(rules *FsmRules) *FSM {
	return NewFSM(rules, nil)
}

func NewFSM(rules *FsmRules, logger *log.Entry) *FSM {
	return &FSM{
		InEvents:  make(chan Event),
		OutEvents: make(chan Event),
		errs:      make(chan error),
		rules:     rules,
		Logger:    logger,
	}
}

func (fsm *FSM) SetIdleState(state State) error {
	if fsm.state != nil {
		return fmt.Errorf("Idle status not empty. Current fsm state: %v, desired idle state: %v", fsm.state, state)
	}

	fsm.state = state
	return nil
}

func (fsm *FSM) State() State {
	return fsm.state
}

// Send a event to the FSM
//
//	name: the name of the event
//	data: the payload of the event
//
//	return: the output chanel for returning events from the fsm
func (fsm *FSM) SendEventToFsm(name string, data interface{}) chan Event {
	if fsm.Logger != nil {
		fsm.Logger.WithField("event", name).Trace("Send event to FSM")
	}

	outevents := make(chan Event)

	event := Event{
		Name:      name,
		Data:      data,
		OutEvents: outevents,
	}

	fsm.InEvents <- event

	return outevents
}

// Send a event from the FSM to the outside world
//
//	name: the name of the event
//	data: the payload of the event
//
//	return: the output chanel for returning events to the fsm
func (fsm *FSM) SendEventFromFsm(name string, data interface{}) chan Event {
	if fsm.Logger != nil {
		fsm.Logger.WithField("event", name).Trace("Send event from FSM")
	}

	outevents := make(chan Event)

	event := Event{
		Name:      name,
		Data:      data,
		OutEvents: outevents,
	}

	fsm.OutEvents <- event

	return outevents
}

// Should be called in a seperate go-routine
func (fsm *FSM) Run() {
	var last State
	var evnt Event
	finiteState := NewFiniteState()

	for {
		next, evnt, reason := fsm.state.Run(last, evnt)
		if next.Compare(finiteState) {
			if fsm.Logger != nil {
				fsm.Logger.WithFields(log.Fields{
					"current_state": fsm.state.String(),
					"new_state":     next.String(),
					"event":         evnt.Name,
					"reason":        reason,
				}).Info("Finite")
			}

			fsm.mu.Lock()
			fsm.state = next
			fsm.mu.Unlock()

			// This break ends the finite state machine...
			break
		}

		if fsm.rules != nil && !fsm.rules.IsTransitionAllowed(fsm.state, next) {
			if fsm.Logger != nil {
				fsm.Logger.WithFields(log.Fields{
					"last_state": fsm.state.String(),
					"new_state":  next.String(),
					"event":      evnt.Name,
					"reason":     reason,
				}).Error("State transition is not allowed")
			}
		}

		if fsm.Logger != nil {
			fsm.Logger.WithFields(log.Fields{
				"current_state": fsm.state.String(),
				"new_state":     next.String(),
				"event":         evnt.Name,
				"reason":        reason,
			}).Info("State change")
		}

		last = fsm.state

		fsm.mu.Lock()
		fsm.state = next
		fsm.mu.Unlock()
	}
}
