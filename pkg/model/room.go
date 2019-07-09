package model

// Room represents game room which members can join and play together.
type Room interface {
	// ID returns the room ID.
	ID() uint64

	// CreateMember creates a member representation in this room.
	CreateMember(data MemberPayload) Member
	// DeleteMember removes a member from this room by the member ID.
	DeleteMember(memberID uint64)
	// Member finds a member based on the member ID or returns nil if not found.
	Member(memberID uint64) Member

	// Events returns feed of events happening in this room.
	Events() RoomEventFeed
}

// RoomPayload represents payload of a game room.
type RoomPayload struct {
	MaxMemberCount uint
}
