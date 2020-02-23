package iohttp

import (
    "encoding/json"
    "io"
    "log"
    "net/http"
    "strconv"

    "github.com/nieltg/quickshoot-party-match-server/pkg/model"
)

func (handler *Handler) fetchRoom(writer http.ResponseWriter, roomIDStr string) model.Room {
    roomID, err := strconv.ParseUint(roomIDStr, 10, 64)
    if err != nil {
        log.Println("Unable to parse room ID:", err)
        writer.WriteHeader(http.StatusInternalServerError)
        return nil
    }

    room := handler.Domain.Room(roomID)
    if room == nil {
        writer.WriteHeader(http.StatusNotFound)
        return nil
    }

    return room
}

func decodeJSONBody(writer http.ResponseWriter, body io.ReadCloser, value interface{}) bool {
    decoder := json.NewDecoder(body)

    if err := decoder.Decode(value); err != nil {
        log.Println("Unable to decode body:", err)
        writer.WriteHeader(http.StatusInternalServerError)
        return false
    }

    return true
}

func writeJSON(writer http.ResponseWriter, value interface{}) bool {
    data, err := json.Marshal(value)
    if err != nil {
        log.Println("Unable to marshal JSON output:", err)
        writer.WriteHeader(http.StatusInternalServerError)
        return false
    }

    writer.Header().Set("Content-Type", "application/json")

    if _, err = writer.Write(data); err != nil {
        log.Println("Unable to write response:", err)
        writer.WriteHeader(http.StatusInternalServerError)
        return false
    }

    return true
}
