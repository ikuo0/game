
package eventid

import (
	"github.com/ikuo0/game/lib/event"
)

const (
	Player      event.Id = event.User + iota
	PlayerDied
	StageClear
	Shot
	Sheld
	Enemy
	Explosion1
	Explosion2
	Vortex
	BigNumber
	Score
	Block
	OccureBlock
)

func String(id event.Id) (string) {
	switch id {
		case event.Wait: return "Wait"
		case event.Jump: return "Jump"

		case PlayerDied: return "PlayerDied"
		case StageClear: return "StageClear"
		case Shot: return "Shot"
		case Sheld: return "Sheld"
		case Enemy: return "Enemy"
		case Explosion1: return "Explosion1"
		case Explosion2: return "Explosion2"
		case Vortex: return "Vortex"
		case BigNumber: return "BigNumber"
		case Score: return "Score"
		case Block: return "Block"
		case OccureBlock: return "OccureBlock"
	}
	return "Unknown"
}

