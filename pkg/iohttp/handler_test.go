package iohttp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"encoding/json"
	"time"
	"bytes"
	"fmt"

	"github.com/gorilla/mux"
	"github.com/nieltg/quickshoot-party-match-server/pkg/modelmemory"
	"github.com/nieltg/quickshoot-party-match-server/pkg/model"
)

var handler = Handler{
	DeferredRequestMaxDuration: 30 * time.Second,

	Domain: &modelmemory.Domain{
		JoinMaxDuration: 5 * time.Minute,
	},
}

var roomID uint64

//TestCreateRoom is a unit test to test the capability of the system to create a room
func TestCreateRoom(t *testing.T) {
	requestData := newRoomRequest{
		Payload: model.RoomPayload{
			MaxMemberCount: 4,
		},
	}

	jsonSent, error := json.Marshal(requestData)
	if error != nil {
		t.Fatal("Fail to create jsonSent data!")
	}

	request, error := http.NewRequest(http.MethodPost, "/room/new", bytes.NewBuffer(jsonSent))
	if error != nil {
		t.Fatal(request, " request can't be created!")
	}

	response := httptest.NewRecorder()

	server := http.HandlerFunc(handler.newRoom)

	server.ServeHTTP(response, request)
	if status := response.Code; status != http.StatusOK {
		t.Fatal("Can't create room! Status: ", status)
	}
	var responseData newRoomResponse
	if error := json.Unmarshal(response.Body.Bytes(), &responseData); error != nil {
		t.Fatal("WRONG RETURN BRO!");
	}

	roomID = responseData.ID
}

func TestJoinRoom(t *testing.T) {
	TestCreateRoom(t)

	requestData := newRoomMemberRequest{
		Payload: model.MemberPayload{
			ID: 1,
			Name: "Giovanni Dejan",
		},
	}
	jsonSent, error := json.Marshal(requestData)
	if error != nil {
		t.Fatal("Fail to create jsonSent data!")
	}

	request, error := http.NewRequest(http.MethodPost, fmt.Sprintf("/room/new/%d/member/new", roomID), bytes.NewBuffer(jsonSent))
	request = mux.SetURLVars(request, map[string]string{"roomID": "1"})
	if error != nil {
		t.Fatal(request, " request can't be created!")
	}

	response := httptest.NewRecorder()

	server := http.HandlerFunc(handler.newRoomMember)

	server.ServeHTTP(response, request)
	if status := response.Code; status != http.StatusOK {
		t.Fatal("Can't join room! Status", status)
	}
}

func TestLeaveRoom(t *testing.T) {
	TestJoinRoom(t)

	request, error := http.NewRequest(http.MethodDelete, fmt.Sprintf("/room/new/%d/member/%d", roomID, 1), nil)
	request = mux.SetURLVars(request, map[string]string{"roomID": fmt.Sprintf("%d", roomID), "memberID": "1"})
	fmt.Println(mux.Vars(request))
	if error != nil {
		t.Fatal(request, " request can't be created!")
	}

	response := httptest.NewRecorder()

	server := http.HandlerFunc(handler.deleteRoomMember)

	server.ServeHTTP(response, request)
	if status := response.Code; status != http.StatusOK {
		t.Fatal("Can't leave room! Status", status)
	}
}