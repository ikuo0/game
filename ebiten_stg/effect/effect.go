
package effect

import (
	"github.com/ikuo0/game/ebiten_stg/eventid"
	"github.com/ikuo0/game/lib/event"
	"github.com/ikuo0/game/lib/fig"
	"github.com/ikuo0/game/lib/radian"
	"github.com/ikuo0/game/lib/timer"
	"math/rand"
)

//########################################
//# ReamExplosion
//########################################
type ReamExplosion struct {
	fig.Point
	Width        int
	ExplodeTimer timer.Frame
	VanishTimer  timer.Frame
}

func (me *ReamExplosion) GetPoint() (fig.Point) {
	return me.Point
}
func (me *ReamExplosion) Direction() (radian.Radian) {
	return 0
}
func (me *ReamExplosion) Update(trigger event.Trigger) (bool) {
	if me.VanishTimer.Up() {
		return true
	} else if me.ExplodeTimer.Up() {
		me.ExplodeTimer.Start(8 + rand.Intn(10))
		w := me.Width
		wh := float64(w / 2)
		trigger.EventTrigger(eventid.Explosion1, fig.Point{float64(rand.Intn(w)) - wh, float64(rand.Intn(w)) - wh}, me)
	}
	return false
}

func NewReamExplosion(width, tup int, pt fig.Point) (*ReamExplosion) {
	res := ReamExplosion {
		Point: pt,
		Width: width,
	}
	res.VanishTimer.Start(tup)
	return &res
}


