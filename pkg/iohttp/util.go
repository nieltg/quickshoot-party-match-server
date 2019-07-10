package iohttp

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/nieltg/quickshoot-party-match-server/pkg/model"
)

func (s *Handler) fetchRoom(w http.ResponseWriter, roomIDStr string) model.Room {
	roomID, err := strconv.ParseUint(roomIDStr, 10, 64)
	if err != nil {
		log.Println("Unable to parse room ID:", err)
		w.WriteHeader(500)
		return nil
	}

	room := s.Domain.Room(roomID)
	if room == nil {
		w.WriteHeader(404)
		return nil
	}

	return room
}

func decodeJSONBody(w http.ResponseWriter, body io.ReadCloser, v interface{}) bool {
	decoder := json.NewDecoder(body)

	if err := decoder.Decode(v); err != nil {
		log.Println("Unable to decode body:", err)
		w.WriteHeader(500)
		return false
	}

	return true
}

func writeJSON(w http.ResponseWriter, v interface{}) bool {
	data, err := json.Marshal(v)
	if err != nil {
		log.Println("Unable to marshal JSON output:", err)
		w.WriteHeader(500)
		return false
	}

	w.Header().Set("Content-Type", "application/json")

	if _, err = w.Write(data); err != nil {
		log.Println("Unable to write response:", err)
		w.WriteHeader(500)
		return false
	}

	return true
}
