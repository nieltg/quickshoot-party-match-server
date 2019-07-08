package domain

// Room is a representation of game room.
type Room struct {
	ID   uint64
	Feed *Feed

	deleteChannel chan struct{}
}

func newRoom(ID uint64) *Room {
	return &Room{
		ID:   ID,
		Feed: NewFeed(),

		deleteChannel: make(chan struct{}),
	}
}
