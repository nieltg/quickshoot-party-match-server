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
func (handler *Handler) Handler() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/room/new", handler.newRoom).Methods(http.MethodPost)
	router.HandleFunc("/room/{roomID}/events", handler.listRoomEvents).Methods(http.MethodGet)
	router.HandleFunc("/room/{roomID}/member/new", handler.newRoomMember).Methods(http.MethodPost)
	router.HandleFunc("/room/{roomID}/member/{memberID}", handler.deleteRoomMember).Methods(http.MethodDelete)

	router.HandleFunc("/room/{roomID}/member/{memberID}/tap", handler.registerTapTime).Methods(http.MethodPost)

	return router
}

func (handler *Handler) newRoom(writer http.ResponseWriter, request *http.Request) {
	var body newRoomRequest
	if decodeJSONBody(writer, request.Body, &body) != true {
		return
	}

	room := handler.Domain.CreateRoom(body.Payload)

	response := newRoomResponse{
		ID: room.ID(),
	}

	writeJSON(writer, response)
}

func (handler *Handler) listRoomEvents(writer http.ResponseWriter, request *http.Request) {
	room := handler.fetchRoom(writer, mux.Vars(request)["roomID"])
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
			writer.WriteHeader(http.StatusInternalServerError)
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
		case <-time.After(handler.DeferredRequestMaxDuration):
		}
	}

	writeJSON(writer, roomEventsResponse{
		Notifications: notifs,
		LastID:        nextLastID,
	})
}

func (handler *Handler) newRoomMember(writer http.ResponseWriter, request *http.Request) {
	room := handler.fetchRoom(writer, mux.Vars(request)["roomID"])
	if room == nil {
		return
	}

	var body newRoomMemberRequest
	if decodeJSONBody(writer, request.Body, &body) != true {
		return
	}

	var status bool
	var maximumCapacity uint64
	if status, maximumCapacity = room.CreateMember(body.Payload); status != true {
		writer.WriteHeader(http.StatusForbidden)
		return
	}

	writeJSON(writer, model.RoomPayload{
		MaxMemberCount: maximumCapacity,
	})
}

func (handler *Handler) deleteRoomMember(writer http.ResponseWriter, request *http.Request) {
	room := handler.fetchRoom(writer, mux.Vars(request)["roomID"])
	if room == nil {
		return
	}

	memberID, err := strconv.ParseUint(mux.Vars(request)["memberID"], 10, 64)
	if err != nil {
		return
	}

	if room.DeleteMember(memberID) != true {
		writer.WriteHeader(http.StatusForbidden)
		return
	}
}

func (handler *Handler) registerTapTime(writer http.ResponseWriter, request *http.Request) {
	room := handler.fetchRoom(writer, mux.Vars(request)["roomID"])
	if room == nil {
		return
	}

	memberID, err := strconv.ParseUint(mux.Vars(request)["memberID"], 10, 64)
	if err != nil {
		return
	}
	var body newTapTimeRequest
	if !decodeJSONBody(writer, request.Body, &body) {
		return
	}

	if room.RecordTapTime(memberID, body.Payload) != true {
		writer.WriteHeader(http.StatusForbidden)
		return
	}
}
