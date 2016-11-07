
package eventid

import (
	"github.com/ikuo0/game/lib/event"
)

const (
	Player event.Id = iota + event.User
	PlayerDied
	StageClear
	Shot
	Sheld
	Bullet1
	Bullet2
	Bullet3
	Heli1
	Heli2
	Heli3
	Aide1
	Aide2
	Aide3
	Boss1
	Boss2
	Boss3
	Explosion1
	Explosion2
	Explosion3
	Vanishing1
	Vanishing21
	Vanishing3
	BigNumber
	Score
	BgmPlay
)

func String(id event.Id) (string) {
	switch id {
		case event.Wait: return "Wait"
		case event.Jump: return "Jump"
		case Player: return "Player"
		case Shot: return "Shot"
		case Bullet1: return "Bullet1"
		case Bullet2: return "Bullet2"
		case Bullet3: return "Bullet3"
		case Heli1: return "Heli1"
		case Heli2: return "Heli2"
		case Heli3: return "Heli3"
		case Aide1: return "Aide1"
		case Aide2: return "Aide2"
		case Aide3: return "Aide3"
		case Boss1: return "Boss1"
		case Boss2: return "Boss2"
		case Boss3: return "Boss3"
		case Explosion1: return "Explosion1"
		case Explosion2: return "Explosion2"
		case Explosion3: return "Explosion3"
		case BigNumber: return "BigNumber"
		case Score: return "Score"
		case BgmPlay: return "BgmPlay"
	}
	return "Unknown"
}

