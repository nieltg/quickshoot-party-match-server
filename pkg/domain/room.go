package domain

// Room is a representation of game room.
type Room struct {
	ID uint64

	deleteChannel chan struct{}
}
