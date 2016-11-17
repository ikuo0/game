
package eventid

import (
	"github.com/ikuo0/game/lib/event"
)

const (
	Player event.Id = iota + event.User
	Wall
	StageClear
	CollisionWall
	BigNumber
	Score
	BgmPlay
)

func String(id event.Id) (string) {
	switch id {
		case event.Wait: return "Wait"
		case event.Jump: return "Jump"
		case Player: return "Player"
		case StageClear: return "StageClear"
		case CollisionWall: return "CollisionWall"
		case BigNumber: return "BigNumber"
		case Score: return "Score"
		case BgmPlay: return "BgmPlay"
	}
	return "Unknown"
}

