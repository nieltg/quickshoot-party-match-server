package event

// Event represents an event.
type Event interface{}

// Value represents an event value.
type value struct {
	Type Type
}

// ValueWithPayload represents an event value with payload.
type valueWithPayload struct {
	Type    Type
	Payload interface{}
}
