package model

type roomEventValue struct {
	Type roomEventType
}

type roomEventValueWithPayload struct {
	Type roomEventType

	Payload interface{}
}

// RoomEventMemberJoin builds a member-join event.
func RoomEventMemberJoin(payload *RoomEventMemberJoinPayload) RoomEvent {
	return roomEventValueWithPayload{
		Type:    roomEventTypeMemberJoin,
		Payload: payload,
	}
}

// RoomEventMemberLeave builds a member-leave event.
func RoomEventMemberLeave(payload *RoomEventMemberLeavePayload) RoomEvent {
	return roomEventValueWithPayload{
		Type:    roomEventTypeMemberLeave,
		Payload: payload,
	}
}

// RoomEventGameBegin builds a game-begin event.
func RoomEventGameBegin() RoomEvent {
	return roomEventValue{
		Type: roomEventTypeGameBegin,
	}
}

// RoomEventMemberTapTime builds a member-tap-time event.
func RoomEventMemberTapTime(payload *RoomEventMemberTapTimePayload) RoomEvent {
	return roomEventValueWithPayload{
		Type:    roomEventTypeMemberTapTime,
		Payload: payload,
	}
}

// RoomEventGameEnd builds a game-end event.
func RoomEventGameEnd(payload *RoomEventGameEndPayload) RoomEvent {
	return roomEventValueWithPayload{
		Type:    roomEventTypeGameEnd,
		Payload: payload,
	}
}
