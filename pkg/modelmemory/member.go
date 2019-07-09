package modelmemory

import "github.com/nieltg/quickshoot-party-match-server/pkg/model"

type member struct {
	payload model.MemberPayload
}

func (m *member) RegisterTapTime(tapTimeMilis uint64) {
	// TODO: Do something here.
}
