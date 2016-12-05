
package player

import (
	"github.com/ikuo0/game/ebiten_stg/action"
	"github.com/ikuo0/game/ebiten_stg/effect"
	"github.com/ikuo0/game/ebiten_stg/eventid"
	"github.com/ikuo0/game/ebiten_stg/world"
	"github.com/ikuo0/game/lib/event"
	"github.com/ikuo0/game/lib/fig"
	"github.com/ikuo0/game/lib/ginput"
	"github.com/ikuo0/game/lib/gradian"
	"github.com/ikuo0/game/lib/kcmd"
	"github.com/ikuo0/game/lib/move"
	"github.com/ikuo0/game/lib/radian"
	"github.com/ikuo0/game/lib/script"
	"github.com/ikuo0/game/lib/timer"
	//"github.com/hajimehoshi/ebiten"
	//"math"
	//"fmt"
)

var SrcSlopeLeft2  = fig.IntRect {218, 0, 218 + 32, 36}
var SrcSlopeLeft1  = fig.IntRect {250, 0, 250 + 32, 36}
var SrcCenter      = fig.IntRect {378, 0, 378 + 32, 36}
var SrcSlopeRight1 = fig.IntRect {314, 0, 314 + 32, 36}
var SrcSlopeRight2 = fig.IntRect {282, 0, 282 + 32, 36}
var ShotCommand    = []ginput.InputBits {ginput.Nkey1, ginput.Key1}
var SheldCommand   = []ginput.InputBits {ginput.Nkey2, ginput.Key2}

type Player struct {
	action.Object
	CurrentSrc  fig.IntRect
	V           *move.Vector
	XYcomponent *move.XYcomponent
	InputBits   ginput.InputBits
	Kbuffer     *kcmd.Buffer
	FireFrame   int
	Endurance   int
	Dead        bool
	ReamExplosion *effect.ReamExplosion
	Invisible   *timer.Frame
	NowEntry    *timer.Frame
}

func (me *Player) Direction() (radian.Radian) {
	return gradian.Up()
}

func (me *Player) Update(trigger event.Trigger) {
	if me.Dead {
		if me.ReamExplosion.Update(trigger) {
			me.Vanish()
			trigger.EventTrigger(eventid.PlayerDied, nil, nil)
		}
	} else if !me.NowEntry.Up() {
		me.Y -= 2
		me.SetCurrentSrc(0)
	} else {
		world.SetPlayer(me)

		bits := me.InputBits

		if bits.And(ginput.Left | ginput.Up) {
			me.V.Degree.Deg = 135
		} else if bits.And(ginput.Left | ginput.Down) {
			me.V.Degree.Deg = -135
		} else if bits.And(ginput.Right | ginput.Up) {
			me.V.Degree.Deg = 45
		} else if bits.And(ginput.Right | ginput.Down) {
			me.V.Degree.Deg = -45
		} else if bits.And(ginput.Left) {
			me.V.Degree.Deg = 180
		} else if bits.And(ginput.Down) {
			me.V.Degree.Deg = -90
		} else if bits.And(ginput.Right) {
			me.V.Degree.Deg = 0
		} else if bits.And(ginput.Up) {
			me.V.Degree.Deg = 90
		}

		if bits.Or(ginput.AxisMask) {
			me.V.Accel(7)
			me.XYcomponent.Accel(me.V.X(), me.V.Y())
		} else {
			//me.V.Frictional(0.2)
		}

		me.Kbuffer.Update(bits)
		if me.FireFrame == 0 && kcmd.Check(ShotCommand, me.Kbuffer, 1) {
			me.FireFrame = 3
		}

		if me.FireFrame > 0 {
			trigger.EventTrigger(eventid.Shot, nil, me)
			me.FireFrame--
		}

		if kcmd.Check(SheldCommand, me.Kbuffer, 1) {
			trigger.EventTrigger(eventid.Sheld, nil, me)
		}

		me.XYcomponent.Update()
		p := me.XYcomponent.Power()
		me.X += p.X
		me.Y += p.Y

		me.SetCurrentSrc(p.X)
	}
}

func (me *Player) Src() (x0, y0, x1, y1 int) {
	r := me.CurrentSrc
	return r.Left, r.Top, r.Right, r.Bottom
}
func (me *Player) Dst() (x0, y0, x1, y1 int) {
	x, y := int(me.X) - 16, int(me.Y) - 18
	return x, y, x + 32, y + 36
}
func (me *Player) SetInput(bits ginput.InputBits) {
	me.InputBits = bits
}

func (me *Player) HitRects() ([]fig.Rect) {
	if me.Dead || !me.Invisible.Up() {
		return nil
	} else {
		x, y := me.X, me.Y
		return []fig.Rect{{x, y, x, y}}
	}
}

func (me *Player) Hit(obj action.Interface) {
	me.Endurance--
	if me.Endurance <= 0 {
		me.Dead = true
		me.ReamExplosion = effect.NewReamExplosion(64, 180, me.Point)
	}
}
func (me *Player) Stack() (*script.Stack) {
	return nil
}

func (me *Player) SetCurrentSrc(p float64) {
	if p >= 6 {
		me.CurrentSrc = SrcSlopeRight2
	} else if p >= 3 {
		me.CurrentSrc = SrcSlopeRight1
	} else if p <= -6 {
		me.CurrentSrc = SrcSlopeLeft2
	} else if p <= -3 {
		me.CurrentSrc = SrcSlopeLeft1
	} else {
		me.CurrentSrc = SrcCenter
	}
}

func NewPlayer(pt fig.Point) (*Player) {
	return &Player{
		Object: action.Object {
			Point: pt,
		},
		V:           move.NewVector(90, 7),
		XYcomponent: move.NewXYcomponent(0.4),
		Kbuffer:    &kcmd.Buffer{},
		Endurance:  100,
		Invisible:  timer.NewFrame(180),
		NowEntry:   timer.NewFrame(30),
	}
}

