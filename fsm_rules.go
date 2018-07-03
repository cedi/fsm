package fsm

type FsmRules struct {
	transitions map[string]map[string]struct{} // because go doesn't know sets i use map[]struct{}
	whitelist   bool                           // should the ruleset act as an whitelist (true) or as an blacklist (false)
}

// Check if the Transition is allowed
func (r *FsmRules) IsTransitionAllowed(from State, to State) bool {
	set, ok := r.transitions[from.String()]
	if !ok {
		// if the entry is not in the rule set
		// in case of a whitelist we deny this state transition
		// in case of a blacklist we allow this state transition
		return !r.whitelist
	}

	_, ok = set[to.String()]
	if ok {
		// if the entry is in the rule set
		// in case of a whitelist we allow this state transition
		// in case of a blacklist we deny this state transition
		return r.whitelist
	}

	// If the entry exists as correct from state, but the to state is not there
	return !r.whitelist
}

// Add a transition to validate
func (r *FsmRules) AddTransition(from string, to string) {
	s := r.transitions[from]
	if s == nil {
		s = make(map[string]struct{})
	}
	s[to] = struct{}{}
	r.transitions[from] = s
}

// Create a new fsm ruleset
func newFsmRules(whitelist bool) *FsmRules {
	return &FsmRules{
		transitions: make(map[string]map[string]struct{}),
		whitelist:   whitelist,
	}
}

func NewFsmWhitelist() *FsmRules {
	return newFsmRules(true)
}

func NewFsmBlacklist() *FsmRules {
	return newFsmRules(false)
}
