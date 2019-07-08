package domain

import (
	"sync"
	"sync/atomic"
	"time"
)

// Domain is a representation of a game domain.
type Domain struct {
	JoinMaxDuration time.Duration

	rooms      sync.Map
	roomNextID uint64
}

func (domain *Domain) autoDeleteRoom(room *Room) {
	select {
	case <-room.deleteChannel:
	case <-time.After(domain.JoinMaxDuration):
	}

	domain.rooms.Delete(room.ID)
}

// CreateRoom creates a new room in current domain.
func (domain *Domain) CreateRoom() *Room {
	room := &Room{
		ID: atomic.AddUint64(&domain.roomNextID, 1),

		deleteChannel: make(chan struct{}),
	}

	domain.rooms.Store(room.ID, room)
	go domain.autoDeleteRoom(room)

	return room
}
