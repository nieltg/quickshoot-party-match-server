package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nieltg/quickshoot-party-match-server/pkg/domain"
)

// Server represents a HTTP Server.
type Server struct {
	Domain domain.Domain
}

// Handler returns new handler for HTTP requests.
func (s *Server) Handler() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/room/new", s.newRoom)

	return router
}

func (s *Server) newRoom(writer http.ResponseWriter, req *http.Request) {
	room := s.Domain.CreateRoom()

	data, err := json.Marshal(newRoomResponse{
		ID: room.ID,
	})
	if err != nil {
		log.Println("Unable to marshal JSON output: ", err)
	}

	_, err = writer.Write(data)
	if err != nil {
		log.Println("Unable to write response: ", err)
	}
}
