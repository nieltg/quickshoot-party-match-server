package iohttp

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/nieltg/quickshoot-party-match-server/pkg/model"
)

func (s *Handler) fetchRoom(writer http.ResponseWriter, roomIDStr string) model.Room {
	roomID, err := strconv.ParseUint(roomIDStr, 10, 64)
	if err != nil {
		log.Println("Unable to parse room ID:", err)
		writer.WriteHeader(500)
		return nil
	}

	room := s.Domain.Room(roomID)
	if room == nil {
		writer.WriteHeader(404)
		return nil
	}

	return room
}

func decodeJSONBody(writer http.ResponseWriter, body io.ReadCloser, v interface{}) bool {
	decoder := json.NewDecoder(body)

	if err := decoder.Decode(v); err != nil {
		log.Println("Unable to decode body:", err)
		writer.WriteHeader(500)
		return false
	}

	return true
}

func writeJSON(writer http.ResponseWriter, v interface{}) bool {
	data, err := json.Marshal(v)
	if err != nil {
		log.Println("Unable to marshal JSON output:", err)
		writer.WriteHeader(500)
		return false
	}

	writer.Header().Set("Content-Type", "application/json")

	if _, err = writer.Write(data); err != nil {
		log.Println("Unable to write response:", err)
		writer.WriteHeader(500)
		return false
	}

	return true
}
