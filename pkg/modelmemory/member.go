package modelmemory

import "github.com/nieltg/quickshoot-party-match-server/pkg/model"

type member struct {
	payload model.MemberPayload
}

func (m *member) Payload() model.MemberPayload {
	return m.payload
}