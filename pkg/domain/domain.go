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
	room := newRoom(atomic.AddUint64(&domain.roomNextID, 1))

	domain.rooms.Store(room.ID, room)
	go domain.autoDeleteRoom(room)

	return room
}

// Room returns room based on specified room ID.
func (domain *Domain) Room(ID uint64) *Room {
	value, ok := domain.rooms.Load(ID)
	if !ok {
		return nil
	}

	return value.(*Room)
}
