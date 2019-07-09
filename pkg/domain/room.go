package domain

import (
	"fmt"
	"sync"

	"github.com/nieltg/quickshoot-party-match-server/pkg/util"
)

// Member is representation of game user who joined room.
type Member struct {
	ID   uint64
	Name string
}

// Room is a representation of game room.
type Room struct {
	ID      uint64
	Events  *util.Feed
	Members sync.Map

	deleteChannel chan struct{}
}

func newRoom(ID uint64) *Room {
	return &Room{
		ID:     ID,
		Events: util.NewFeed(),

		deleteChannel: make(chan struct{}),
	}
}

// Join is a function to let users join
func (room *Room) Join(member *Member) {
	room.Members.Store(member.ID, member)

	room.Events.Put(fmt.Sprintf("User %d has joined!", member.ID))
}
