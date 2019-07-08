package util

import "sync"

// FeedItem ...
type FeedItem interface{}

// Feed ...
type Feed struct {
	items      []FeedItem
	itemsMutex sync.RWMutex

	waitChannel chan struct{}
}

// NewFeed initializes this data structure.
func NewFeed() *Feed {
	return &Feed{
		waitChannel: make(chan struct{}),
	}
}

// Put a new item to this data structure.
func (n *Feed) Put(notif FeedItem) {
	n.itemsMutex.Lock()
	defer n.itemsMutex.Unlock()

	n.items = append(n.items, notif)

	close(n.waitChannel)
	n.waitChannel = make(chan struct{})
}

// List fetches all available items after last ID. (-1 if no last ID)
func (n *Feed) List(lastID int) ([]FeedItem, int, <-chan struct{}) {
	n.itemsMutex.RLock()
	defer n.itemsMutex.RUnlock()

	return n.items[lastID+1:], len(n.items) - 1, n.waitChannel
}
