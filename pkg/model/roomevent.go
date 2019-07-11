package model

// RoomEventFeed represents feed of room events.
type RoomEventFeed interface {
	List(lastID int) ([]RoomEvent, int, <-chan struct{})
}

// RoomEvent represents an event of the room.
type RoomEvent interface{}

type roomEventType string

const (
	roomEventTypeMemberJoin    roomEventType = "MEMBER_JOIN"
	roomEventTypeMemberLeave   roomEventType = "MEMBER_LEAVE"
	roomEventTypeGameBegin     roomEventType = "GAME_BEGIN"
	roomEventTypeMemberTapTime roomEventType = "MEMBER_TAP_TIME"
	roomEventTypeGameEnd       roomEventType = "GAME_END"
)

// RoomEventMemberJoinPayload represents payload for member-join event.
type RoomEventMemberJoinPayload struct {
	ID   uint64
	Name string
}

// RoomEventMemberLeavePayload represents payload for member-leave event.
type RoomEventMemberLeavePayload struct {
	MemberID uint64
}

// RoomEventMemberTapTimePayload represents payload for member-tap-time event.
type RoomEventMemberTapTimePayload struct {
	MemberID uint64
	TapTime uint64
}

// RoomEventGameEndPayload represents payload for game-end event.
type RoomEventGameEndPayload struct {
	BestTapTime uint
}
