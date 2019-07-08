package domain

import "github.com/nieltg/quickshoot-party-match-server/pkg/util"

// Room is a representation of game room.
type Room struct {
	ID   uint64
	Feed *util.Feed

	deleteChannel chan struct{}
}

func newRoom(ID uint64) *Room {
	return &Room{
		ID:   ID,
		Feed: util.NewFeed(),

		deleteChannel: make(chan struct{}),
	}
}
