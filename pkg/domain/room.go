package domain

import (
	"sync"

	"github.com/nieltg/quickshoot-party-match-server/pkg/domain/event"
	"github.com/nieltg/quickshoot-party-match-server/pkg/util"
)

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

// CreateMember is a function to let users join
func (room *Room) CreateMember(member *Member) {
	room.Members.Store(member.ID, member)

	room.Events.Put(event.MemberJoin(&event.MemberJoinPayload{
		ID:   member.ID,
		Name: member.Name,
	}))
}
