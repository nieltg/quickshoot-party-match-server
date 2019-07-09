package model

// Member represents a mmeber of game room.
type Member interface {
	RegisterTapTime(tapTimeMilis uint64)
}

// MemberPayload represents payload of a room member.
type MemberPayload struct {
	ID   uint64
	Name string
}
