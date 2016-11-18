
package vortex

import (
	"github.com/ikuo0/game/ebiten_act/eventid"
	"github.com/ikuo0/game/lib/action"
	"github.com/ikuo0/game/lib/event"
	"github.com/ikuo0/game/lib/script"
	"github.com/ikuo0/game/lib/fig"
	"github.com/ikuo0/game/lib/radian"
)

//########################################
//# Vortex
//########################################
const Width = 32
const Height = 32
const AdjustX = -16
const AdjustY = -16

var ImageSource = []fig.IntRect {
	{
		0,
		0,
		Width,
		Height,
	},
}

type Vortex struct {
	fig.Point
	Taken    bool
	Vanished bool
}

func (me *Vortex) GetPoint() (fig.Point) {
	return me.Point
}

func (me *Vortex) Direction() (radian.Radian) {
	return 0
}

func (me *Vortex) Update(trigger event.Trigger) {
	if me.Taken {
		trigger.EventTrigger(eventid.VortexTaken, nil, me)
		me.Vanish()
	}
}

func (me *Vortex) Vanish() {
	me.Vanished = true
}

func (me *Vortex) IsVanish() (bool) {
	return me.Vanished
}
func (me *Vortex) Src() (x0, y0, x1, y1 int) {
	x := ImageSource[0]
	return x.Left, x.Top, x.Right, x.Bottom
}
func (me *Vortex) Dst() (x0, y0, x1, y1 int) {
	x, y := int(me.X) + AdjustX, int(me.Y) + AdjustY
	return x, y, x + Width, y + Height
}
func (me *Vortex) HitRects() ([]fig.Rect) {
	x, y := me.X + AdjustX, me.Y + AdjustY
	return []fig.Rect{{x, y, x + Width, y + Height}}
}

func (me *Vortex) Hit(obj action.Object) {
	me.Taken = true
}

func (me *Vortex) Stack() (*script.Stack) {
	return nil
}

func New(pt fig.Point) (*Vortex) {
	return &Vortex {
		Point: pt,
	}
}
