package event

// Type represents type of an event.
type Type string

// Available event types.
const (
	TypeMemberJoin    Type = "MEMBER_JOIN"
	TypeMemberLeave   Type = "MEMBER_LEAVE"
	TypeGameBegin     Type = "GAME_BEGIN"
	TypeMemberTapTime Type = "MEMBER_TAP_TIME"
	TypeGameEnd       Type = "GAME_END"
)
