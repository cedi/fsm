package fsm

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
)

const (
	Finite = "Finite"
)

type Event struct {
	Name string
	Data interface{}
}

type State interface {
	// Execute the work in this State.
	// To transition into the next state, return the new state and give a reason
	// last: the previous state. For the initial state: empty
	Run(last State) (State, string)

	// MUST return true if "to" is the same state
	Compare(to State) bool

	// return the State-Name as String
	String() string
}

type FiniteState struct {
	fsm *FSM
}

func (s FiniteState) Run(last State) (State, string) {
	return s, Finite
}

func (s FiniteState) Compare(to State) bool {
	return s.String() == to.String()
}

func (s FiniteState) String() string {
	return Finite
}

func NewFiniteState(fsm *FSM) *FiniteState {
	return &FiniteState{fsm: fsm}
}

type FSM struct {
	Events chan Event
	errs   chan error
	state  State
	rules  *FsmRules
	log    bool
	mu     sync.RWMutex
}

func NewLoggingFSM(rules *FsmRules) *FSM {
	return NewFSM(rules, true)
}

func NewNonLoggingFSM(rules *FsmRules) *FSM {
	return NewFSM(rules, false)
}

func NewFSM(rules *FsmRules, log bool) *FSM {
	return &FSM{
		Events: make(chan Event),
		errs:   make(chan error),
		rules:  rules,
		log:    log,
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

func (fsm *FSM) Run() {
	var last State
	finiteState := NewFiniteState(fsm)

	for {
		next, reason := fsm.state.Run(last)
		if next.Compare(finiteState) {
			if fsm.log {
				log.WithFields(log.Fields{
					"current_state": fsm.state.String(),
					"new_state":     next.String(),
					"reason":        reason,
				}).Info("Finite")
			}

			break
		}

		if fsm.rules != nil && !fsm.rules.IsTransitionAllowed(fsm.state, next) {
			if fsm.log {
				log.WithFields(log.Fields{
					"last_state": fsm.state.String(),
					"new_state":  next.String(),
					"reason":     reason,
				}).Error("State transition is not allowed")
			}
		}

		if fsm.log {
			log.WithFields(log.Fields{
				"current_state": fsm.state.String(),
				"new_state":     next.String(),
				"reason":        reason,
			}).Info("State change")
		}

		last = fsm.state

		fsm.mu.Lock()
		fsm.state = next
		fsm.mu.Unlock()
	}
}
