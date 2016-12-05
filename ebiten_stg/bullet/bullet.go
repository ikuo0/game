
package bullet

import (
	"github.com/ikuo0/game/ebiten_stg/eventid"
	"github.com/ikuo0/game/ebiten_stg/action"
	"github.com/ikuo0/game/lib/event"
	"github.com/ikuo0/game/lib/anime"
	"github.com/ikuo0/game/lib/fig"
	"github.com/ikuo0/game/lib/move"
	"github.com/ikuo0/game/lib/radian"
	//"fmt"
)

//########################################
//# Bullet1
//########################################
var bullet1Source = []fig.Rect {
	{
		35,
		3,
		35 + 26,
		3 + 26,
	},
	{
		35,
		35,
		35 + 26,
		35 + 26,
	},
	{
		35,
		35,
		35 + 26,
		35 + 26,
	},
	{
		35,
		99,
		35 + 26,
		99 + 26,
	},
}

type Bullet1 struct {
	action.Object
	V         *move.FixedVector
	Anime     *anime.Frames
	Endurance int
}

func (me *Bullet1) Update(trigger event.Trigger) {
	if me.Endurance < 1 {
		me.Vanish()
		trigger.EventTrigger(eventid.Vanishing1, nil, me)
	} else {
		me.V.Accel(0.2)
		me.X += me.V.X()
		me.Y += me.V.Y()
		me.Anime.Update()
	}
}

func (me *Bullet1) Src() (x0, y0, x1, y1 int) {
	x := bullet1Source[me.Anime.Index()]
	return int(x.Left), int(x.Top), int(x.Right), int(x.Bottom)
}
func (me *Bullet1) Dst() (x0, y0, x1, y1 int) {
	x, y := int(me.X) - 13, int(me.Y) - 13
	return x, y, x + 26, y + 26
}
func (me *Bullet1) HitRects() ([]fig.Rect) {
	x, y := me.X - 13, me.Y - 13
	return []fig.Rect{{x, y, x + 26, y + 26}}
}

func (me *Bullet1) Hit(obj action.Interface) {
	me.Endurance--
}

func NewBullet1(pt fig.Point, direction radian.Radian) (*Bullet1) {
	return &Bullet1 {
		//Point:      pt,
		//Radian:     direction,
		Object: action.Object {
			Point:  pt,
			Radian: direction,
		},
		V:          move.NewFixedVector(direction, 6),
		Anime:      anime.NewFrames(15, 15, 15, 15),
		Endurance:  1,
	}
}

