package iohttp

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/nieltg/quickshoot-party-match-server/pkg/model"
)

// Handler represents a HTTP Handler.
type Handler struct {
	DeferredRequestMaxDuration time.Duration

	Domain model.Domain
}

// Handler returns new handler for HTTP requests.
func (s *Handler) Handler() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/room/new", s.newRoom)
	router.HandleFunc("/room/{id}/events", s.listRoomEvents)
	router.HandleFunc("/room/{id}/member/new", s.newRoomMember)

	return router
}

func (s *Handler) newRoom(w http.ResponseWriter, req *http.Request) {
	var body newRoomRequest
	if !decodeJSONBody(w, req.Body, &body) {
		return
	}

	room := s.Domain.CreateRoom(body.Payload)

	writeJSON(w, newRoomResponse{
		ID: room.ID(),
	})
}

func (s *Handler) listRoomEvents(w http.ResponseWriter, req *http.Request) {
	room := s.fetchRoom(w, mux.Vars(req)["id"])
	if room == nil {
		return
	}

	var lastID int

	lastIDStr := req.URL.Query().Get("lastID")
	if lastIDStr == "" {
		lastID = -1
	} else {
		var err error

		if lastID, err = strconv.Atoi(lastIDStr); err != nil {
			log.Println("Unable to parse last feed ID:", err)
			w.WriteHeader(500)
			return
		}
	}

	var notifs []model.RoomEvent
	var nextLastID int
	var waitChannel <-chan struct{}

	notifs, nextLastID, waitChannel = room.Events().List(lastID)
	if len(notifs) == 0 {
		select {
		case <-waitChannel:
			notifs, nextLastID, _ = room.Events().List(lastID)
		case <-time.After(s.DeferredRequestMaxDuration):
		}
	}

	writeJSON(w, roomEventsResponse{
		Notifications: notifs,
		LastID:        nextLastID,
	})
}

func (s *Handler) newRoomMember(w http.ResponseWriter, req *http.Request) {
	room := s.fetchRoom(w, mux.Vars(req)["id"])
	if room == nil {
		return
	}

	var body newRoomMemberRequest
	if !decodeJSONBody(w, req.Body, &body) {
		return
	}

	room.CreateMember(body.Payload)
}
