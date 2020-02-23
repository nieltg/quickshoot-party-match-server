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

func (feed *roomEventFeed) put(event model.RoomEvent) {
    feed.itemsMutex.Lock()
    defer feed.itemsMutex.Unlock()

    feed.items = append(feed.items, event)

    close(feed.waitChannel)
    feed.waitChannel = make(chan struct{})
}

func (feed *roomEventFeed) List(lastID int) ([]model.RoomEvent, int, <-chan struct{}) {
    feed.itemsMutex.RLock()
    defer feed.itemsMutex.RUnlock()

    return feed.items[lastID+1:], len(feed.items) - 1, feed.waitChannel
}
