package http

import (
	"github.com/nieltg/quickshoot-party-match-server/pkg/domain"
)

type newRoomResponse struct {
	ID uint64
}

type roomNotificationsResponse struct {
	Notifications interface{}
	LastID        int
}

type joinRoomRequest struct {
	Member domain.Member
}
