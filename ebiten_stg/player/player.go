
package player

import (
	"github.com/ikuo0/game/ebiten_stg/effect"
	"github.com/ikuo0/game/ebiten_stg/eventid"
	"github.com/ikuo0/game/ebiten_stg/world"
	"github.com/ikuo0/game/lib/anime"
	"github.com/ikuo0/game/lib/event"
	"github.com/ikuo0/game/lib/fig"
	"github.com/ikuo0/game/lib/ginput"
	"github.com/ikuo0/game/lib/kcmd"
	"github.com/ikuo0/game/lib/move"
	"github.com/ikuo0/game/lib/orig"
	"github.com/ikuo0/game/lib/radian"
	"github.com/ikuo0/game/lib/script"
	"github.com/ikuo0/game/lib/sprites"
	"github.com/ikuo0/game/lib/timer"
	"github.com/hajimehoshi/ebiten"
	//"math"
	//"fmt"
)

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
	Hit()
	SetInput(ginput.InputBits)
}

var SrcSlopeLeft2  = fig.Rect {218, 0, 218 + 32, 36}
var SrcSlopeLeft1  = fig.Rect {250, 0, 250 + 32, 36}
var SrcCenter      = fig.Rect {378, 0, 378 + 32, 36}
var SrcSlopeRight1 = fig.Rect {314, 0, 314 + 32, 36}
var SrcSlopeRight2 = fig.Rect {282, 0, 282 + 32, 36}
var ShotCommand    = []ginput.InputBits {ginput.Nkey1, ginput.Key1}
var SheldCommand   = []ginput.InputBits {ginput.Nkey2, ginput.Key2}

type Player struct {
	fig.FloatPoint
	Vanished   bool
	CurrentSrc fig.Rect
	V          *move.Inertia
	InputBits  ginput.InputBits
	Kbuffer    *kcmd.Buffer
	FireFrame  int
	Endurance  int
	Dead       bool
	ReamExplosion *effect.ReamExplosion
	Invisible  *timer.Frame
	NowEntry   *timer.Frame
}

func (me *Player) Point() (fig.FloatPoint) {
	return me.FloatPoint
}

func (me *Player) Direction() (radian.Radian) {
	return radian.Up()
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
			me.V.Radian = radian.LeftUp()
		} else if bits.And(ginput.Left | ginput.Down) {
			me.V.Radian = radian.LeftDown()
		} else if bits.And(ginput.Right | ginput.Up) {
			me.V.Radian = radian.RightUp()
		} else if bits.And(ginput.Right | ginput.Down) {
			me.V.Radian = radian.RightDown()
		} else if bits.And(ginput.Left) {
			me.V.Radian = radian.Left()
		} else if bits.And(ginput.Down) {
			me.V.Radian = radian.Down()
		} else if bits.And(ginput.Right) {
			me.V.Radian = radian.Right()
		} else if bits.And(ginput.Up) {
			me.V.Radian = radian.Up()
		}

		if bits.Or(ginput.AxisMask) {
			me.V.Accel()
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

		me.V.Chafe(0.2)
		p := me.V.Power()
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

func (me *Player) Hit() {
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
		V:          move.NewInertia(radian.Up(), 0, 0.6, 7),
		Kbuffer:    &kcmd.Buffer{},
		Endurance:  100,
		Invisible:  timer.NewFrame(180),
		NowEntry:   timer.NewFrame(30),
	}
}


//########################################
//# Players
//########################################
type Players struct {
	Objs     []Interface
}

func (me *Players) Len() (int) {
	return len(me.Objs)
}
func (me *Players) Src(i int) (x0, y0, x1, y1 int) {
	return me.Objs[i].Src()
}
func (me *Players) Dst(i int) (x0, y0, x1, y1 int) {
	return me.Objs[i].Dst()
}
func (me *Players) Update(i int, trigger event.Trigger) {
	me.Objs[i].Update(trigger)
}
func (me *Players) Origin(i int) (orig.Interface) {
	return me.Objs[i]
}
func (me *Players) HitRects(i int) ([]fig.Rect) {
	return me.Objs[i].HitRects()
}

func (me *Players) Hit(i int) {
	me.Objs[i].Hit()
}
func (me *Players) Vanish(i int) {
	me.Objs[i].Vanish()
}
func (me *Players) Clean(i int) {
	newObjs := []Interface{}
	for _, v := range me.Objs {
		if !v.IsVanish() {
			newObjs = append(newObjs, v)
		}
	}
	me.Objs = newObjs
}
func (me *Players) Options() (*ebiten.DrawImageOptions) {
	return &ebiten.DrawImageOptions {
		ImageParts: me,
	}
}

