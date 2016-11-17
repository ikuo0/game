
package player

import (
	"github.com/ikuo0/game/ebiten_stg/effect"
	"github.com/ikuo0/game/ebiten_stg/eventid"
	"github.com/ikuo0/game/ebiten_stg/world"
	"github.com/ikuo0/game/lib/action"
	"github.com/ikuo0/game/lib/event"
	"github.com/ikuo0/game/lib/fig"
	"github.com/ikuo0/game/lib/ginput"
	"github.com/ikuo0/game/lib/gradian"
	"github.com/ikuo0/game/lib/kcmd"
	"github.com/ikuo0/game/lib/move"
	"github.com/ikuo0/game/lib/radian"
	"github.com/ikuo0/game/lib/script"
	"github.com/ikuo0/game/lib/sprites"
	"github.com/ikuo0/game/lib/timer"
	"github.com/hajimehoshi/ebiten"
	//"math"
	//"fmt"
)

var SrcSlopeLeft2  = fig.Rect {218, 0, 218 + 32, 36}
var SrcSlopeLeft1  = fig.Rect {250, 0, 250 + 32, 36}
var SrcCenter      = fig.Rect {378, 0, 378 + 32, 36}
var SrcSlopeRight1 = fig.Rect {314, 0, 314 + 32, 36}
var SrcSlopeRight2 = fig.Rect {282, 0, 282 + 32, 36}
var ShotCommand    = []ginput.InputBits {ginput.Nkey1, ginput.Key1}
var SheldCommand   = []ginput.InputBits {ginput.Nkey2, ginput.Key2}

type Player struct {
	fig.FloatPoint
	Vanished    bool
	CurrentSrc  fig.Rect
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

func (me *Player) Point() (fig.FloatPoint) {
	return me.FloatPoint
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
			me.XYcomponent.Set(me.V.X(), me.V.Y())
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

func (me *Player) Vanish() {
	me.Vanished = true
}
func (me *Player) IsVanish() (bool) {
	return me.Vanished
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
func (me *Player) SetPoint(pt fig.FloatPoint) {
	me.FloatPoint = pt
}

func (me *Player) HitRects() ([]fig.Rect) {
	if me.Dead || !me.Invisible.Up() {
		return nil
	} else {
		x, y := int(me.X), int(me.Y)
		return []fig.Rect{{x, y, x, y}}
	}
}

func (me *Player) Hit(obj action.Object) {
	me.Endurance--
	if me.Endurance <= 0 {
		me.Dead = true
		me.ReamExplosion = effect.NewReamExplosion(64, 180, me.FloatPoint)
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

func NewPlayer(pt fig.FloatPoint) (*Player) {
	return &Player{
		FloatPoint: pt,
		V:           move.NewVector(90, 7),
		XYcomponent: move.NewXYcomponent(0.2),
		Kbuffer:    &kcmd.Buffer{},
		Endurance:  100,
		Invisible:  timer.NewFrame(180),
		NowEntry:   timer.NewFrame(30),
	}
}


//########################################
//# Objects
//########################################
type Interface interface {
	Point() (fig.FloatPoint)
	Direction() (radian.Radian)
	Update(trigger event.Trigger)
	Vanish()
	IsVanish() (bool)
	Src() (x0, y0, x1, y1 int)
	Dst() (x0, y0, x1, y1 int)
	SetPoint(fig.FloatPoint)
	HitRects() ([]fig.Rect)
	Hit(action.Object)
	SetInput(ginput.InputBits)
}

type Objects struct {
	*sprites.Objects
}
func (me *Objects) Get(i int) (Interface) {
	return me.Objs[i].(Interface)
}
func (me *Objects) SetInput(i int, bits ginput.InputBits) {
	me.Get(i).SetInput(bits)
}
func (me *Objects) SetPoint(i int, pt fig.FloatPoint) {
	me.Get(i).SetPoint(pt)
}
func (me *Objects) DrawOption(i int) (*ebiten.DrawImageOptions) {
	sx0, sy0, sx1, sy1 := me.Src(i)
	dx0, dy0, dx1, dy1 := me.Dst(i)
	opt := ebiten.DrawImageOptions {
		ImageParts: sprites.NewOneSprites(sx0, sy0, sx1, sy1, dx0, dy0, dx1, dy1),
	}

	o := me.Objs[i]
	pt := o.Point()
	opt.GeoM.Translate(-pt.X, -pt.Y)
	opt.GeoM.Rotate(float64(o.Direction()))
	opt.GeoM.Translate(pt.X, pt.Y)
	return &opt
}
func NewObjects() (*Objects) {
	return &Objects {
		Objects: sprites.NewObjects(),
	}
}

