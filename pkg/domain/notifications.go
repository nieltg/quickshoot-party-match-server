package domain

import "sync"

// Notification ...
type Notification interface{}

// Notifications ...
type Notifications struct {
	notifs      []Notification
	notifsMutex sync.RWMutex

	WaitChannel chan struct{}
}

// NewNotifications initializes this data structure.
func NewNotifications() *Notifications {
	return &Notifications{
		WaitChannel: make(chan struct{}),
	}
}

// Put a new notification to this data structure and notify listeners.
func (n *Notifications) Put(notif Notification) {
	n.notifsMutex.Lock()
	defer n.notifsMutex.Unlock()

	n.notifs = append(n.notifs, notif)

	close(n.WaitChannel)
	n.WaitChannel = make(chan struct{})
}

// RLock locks data structure for reading.
func (n *Notifications) RLock() {
	n.notifsMutex.RLock()
}

// RUnlock unlocks data structure so new notification can be put to.
func (n *Notifications) RUnlock() {
	n.notifsMutex.RUnlock()
}

// List fetches all available notifications. Call RLock first.
func (n *Notifications) List() []Notification {
	return n.notifs
}

// ListAfter fetches all available notifications after lastID. Call RLock first.
func (n *Notifications) ListAfter(lastID uint64) []Notification {
	return n.notifs[lastID:]
}