func (me *Players) DrawOption(i int) (*ebiten.DrawImageOptions) {
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

func (me *Players) Occure(objIf Interface) {
	me.Objs = append(me.Objs, objIf)
}

func (me *Players) SetInput(i int, bits ginput.InputBits) {
	me.Objs[i].SetInput(bits)
}
func (me *Players) SetPoint(i int, pt fig.FloatPoint) {
	me.Objs[i].SetPoint(pt)
}

func NewPlayers() (*Players) {
	return &Players {}
}

//########################################
//# Shot
//########################################
var SrcShot = fig.Rect {0, 66, 0 + 60, 66 + 66}
type Shot struct {
	fig.FloatPoint
	Vanished   bool
	V          *move.Constant
	Endurance  int
}

func (me *Shot) Point() (fig.FloatPoint) {
	return me.FloatPoint
}

func (me *Shot) Direction() (radian.Radian) {
	return radian.Up()
}

func (me *Shot) Update(trigger event.Trigger) {
	p := me.V.Power()
	me.X += p.X
	me.Y += p.Y
}

func (me *Shot) Vanish() {
	me.Vanished = true
}
func (me *Shot) IsVanish() (bool) {
	return me.Vanished
}
func (me *Shot) Src() (x0, y0, x1, y1 int) {
	return SrcShot.Left, SrcShot.Top, SrcShot.Right, SrcShot.Bottom
}
func (me *Shot) Dst() (x0, y0, x1, y1 int) {
	x, y := int(me.X) - 24, int(me.Y) - 30
	return x, y, x + 48, y + 60
}
func (me *Shot) SetPoint(pt fig.FloatPoint) {
	me.FloatPoint = pt
}
func (me *Shot) HitRects() ([]fig.Rect) {
	if me.Endurance <= 0 {
		return nil
	} else {
		x, y := int(me.X) - 24, int(me.Y) - 30
		return []fig.Rect{{x, y, x + 48, y + 60}}
	}
}

func (me *Shot) Hit() {
	me.Endurance--
	me.Vanish()
}

func (me *Shot) Stack() (*script.Stack) {
	return nil
}

func NewShot(pt fig.FloatPoint) (*Shot) {
	return &Shot{
		FloatPoint: pt,
		V:          move.NewConstant(radian.Up(), 32),
		Endurance:  6,
	}
}

//########################################
//# Shield
//########################################
var SrcSheld = fig.Rect {64, 64, 64 + 320, 64 + 320}
type Sheld struct {
	fig.FloatPoint
	Vanished   bool
	Anime     *anime.Frames
	V          *move.Accel
}

func (me *Sheld) Point() (fig.FloatPoint) {
	return me.FloatPoint
}

func (me *Sheld) Direction() (radian.Radian) {
	return radian.Up()
}

func (me *Sheld) Update(trigger event.Trigger) {
	me.V.Accel()
	p := me.V.Power()
	me.X += p.X
	me.Y += p.Y
	me.Anime.Update()
}

func (me *Sheld) Vanish() {
	me.Vanished = true
}
func (me *Sheld) IsVanish() (bool) {
	return me.Vanished
}
func (me *Sheld) Src() (x0, y0, x1, y1 int) {
	return SrcSheld.Left, SrcSheld.Top, SrcSheld.Right, SrcSheld.Bottom
}
func (me *Sheld) Dst() (x0, y0, x1, y1 int) {
// 96
// 16 + 80
	width := me.Anime.Index() * 15 + 36
	adjust := width / 2
	x, y := int(me.X) - adjust, int(me.Y) - adjust
	return x, y, x + width, y + width
}
func (me *Sheld) SetPoint(pt fig.FloatPoint) {
	me.FloatPoint = pt
}
func (me *Sheld) HitRects() ([]fig.Rect) {
	if me.IsVanish() {
		return nil
	} else {
		x, y := int(me.X) - 48, int(me.Y) - 48
		return []fig.Rect{{x, y, x + 96, y + 96}}
	}
}

func (me *Sheld) Hit() {
}

func (me *Sheld) Pushed() {
	me.V.Speed.MaxPower += 0.2
}

func (me *Sheld) Stack() (*script.Stack) {
	return nil
}

func NewSheld(pt fig.FloatPoint) (*Sheld) {
	return &Sheld{
		FloatPoint: pt,
		Anime:      anime.NewFrames(8, 7, 3, 2, 8),
		V:          move.NewAccel(radian.Up(), 0.1, 0.1, 0.7),
	}
}


type ShieldInterface interface {
	sprites.Object
	Pushed()
}

type Shields struct {
	*sprites.Objects
}

func (me *Shields) HitRects(i int) ([]fig.Rect) {
	return me.Objects.Objs[i].HitRects()
}

func (me *Shields) Pushed(i int) {
	me.Objs[i].(ShieldInterface).Pushed()
}

func NewShelds() (*Shields) {
	return &Shields {
		Objects: sprites.NewObjects(),
	}
}
