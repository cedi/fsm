package fsm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockState struct {
	name string
}

func (s MockState) Run(lastS State, lastE Event) (State, Event, string) {
	return lastS, lastE, ""
}

func (s MockState) Compare(to State) bool {
	return s.String() == to.String()
}

func (s MockState) String() string {
	return s.name
}

func TestWhitelistRuleset(t *testing.T) {
	whitelist := NewFsmWhitelist()
	assert.NotNil(t, whitelist)

	a := MockState{name: "A"}
	b := MockState{name: "B"}

	whitelist.AddTransition(a.String(), b.String())
	assert.True(t, whitelist.IsTransitionAllowed(a, b))
	assert.False(t, whitelist.IsTransitionAllowed(b, a))

}

func TestEmptyWhiteRuleset(t *testing.T) {
	whitelist := NewFsmWhitelist()
	assert.NotNil(t, whitelist)

	a := MockState{name: "A"}
	b := MockState{name: "B"}

	assert.False(t, whitelist.IsTransitionAllowed(a, b))
	assert.False(t, whitelist.IsTransitionAllowed(b, a))
}

func TestBlacklistRuleset(t *testing.T) {
	blacklist := NewFsmBlacklist()
	assert.NotNil(t, blacklist)

	a := MockState{name: "A"}
	b := MockState{name: "B"}

	blacklist.AddTransition(a.String(), b.String())
	assert.False(t, blacklist.IsTransitionAllowed(a, b))
	assert.True(t, blacklist.IsTransitionAllowed(b, a))
}

func TestEmptyBlacklistRuleset(t *testing.T) {
	blacklist := NewFsmBlacklist()
	assert.NotNil(t, blacklist)

	a := MockState{name: "A"}
	b := MockState{name: "B"}

	assert.True(t, blacklist.IsTransitionAllowed(a, b))
	assert.True(t, blacklist.IsTransitionAllowed(b, a))
}
