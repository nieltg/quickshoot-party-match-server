package event

// Value represents an event value.
type Value struct {
	Type Type
}

// ValueWithPayload represents an event value with payload.
type ValueWithPayload struct {
	Type    Type
	Payload interface{}
}
