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
	gameStarted bool

	deleteChannel chan struct{}
}

func newRoom(ID uint64, payload model.RoomPayload) *room {
	return &room{
		id:           ID,
		payload:      payload,
		events:       newRoomEventFeed(),
		counter:      0, // TODO: clarifiy if creating room still needs separate request for join / not
		gameStarted: false,

		deleteChannel: make(chan struct{}),
	}
}

// ID returns the room ID.
func (r *room) ID() uint64 {
	return r.id
}

// CreateMember creates a member representation in this room.
func (r *room) CreateMember(payload model.MemberPayload) (m model.Member) {
	m = &member{payload: payload}
	r.members.Store(payload.ID, m)
	atomic.AddInt32(&r.counter, 1)

	r.events.put(model.RoomEventMemberJoin(&model.RoomEventMemberJoinPayload{
		ID:   payload.ID,
		Name: payload.Name,
	}))
	return
}

// DeleteMember removes a member from this room by the member ID.
func (r *room) DeleteMember(memberID uint64) {
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

func (r *room) MaximumCapacity() uint {
	return r.payload.MaxMemberCount
}

// Size returns size of the member map
func (r *room) Size() int32 {
	return r.counter
}

func (r *room) IsFull() bool {
	return string(r.Size()) == string(r.MaximumCapacity())
}

func (r *room) StartGame() {
	r.gameStarted = true
	r.events.put(model.RoomEventGameBegin())
}

func (r *room) IsGameStarted() bool {
	return r.gameStarted
}
