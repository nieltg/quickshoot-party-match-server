package event

import "github.com/nieltg/quickshoot-party-match-server/pkg/domain"

// MemberJoin builds a member-join event value.
func MemberJoin(member domain.RoomMember) ValueWithPayload {
	return ValueWithPayload{
		Type: TypeMemberJoin,
		Payload: &struct {
			Member domain.RoomMember
		}{
			Member: member,
		},
	}
}

// MemberLeave builds a member-leave event value.
func MemberLeave(memberID uint64) ValueWithPayload {
	return ValueWithPayload{
		Type: TypeMemberLeave,
		Payload: &struct {
			MemberID uint64
		}{
			MemberID: memberID,
		},
	}
}

// GameBegin builds a game-begin event value.
func GameBegin() Value {
	return Value{
		Type: TypeGameBegin,
	}
}

// MemberTapTime builds a member-tap-time event value.
func MemberTapTime(memberID uint64) ValueWithPayload {
	return ValueWithPayload{
		Type: TypeMemberTapTime,
		Payload: &struct {
			MemberID uint64
		}{
			MemberID: memberID,
		},
	}
}

// GameEnd builds a game-end event value.
func GameEnd(bestTapTime uint) ValueWithPayload {
	return ValueWithPayload{
		Type: TypeGameEnd,
		Payload: &struct {
			BestTapTime uint
		}{
			BestTapTime: bestTapTime,
		},
	}
}
