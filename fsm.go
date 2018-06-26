package fsm

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

type Event string

type State interface {
	// Execute the work in this State.
	// To transition into the next state, return the new state and give a reason
	Run() (State, string)

	// MUST return true if "to" is the same state
	Compare(to State) bool

	// return the State-Name as String
	String() string
}

type FSM struct {
	Events chan Event
	errs   chan error
	state  State
	rules  *FsmRules
	mu     sync.RWMutex
}

func NewFSM(rules *FsmRules) *FSM {
	return &FSM{
		Events: make(chan Event),
		errs:   make(chan error),
		rules:  rules,
	}
}

func (fsm *FSM) SetIdleState(state State) {
	if fsm.state != nil {
		log.WithFields(log.Fields{
			"current fsm state":  fsm.state,
			"desired idle state": state,
		}).Error("Idle state not empty")
		return
	}

	fsm.state = state
}

func (fsm *FSM) Run() {
	next, reason := fsm.state.Run()

	for {
		if fsm.state.Compare(next) {
			continue
		}

		if fsm.rules != nil && !fsm.rules.IsTransitionAllowed(fsm.state, next) {
			log.WithFields(log.Fields{
				"last_state": fsm.state.String(),
				"new_state":  next.String(),
				"reason":     reason,
			}).Panic("FSM: State transition is not allowed")
		}

		log.WithFields(log.Fields{
			"last_state": fsm.state.String(),
			"new_state":  next.String(),
			"reason":     reason,
		}).Info("FSM: State change")

		fsm.mu.Lock()
		fsm.state = next
		fsm.mu.Unlock()

		next, reason = fsm.state.Run()
	}
}
