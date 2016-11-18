
package bullet

import (
	"github.com/ikuo0/game/ebiten_act/event"
	"github.com/ikuo0/game/lib/anime"
	"github.com/ikuo0/game/lib/fig"
	"github.com/ikuo0/game/lib/move"
	"github.com/ikuo0/game/lib/radian"
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
	fig.Point
	V *move.Accel
	Anime *anime.Frames
	Vanished bool
	Endurance int
}

func (me *Bullet1) GetPoint() (fig.Point) {
	return me.Point
}

func (me *Bullet1) Direction() (radian.Radian) {
	return me.V.Radian
}

func (me *Bullet1) Update(trigger event.Trigger) {
	if me.Endurance < 1 {
		me.Vanish()
		trigger.EventTrigger(event.Vanishing1, nil, me)
	} else {
		me.V.Accel()
		p := me.V.Power()
		me.X += p.X
		me.Y += p.Y
		me.Anime.Update()
	}
}

func (me *Bullet1) Vanish() {
	me.Vanished = true
}
func (me *Bullet1) IsVanish() (bool) {
	return me.Vanished
}
func (me *Bullet1) Src() (x0, y0, x1, y1 int) {
	x := bullet1Source[me.Anime.Index()]
	return x.Left, x.Top, x.Right, x.Bottom
}
func (me *Bullet1) Dst() (x0, y0, x1, y1 int) {
	x, y := int(me.X) - 13, int(me.Y) - 13
	return x, y, x + 26, y + 26
}
func (me *Bullet1) HitRects() ([]fig.Rect) {
	x, y := int(me.X) - 13, int(me.Y) - 13
	return []fig.Rect{{x, y, x + 26, y + 26}}
}

func (me *Bullet1) Hit() {
	me.Endurance--
}

func NewBullet1(pt fig.Point, direction radian.Radian) (*Bullet1) {
	return &Bullet1 {
		Point: pt,
		V:          move.NewAccel(direction, 1, 0.5, 6),
		Anime:      anime.NewFrames(15, 15, 15, 15),
		Endurance:  1,
	}
}

