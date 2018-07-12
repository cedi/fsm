package fsm

const (
	Finite = "Finite"
)

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

// The generic finite state which ends the state machine

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
