
package eventid

import (
	"github.com/ikuo0/game/lib/event"
)

const (
	Player      event.Id = event.User + iota
	PlayerDied
	StageClear
	Shot
	Jump
	Beat
	Enemy
	Explosion1
	Explosion2
	Vortex
	VortexTaken
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
		case Jump: return "Jump"
		case Beat: return "Beat"
		case Enemy: return "Enemy"
		case Explosion1: return "Explosion1"
		case Explosion2: return "Explosion2"
		case Vortex: return "Vortex"
		case VortexTaken: return "VortexTaken"
		case BigNumber: return "BigNumber"
		case Score: return "Score"
		case Block: return "Block"
		case OccureBlock: return "OccureBlock"
	}
	return "Unknown"
}

