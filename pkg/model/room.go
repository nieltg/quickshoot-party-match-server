package model

// Room represents game room which members can join and play together.
type Room interface {
    // ID returns the room ID.
    ID() uint64

    // CreateMember creates a member representation in this room.
    CreateMember(data MemberPayload) (bool, uint64)
    // DeleteMember removes a member from this room by the member ID.
    DeleteMember(memberID uint64) bool
    // Member finds a member based on the member ID or returns nil if not found.
    Member(memberID uint64) Member

    // Events returns feed of events happening in this room.
    Events() RoomEventFeed

    // RecordTapTime records user tap time and return 200 if succeed
    RecordTapTime(userID uint64, data MemberTapTimePayload) bool
}

// RoomPayload represents payload of a game room.
type RoomPayload struct {
    MaxMemberCount uint64
}
