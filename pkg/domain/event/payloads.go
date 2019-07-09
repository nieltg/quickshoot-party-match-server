package event

// MemberJoinPayload represents payload for member-join event.
type MemberJoinPayload struct {
	ID   uint64
	Name string
}

// MemberLeavePayload represents payload for member-leave event.
type MemberLeavePayload struct {
	MemberID uint64
}

// MemberTapTimePayload represents payload for member-tap-time event.
type MemberTapTimePayload struct {
	MemberID uint64
}

// GameEndPayload represents payload for game-end event.
type GameEndPayload struct {
	BestTapTime uint
}
