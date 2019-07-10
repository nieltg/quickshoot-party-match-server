package modelmemory

import (
	"sync"

	"github.com/nieltg/quickshoot-party-match-server/pkg/model"
)

type room struct {
	id      uint64
	payload model.RoomPayload
	events  *roomEventFeed
	members sync.Map

	memberCount      uint64
	memberCountMutex sync.Mutex

	deleteChannel chan struct{}
}

func newRoom(ID uint64, payload model.RoomPayload) *room {
	return &room{
		id:      ID,
		payload: payload,
		events:  newRoomEventFeed(),

		deleteChannel: make(chan struct{}),
	}
}

// ID returns the room ID.
func (r *room) ID() uint64 {
	return r.id
}

func (r *room) isRoomFull() bool {
	return r.memberCount == r.payload.MaxMemberCount
}

func (r *room) incrMemberCountIfAllowed() bool {
	r.memberCountMutex.Lock()
	defer r.memberCountMutex.Unlock()

	if r.isRoomFull() {
		return false
	}

	r.memberCount++
	return true
}

func (r *room) startGame() {
	r.events.put(model.RoomEventGameBegin())
}

// CreateMember creates a member representation in this room.
func (r *room) CreateMember(payload model.MemberPayload) (m model.Member) {
	if !r.incrMemberCountIfAllowed() {
		return nil
	}

	m = &member{payload: payload}
	r.members.Store(payload.ID, m)

	r.events.put(model.RoomEventMemberJoin(&model.RoomEventMemberJoinPayload{
		ID:   payload.ID,
		Name: payload.Name,
	}))

	if r.isRoomFull() {
		r.startGame()
	}

	return
}

func (r *room) decrMemberCountIfAllowed() bool {
	r.memberCountMutex.Lock()
	defer r.memberCountMutex.Unlock()

	// Full means game has been started.
	if r.isRoomFull() {
		return false
	}

	r.memberCount--
	return true
}

// DeleteMember removes a member from this room by the member ID.
func (r *room) DeleteMember(memberID uint64) {
	if !r.decrMemberCountIfAllowed() {
		return
	}

	r.members.Delete(memberID)

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
