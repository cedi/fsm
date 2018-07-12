package fsm

// The Event structure which defines an event
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
