package domain

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

// List fetches all available items.
func (n *Feed) List() ([]FeedItem, int, <-chan struct{}) {
	n.itemsMutex.RLock()
	defer n.itemsMutex.RUnlock()

	return n.items, len(n.items) - 1, n.waitChannel
}

// ListAfter fetches all available items after lastID.
func (n *Feed) ListAfter(lastID uint64) ([]FeedItem, int, <-chan struct{}) {
	n.itemsMutex.RLock()
	defer n.itemsMutex.RUnlock()

	return n.items[lastID:], len(n.items) - 1, n.waitChannel
}
