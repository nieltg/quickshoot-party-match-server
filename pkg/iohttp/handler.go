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
	router.HandleFunc("/room/new", s.newRoom).Methods("POST")
	router.HandleFunc("/room/{id}/events", s.listRoomEvents).Methods("GET")
	router.HandleFunc("/room/{id}/member/new", s.newRoomMember).Methods("POST")

	return router
}

func (s *Handler) newRoom(writer http.ResponseWriter, request *http.Request) {
	var body newRoomRequest
	if !decodeJSONBody(writer, request.Body, &body) {
		return
	}

	room := s.Domain.CreateRoom(body.Payload)

	response := newRoomResponse{
		ID:                       room.ID(),
		MaximumDurationInSeconds: s.Domain.MaximumJoinDurationInSeconds(),
	}

	writeJSON(writer, response)
}

func (s *Handler) listRoomEvents(writer http.ResponseWriter, request *http.Request) {
	room := s.fetchRoom(writer, mux.Vars(request)["id"])
	if room == nil {
		return
	}

	var lastID int

	lastIDStr := request.URL.Query().Get("lastID")
	if lastIDStr == "" {
		lastID = -1
	} else {
		var err error

		if lastID, err = strconv.Atoi(lastIDStr); err != nil {
			log.Println("Unable to parse last feed ID:", err)
			writer.WriteHeader(500)
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

	writeJSON(writer, roomEventsResponse{
		Notifications: notifs,
		LastID:        nextLastID,
	})
}

func (s *Handler) newRoomMember(writer http.ResponseWriter, request *http.Request) {
	room := s.fetchRoom(writer, mux.Vars(request)["id"])
	if room == nil {
		return
	}

	var body newRoomMemberRequest
	if !decodeJSONBody(writer, request.Body, &body) {
		return
	}

	room.CreateMember(body.Payload)

	if room.IsFull() {
		// TODO: Publish GAME_START event
		room.StartGame()

		// TODO: lock users from joining
	}
}
