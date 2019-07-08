package http

type newRoomResponse struct {
	ID uint64
}

type roomNotificationsResponse struct {
	Notifications interface{}
	LastID        int
}
