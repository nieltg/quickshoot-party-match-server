package modelmemory

import (
	"sync"
	"sync/atomic"

	"github.com/nieltg/quickshoot-party-match-server/pkg/model"
)

type room struct {
	id      uint64
	payload model.RoomPayload
	events  *roomEventFeed
	members sync.Map

	counter      int32
	counterMutex sync.Mutex

	deleteChannel chan struct{}
}

func newRoom(ID uint64, payload model.RoomPayload) *room {
	return &room{
		id:      ID,
		payload: payload,
		events:  newRoomEventFeed(),
		counter: 0,

		deleteChannel: make(chan struct{}),
	}
}

// ID returns the room ID.
func (r *room) ID() uint64 {
	return r.id
}

func (r *room) incrementCounter() bool {
	r.counterMutex.Lock()
	defer r.counterMutex.Unlock()

	if r.isFull() {
		return false
	}

	atomic.AddInt32(&r.counter, 1)
	return true
}

// CreateMember creates a member representation in this room.
func (r *room) CreateMember(payload model.MemberPayload) (m model.Member) {
	if !r.incrementCounter() {
		return nil
	}

	m = &member{payload: payload}
	r.members.Store(payload.ID, m)

	r.events.put(model.RoomEventMemberJoin(&model.RoomEventMemberJoinPayload{
		ID:   payload.ID,
		Name: payload.Name,
	}))

	if r.isFull() {
		r.startGame()
	}

	return
}

func (r *room) decrementCounter() bool {
	r.counterMutex.Lock()
	defer r.counterMutex.Unlock()

	if r.isGameStarted() {
		return false
	}

	atomic.AddInt32(&r.counter, -1)
	return true
}

// DeleteMember removes a member from this room by the member ID.
func (r *room) DeleteMember(memberID uint64) {
	if !r.decrementCounter() {
		return
	}

	r.members.Delete(memberID)

	atomic.AddInt32(&r.counter, -1)

	r.events.put(model.RoomEventMemberLeave(&model.RoomEventMemberLeavePayload{
		MemberID: memberID,
	}))
}

// Member finds a member based on the member ID or returns nil if not found.
func (r *room) Member(memberID uint64) model.Member {
	value, ok := r.members.Load(memberID)
	if !ok {
		return nil
	}

	return value.(*member)
}

// Events returns feed of events happening in this room.
func (r *room) Events() model.RoomEventFeed {
	return r.events
}

func (r *room) maximumCapacity() uint {
	return r.payload.MaxMemberCount
}

// size returns size of the member map
func (r *room) size() int32 {
	return r.counter
}

func (r *room) isFull() bool {
	return string(r.size()) == string(r.maximumCapacity())
}

func (r *room) startGame() {
	r.events.put(model.RoomEventGameBegin())
}

func (r *room) isGameStarted() bool {
	return r.isFull()
}
