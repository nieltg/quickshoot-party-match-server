package domain

import "sync"

// Notification ...
type Notification interface{}

// Notifications ...
type Notifications struct {
	notifs      []Notification
	notifsMutex sync.RWMutex

	notifChannel chan struct{}
}

// NewNotifications ...
func NewNotifications() *Notifications {
	return &Notifications{
		notifChannel: make(chan struct{}),
	}
}

// Put ...
func (n *Notifications) Put(notif Notification) {
	n.notifsMutex.Lock()
	defer n.notifsMutex.Unlock()

	n.notifs = append(n.notifs, notif)

	close(n.notifChannel)
	n.notifChannel = make(chan struct{})
}

func (n *Notifications) listUnlocked(useLastID bool, lastID uint64) []Notification {
	if useLastID {
		return n.notifs[lastID:]
	}
	return n.notifs
}

// List ...
func (n *Notifications) List(useLastID bool, lastID uint64) []Notification {
	n.notifsMutex.RLock()
	defer n.notifsMutex.RUnlock()

	notifs := n.listUnlocked(useLastID, lastID)
	if len(notifs) > 0 {
		return notifs
	}

	<-n.notifChannel
	return n.listUnlocked(useLastID, lastID)
}
