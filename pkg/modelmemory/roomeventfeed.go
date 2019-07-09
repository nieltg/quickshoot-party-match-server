package modelmemory

import (
	"sync"

	"github.com/nieltg/quickshoot-party-match-server/pkg/model"
)

type roomEventFeed struct {
	items      []model.RoomEvent
	itemsMutex sync.RWMutex

	waitChannel chan struct{}
}

func newRoomEventFeed() *roomEventFeed {
	return &roomEventFeed{
		waitChannel: make(chan struct{}),
	}
}

func (f *roomEventFeed) put(event model.RoomEvent) {
	f.itemsMutex.Lock()
	defer f.itemsMutex.Unlock()

	f.items = append(f.items, event)

	close(f.waitChannel)
	f.waitChannel = make(chan struct{})
}

func (f *roomEventFeed) List(lastID int) ([]model.RoomEvent, int, <-chan struct{}) {
	f.itemsMutex.RLock()
	defer f.itemsMutex.RUnlock()

	return f.items[lastID+1:], len(f.items) - 1, f.waitChannel
}
