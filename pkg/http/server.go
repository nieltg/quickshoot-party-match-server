package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/nieltg/quickshoot-party-match-server/pkg/domain"
	"github.com/nieltg/quickshoot-party-match-server/pkg/util"
)

// Server represents a HTTP Server.
type Server struct {
	DeferredRequestMaxDuration time.Duration

	Domain domain.Domain
}

// Handler returns new handler for HTTP requests.
func (s *Server) Handler() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/room/new", s.newRoom)
	router.HandleFunc("/room/{id}/notifications", s.listRoomNotifs)
	router.HandleFunc("/room/{id}/join", s.joinRoom)

	return router
}

func (s *Server) newRoom(w http.ResponseWriter, req *http.Request) {
	room := s.Domain.CreateRoom()

	data, err := json.Marshal(newRoomResponse{
		ID: room.ID,
	})
	if err != nil {
		log.Println("Unable to marshal JSON output:", err)
		w.WriteHeader(500)
		return
	}

	if _, err = w.Write(data); err != nil {
		log.Println("Unable to write response:", err)
		w.WriteHeader(500)
		return
	}
}

func (s *Server) listRoomNotifs(w http.ResponseWriter, req *http.Request) {
	roomID, err := strconv.ParseUint(mux.Vars(req)["id"], 10, 64)
	if err != nil {
		log.Println("Unable to parse room ID:", err)
		w.WriteHeader(500)
		return
	}

	room := s.Domain.Room(roomID)
	if room == nil {
		w.WriteHeader(404)
		return
	}

	var lastID int
	lastIDStr := req.URL.Query().Get("lastID")
	if lastIDStr == "" {
		lastID = -1
	} else if lastID, err = strconv.Atoi(lastIDStr); err != nil {
		log.Println("Unable to parse last feed ID:", err)
		w.WriteHeader(500)
		return
	}

	var notifs []util.FeedItem
	var nextLastID int
	var waitChannel <-chan struct{}

	notifs, nextLastID, waitChannel = room.Feed.List(lastID)
	if len(notifs) == 0 {
		select {
		case <-waitChannel:
			notifs, nextLastID, _ = room.Feed.List(lastID)
		case <-time.After(s.DeferredRequestMaxDuration):
		}
	}

	data, err := json.Marshal(roomNotificationsResponse{
		Notifications: notifs,
		LastID:        nextLastID,
	})
	if err != nil {
		log.Println("Unable to marshal JSON output:", err)
		w.WriteHeader(500)
		return
	}

	if _, err = w.Write(data); err != nil {
		log.Println("Unable to write response:", err)
		w.WriteHeader(500)
		return
	}
}

func (s *Server) joinRoom(w http.ResponseWriter, req *http.Request) {
	roomID, err := strconv.ParseUint(mux.Vars(req)["id"], 10, 64)
	if err != nil {
		log.Println("Unable to parse room ID:", err)
		w.WriteHeader(500)
		return
	}

	room := s.Domain.Room(roomID)
	if room == nil {
		w.WriteHeader(404)
		return
	}

	decoder := json.NewDecoder(req.Body)
	var requestBody joinRoomRequest
	if err = decoder.Decode(&requestBody); err != nil {
		log.Println("Unable to decode body:", err)
		w.WriteHeader(500)
		return
	}

	room.Join(&requestBody.Member)
}
