package modelmemory

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/nieltg/quickshoot-party-match-server/pkg/model"
)

// Domain represents game domain.
type Domain struct {
	JoinMaxDuration time.Duration

	rooms      sync.Map
	roomNextID uint64
}

// CreateRoom creates a room in this game domain.
func (domain *Domain) CreateRoom(payload model.RoomPayload) model.Room {
	room := newRoom(atomic.AddUint64(&domain.roomNextID, 1), payload)

	domain.rooms.Store(room.id, room)
	go domain.autoDeleteRoom(room)

	return room
}

func (domain *Domain) autoDeleteRoom(room *room) {
	select {
	case <-room.deleteChannel:
	case <-time.After(domain.JoinMaxDuration):
	}

	domain.rooms.Delete(room.id)
}

// Room finds a room by the room ID or return nil if room is missing.
func (domain *Domain) Room(ID uint64) model.Room {
	value, ok := domain.rooms.Load(ID)
	if !ok {
		return nil
	}

	return value.(*room)
}
