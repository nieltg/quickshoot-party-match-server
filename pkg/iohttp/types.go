package iohttp

import "github.com/nieltg/quickshoot-party-match-server/pkg/model"

type newRoomRequest struct {
	Payload model.RoomPayload
}

type newRoomResponse struct {
	ID uint64
}

type roomEventsResponse struct {
	Notifications interface{}
	LastID        int
}

type newRoomMemberRequest struct {
	Payload model.MemberPayload
}
