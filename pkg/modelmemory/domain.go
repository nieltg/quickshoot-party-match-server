package modelmemory

import (
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
	"math"

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

func duplicateMap(source, destination *sync.Map) {
	source.Range(func(key, value interface{}) bool {
		id, _ := strconv.ParseUint(fmt.Sprintf("%d", key), 10, 64)
		destination.Store(id, value)
		return true
	})
}

func autoTap(room *room) {
	var membersClone sync.Map
	duplicateMap(&(*room).members, &membersClone)
	for id := range room.tapTimes {
		membersClone.Delete(id)
	}

	membersClone.Range(func(key, value interface{}) bool {
		id, _ := strconv.ParseUint(fmt.Sprintf("%d", key), 10, 64)
		recordStatus := room.RecordTapTime(id, model.MemberTapTimePayload{
			TimeInMilis: uint64(5.0 * math.Round(time.Minute.Seconds() * 1000.0)),
		})

		return recordStatus
	})
}

func (domain *Domain) autoDeleteRoom(room *room) {
	select {
	case <-room.deleteChannel:
	case <-time.After(domain.JoinMaxDuration):
		gameBegins := len(room.tapTimes) > 0
		if gameBegins {
			autoTap(room)
		}
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
