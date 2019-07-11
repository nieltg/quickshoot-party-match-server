package model

// Member represents a member of game room.
type Member interface {
	RegisterTapTime(tapTimeMilis uint64)
}

// MemberPayload represents payload of a room member.
type MemberPayload struct {
	ID   uint64
	Name string
}

// MemberTapTimePayload represents tap time in ms
type MemberTapTimePayload struct {
	TimeInMilis uint64
}