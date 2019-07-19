package modelmemory

import (
	"sync"
	"math"
	"time"

	"github.com/nieltg/quickshoot-party-match-server/pkg/model"
)

type room struct {
	id          uint64
	payload     model.RoomPayload
	events      *roomEventFeed
	memberCount uint64

	members      sync.Map //key: userID, value: user object
	membersMutex sync.RWMutex

	tapTimes      map[uint64]uint64 //key: userID, value: tap time
	tapTimesMutex sync.RWMutex

	deleteChannel chan struct{}
}

func newRoom(ID uint64, payload model.RoomPayload) *room {
	return &room{
		id:       ID,
		payload:  payload,
		events:   newRoomEventFeed(),
		tapTimes: make(map[uint64]uint64),

		deleteChannel: make(chan struct{}),
	}
}

func waitForCompletion() {
	select {
	case <- time.After(3 * time.Second):
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
func (r *room) CreateMember(payload model.MemberPayload) bool {
	r.membersMutex.Lock()
	defer r.membersMutex.Unlock()

	if r.incrMemberCountIfAllowed() != true {
		return false
	}

	m := &member{payload: payload}
	r.members.Store(payload.ID, m)

	r.events.put(model.RoomEventMemberJoin(&model.RoomEventMemberJoinPayload{
		ID:   payload.ID,
		Name: payload.Name,
	}))

	if r.isRoomFull() {
		r.startGame()
	}

	return true
}

func (r *room) decrMemberCountIfAllowed() bool {
	// Full means game has been started.
	if r.isRoomFull() {
		return false
	}

	r.memberCount--
	return true
}

func (r *room) isEmpty() bool {
	return r.memberCount < 1;
}

// DeleteMember removes a member from this room by the member ID.
func (r *room) DeleteMember(memberID uint64) bool {
	r.membersMutex.Lock()
	defer r.membersMutex.Unlock()

	if r.decrMemberCountIfAllowed() != true {
		return false
	}

	r.members.Delete(memberID)
	r.events.put(model.RoomEventMemberLeave(&model.RoomEventMemberLeavePayload{
		MemberID: memberID,
	}))

	defer func() {
		if r.isEmpty() {
			waitForCompletion()
			close(r.deleteChannel)
		}
	}()

	return true
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

func (r *room) RecordTapTime(userID uint64, data model.MemberTapTimePayload) bool {
	r.tapTimesMutex.Lock()
	defer r.tapTimesMutex.Unlock()

	if r.isRoomFull() != true {
		return false
	}

	r.tapTimes[userID] = data.TimeInMilis

	r.events.put(model.RoomEventMemberTapTime(&model.RoomEventMemberTapTimePayload{
		MemberID: userID,
		TapTime:  data.TimeInMilis,
	}))

	winnerUserTime, winner := r.findWinner()
	if winner != nil {
		r.events.put(model.RoomEventGameEnd(&model.RoomEventGameEndPayload{
			BestTapTime: winnerUserTime,
			Winner:      winner.Payload(),
		}))

		waitForCompletion()
		close(r.deleteChannel)
	}

	return true
}

func (r *room) findWinner() (uint64, model.Member) {
	if uint64(len(r.tapTimes)) != r.memberCount {
		return math.MaxUint64, nil
	}

	var winnerUserTime uint64 = math.MaxUint64
	var winnerUserID uint64
	for userID, time := range r.tapTimes {
		if time < winnerUserTime {
			winnerUserTime = time
			winnerUserID = userID
		}
	}

	member := r.Member(winnerUserID)

	return winnerUserTime, member
}
