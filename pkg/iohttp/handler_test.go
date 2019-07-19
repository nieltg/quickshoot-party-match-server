package iohttp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/mux"
	"github.com/nieltg/quickshoot-party-match-server/pkg/model"
	"github.com/nieltg/quickshoot-party-match-server/pkg/modelmemory"
)

var handler = Handler{
	DeferredRequestMaxDuration: 30 * time.Second,

	Domain: &modelmemory.Domain{
		JoinMaxDuration: 5 * time.Minute,
	},
}

var roomID uint64
var joinChannel = make(chan int, 2)

//TestCreateRoom is a unit test to test the capability of the system to create a room
func TestCreateRoom(t *testing.T) {
	requestData := newRoomRequest{
		Payload: model.RoomPayload{
			MaxMemberCount: 2,
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
		t.Fatal("WRONG RETURN BRO!")
	}

	roomID = responseData.ID
}

func TestFirstPlayerJoinRoom(t *testing.T) {
	TestCreateRoom(t)

	requestData := newRoomMemberRequest{
		Payload: model.MemberPayload{
			ID:   1,
			Name: "Giovanni Dejan",
		},
	}
	jsonSent, error := json.Marshal(requestData)
	if error != nil {
		t.Fatal("Fail to create jsonSent data!")
	}

	request, error := http.NewRequest(http.MethodPost, fmt.Sprintf("/room/new/%d/member/new", roomID), bytes.NewBuffer(jsonSent))
	request = mux.SetURLVars(request, map[string]string{"roomID": fmt.Sprintf("%d", roomID)})
	if error != nil {
		t.Fatal(request, " request can't be created!")
	}

	response := httptest.NewRecorder()

	server := http.HandlerFunc(handler.newRoomMember)

	func() {
		joinChannel <- 1
		server.ServeHTTP(response, request)
		if status := response.Code; status != http.StatusOK {
			t.Fatal("Can't join room! Status", status)
		}
	}()
}

func TestFirstPlayerJoinNonexistentRoom(t *testing.T) {
	requestData := newRoomMemberRequest{
		Payload: model.MemberPayload{
			ID:   1,
			Name: "Giovanni Dejan",
		},
	}
	jsonSent, error := json.Marshal(requestData)
	if error != nil {
		t.Fatal("Fail to create jsonSent data!")
	}

	request, error := http.NewRequest(http.MethodPost, fmt.Sprintf("/room/new/%d/member/new", 100), bytes.NewBuffer(jsonSent))
	request = mux.SetURLVars(request, map[string]string{"roomID": fmt.Sprintf("%d", 100)})
	if error != nil {
		t.Fatal(request, " request can't be created!")
	}

	response := httptest.NewRecorder()

	server := http.HandlerFunc(handler.newRoomMember)

	server.ServeHTTP(response, request)
	if status := response.Code; status != http.StatusNotFound {
		t.Fatal("Status isn't 404! Status", status)
	}
}

func TestFirstPlayerLeaveRoom(t *testing.T) {
	if roomID == 0 {
		TestCreateRoom(t)
	}
	TestFirstPlayerJoinRoom(t)

	request, error := http.NewRequest(http.MethodDelete, fmt.Sprintf("/room/new/%d/member/%d", roomID, 1), nil)
	request = mux.SetURLVars(request, map[string]string{"roomID": fmt.Sprintf("%d", roomID), "memberID": "1"})
	fmt.Println(mux.Vars(request))
	if error != nil {
		t.Fatal(request, " request can't be created!")
	}

	response := httptest.NewRecorder()

	server := http.HandlerFunc(handler.deleteRoomMember)

	go func(){
		server.ServeHTTP(response, request)
		if status := response.Code; status != http.StatusOK {
			t.Fatal("Can't leave room! Status", status)
		}
	}()
}

func TestSecondPlayerJoinRoom(t *testing.T) {
	if roomID == 0 {
		TestCreateRoom(t)
	}

	requestData := newRoomMemberRequest{
		Payload: model.MemberPayload{
			ID:   2,
			Name: "Daniel Pintara",
		},
	}
	jsonSent, error := json.Marshal(requestData)
	if error != nil {
		t.Fatal("Fail to create jsonSent data!")
	}

	request, error := http.NewRequest(http.MethodPost, fmt.Sprintf("/room/new/%d/member/new", roomID), bytes.NewBuffer(jsonSent))
	request = mux.SetURLVars(request, map[string]string{"roomID": fmt.Sprintf("%d", roomID)})
	if error != nil {
		t.Fatal(request, " request can't be created!")
	}

	response := httptest.NewRecorder()

	server := http.HandlerFunc(handler.newRoomMember)

	go func(){
		joinChannel <- 1
		server.ServeHTTP(response, request)
		if status := response.Code; status != http.StatusOK {
			t.Fatal("Can't join room! Status", status)
		}
	}()
}

func TestJoinAndLeaveRoomNotifs(t *testing.T) {
	TestFirstPlayerLeaveRoom(t)

	t.Logf("Room ID: %d", roomID)

	request, error := http.NewRequest(http.MethodGet, fmt.Sprintf("/room/%d/events", roomID), nil)
	request = mux.SetURLVars(request, map[string]string{"roomID": fmt.Sprintf("%d", roomID)})
	if error != nil {
		t.Fatal(request, " request can't be created!")
	}

	response := httptest.NewRecorder()

	server := http.HandlerFunc(handler.listRoomEvents)

	server.ServeHTTP(response, request)
	if status := response.Code; status != http.StatusOK {
		t.Fatal("Can't get list of room events! Status", status)
	}
}

func TestJoinRoomWhenFull(t *testing.T) {
	TestFirstPlayerJoinRoom(t)
	<- joinChannel
	TestSecondPlayerJoinRoom(t)
	<- joinChannel

	requestData := newRoomMemberRequest{
		Payload: model.MemberPayload{
			ID:   3,
			Name: "Daniel Agatan",
		},
	}
	jsonSent, error := json.Marshal(requestData)
	if error != nil {
		t.Fatal("Fail to create jsonSent data!")
	}

	request, error := http.NewRequest(http.MethodPost, fmt.Sprintf("/room/new/%d/member/new", roomID), bytes.NewBuffer(jsonSent))
	request = mux.SetURLVars(request, map[string]string{"roomID": fmt.Sprintf("%d", roomID)})
	if error != nil {
		t.Fatal(request, " request can't be created!")
	}

	response := httptest.NewRecorder()

	server := http.HandlerFunc(handler.newRoomMember)

	server.ServeHTTP(response, request)
	if status := response.Code; status != http.StatusForbidden {
		t.Fatal("Room isn't full! Status", status)
	}
}

func TestGameForTwoGamers(t *testing.T) {
	TestFirstPlayerJoinRoom(t)
	<- joinChannel
	TestSecondPlayerJoinRoom(t)
	<- joinChannel

	t.Log("Room ID: " + fmt.Sprint(roomID));

	var requestsData []newTapTimeRequest
	requestsData = append(requestsData, newTapTimeRequest{
		Payload: model.MemberTapTimePayload{
			TimeInMilis: 250,
		},
	})
	requestsData = append(requestsData, newTapTimeRequest{
		Payload: model.MemberTapTimePayload{
			TimeInMilis: 249,
		},
	})

	for i, requestData := range requestsData {
		jsonSent, error := json.Marshal(requestData)
		if error != nil {
			t.Fatal("Fail to create jsonSent data!")
		}

		request, error := http.NewRequest(http.MethodPost, fmt.Sprintf("/room/new/%d/member/%d/tap", roomID, (i + 1)), bytes.NewBuffer(jsonSent))
		request = mux.SetURLVars(request, map[string]string{"roomID": fmt.Sprintf("%d", roomID), "memberID": fmt.Sprintf("%d", (i + 1))})
		if error != nil {
			t.Fatal(request, " request can't be created!")
		}

		response := httptest.NewRecorder()

		server := http.HandlerFunc(handler.registerTapTime)

		go func(){
			server.ServeHTTP(response, request)
			if status := response.Code; status != http.StatusOK {
				t.Fatal("User can't tap! Status", status)
			}
		}()

		listenNotif(t)
	}
}

func listenNotif(t *testing.T) {
	request, error := http.NewRequest(http.MethodGet, fmt.Sprintf("/room/%d/events", roomID), nil)
	request = mux.SetURLVars(request, map[string]string{"roomID": fmt.Sprintf("%d", roomID)})
	if error != nil {
		t.Fatal(request, " request can't be created!")
	}

	response := httptest.NewRecorder()

	server := http.HandlerFunc(handler.listRoomEvents)

	server.ServeHTTP(response, request)
	if status := response.Code; status != http.StatusOK {
		t.Fatal("Can't get list of room events! Status", status)
	}
}