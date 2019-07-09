package event

// MemberJoin builds a member-join event.
func MemberJoin(payload *MemberJoinPayload) Event {
	return valueWithPayload{
		Type:    TypeMemberJoin,
		Payload: payload,
	}
}

// MemberLeave builds a member-leave event.
func MemberLeave(payload *MemberLeavePayload) Event {
	return valueWithPayload{
		Type:    TypeMemberLeave,
		Payload: payload,
	}
}

// GameBegin builds a game-begin event.
func GameBegin() Event {
	return value{
		Type: TypeGameBegin,
	}
}

// MemberTapTime builds a member-tap-time event.
func MemberTapTime(payload *MemberTapTimePayload) Event {
	return valueWithPayload{
		Type:    TypeMemberTapTime,
		Payload: payload,
	}
}

// GameEnd builds a game-end event.
func GameEnd(payload *GameEndPayload) Event {
	return valueWithPayload{
		Type:    TypeGameEnd,
		Payload: payload,
	}
}
