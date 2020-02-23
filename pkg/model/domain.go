package model

// Domain represents game domain.
type Domain interface {
    // CreateRoom creates a room in this game domain.
    CreateRoom(data RoomPayload) Room
    // Room finds a room by the room ID or return nil if room is missing.
    Room(roomID uint64) Room
}
