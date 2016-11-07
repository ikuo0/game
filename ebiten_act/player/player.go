package player

// 漢字

import (
	"github.com/ikuo0/game/ebiten_act/eventid"
	"github.com/ikuo0/game/ebiten_act/funcs"
	"github.com/ikuo0/game/ebiten_act/world"
	"github.com/ikuo0/game/lib/action"
	"github.com/ikuo0/game/lib/anime"
	"github.com/ikuo0/game/lib/event"
	"github.com/ikuo0/game/lib/fig"
	"github.com/ikuo0/game/lib/ginput"
	"github.com/ikuo0/game/lib/kcmd"
	"github.com/ikuo0/game/lib/move"
	"github.com/ikuo0/game/lib/radian"
	"github.com/ikuo0/game/lib/script"
	"github.com/ikuo0/game/lib/sprites"
	"github.com/ikuo0/game/lib/timer"
	//"fmt"
)

const Width = 48
const Height = 64
const AdjustX = -24
const AdjustY = -64

var ImageSources []fig.Rect = []fig.Rect {
	{0, 0, 48, 64},
	{48, 0, 48 + 48, 64},
	{48, 0, 0, 64},
	{48 + 48, 0, 48, 64},
}

var ShotCommand    = []ginput.InputBits {ginput.Nkey1, ginput.Key1}
var JumpCommand    = []ginput.InputBits {ginput.Nkey2, ginput.Key2}

type Gun struct {
	fig.FloatPoint
	FaceDirection funcs.FaceDirection
}
func (me *Gun) Point() (fig.FloatPoint) {
	return me.FloatPoint
}

func (me *Gun) Direction() (radian.Radian) {
	if me.FaceDirection == funcs.FaceLeft {
		return radian.Left()
	} else {
		return radian.Right()
	}
}

func (me *Gun) SetPoint(pt fig.FloatPoint) {
	me.X = pt.X
	me.Y = pt.Y - 32
}

func (me *Gun) SetLeft(pt fig.FloatPoint) {
	me.SetPoint(pt)
	me.FaceDirection = funcs.FaceLeft
}

func (me *Gun) SetRight(pt fig.FloatPoint) {
	me.SetPoint(pt)
	me.FaceDirection = funcs.FaceRight
}

type Player struct {
	fig.FloatPoint
	FaceDirection funcs.FaceDirection
	Vanished      bool
	Dead          bool
	CanJump       bool
	Beaten        bool
	Blackout      bool
	BlackoutTimer   *timer.Frame
	CurrentSrc    fig.Rect
	V             *move.FallingInertia
	InputBits     ginput.InputBits
	Kbuffer       *kcmd.Buffer
	Endurance     int
	Anime         *anime.Frames
	FrameCounter  int
	FallingRects  *funcs.FallingRects
	Gun           Gun
}

func (me *Player) Point() (fig.FloatPoint) {
	return me.FloatPoint
}

func (me *Player) Direction() (radian.Radian) {
	return radian.Up()
}

func (me *Player) Update(trigger event.Trigger) {
	defer func () {
		me.FrameCounter++
	} ()

	if me.Dead {
		trigger.EventTrigger(eventid.Explosion2, nil, nil)
		me.Vanish()
	} else {

		bits := me.InputBits
		me.Kbuffer.Update(bits)

		if me.Blackout {
			if me.BlackoutTimer.Up() {
				me.Blackout = false
				me.Beaten = false
			}
		} else if me.Beaten {
			trigger.EventTrigger(eventid.Beat, nil, me)
			me.Blackout = true
			me.BlackoutTimer.Start(60)
		} else {
		}

		if bits.And(ginput.Left) {
			me.V.Radian = radian.Left()
		} else if bits.And(ginput.Right) {
			me.V.Radian = radian.Right()
		}

		if bits.Or(ginput.Left | ginput.Right) {
			me.V.Accel()
			me.Anime.Update()
			if bits.Or(ginput.Left) {
				me.FaceDirection = funcs.FaceLeft
			} else {
				me.FaceDirection = funcs.FaceRight
			}
		}

		if kcmd.Check(ShotCommand, me.Kbuffer, 1) {
			if me.FaceDirection == funcs.FaceLeft {
				me.Gun.SetLeft(me.FloatPoint)
			} else {
				me.Gun.SetRight(me.FloatPoint)
			}
			trigger.EventTrigger(eventid.Shot, &me.Gun, me)
		}

		if me.CanJump && kcmd.Check(JumpCommand, me.Kbuffer, 1) {
			trigger.EventTrigger(eventid.Jump, nil, me)
			me.V.Jump(17)
		}

		me.V.Fall()
		me.V.Chafe(0.4)
		p := me.V.Power()
		me.X += p.X
		me.Y += p.Y
		world.SetPlayer(me)
	}
}

func (me *Player) Vanish() {
	me.Vanished = true
}
func (me *Player) IsVanish() (bool) {
	return me.Vanished
}
func (me *Player) Src() (x0, y0, x1, y1 int) {
	idx := me.Anime.Index()
	if me.FaceDirection == funcs.FaceLeft {
		idx += 2
	}
	r := ImageSources[idx]
	return r.Left, r.Top, r.Right, r.Bottom
}
func (me *Player) Dst() (x0, y0, x1, y1 int) {
	x, y := int(me.X) + AdjustX, int(me.Y) + AdjustY
	return x, y, x + Width, y + Height
}
func (me *Player) SetInput(bits ginput.InputBits) {
	me.InputBits = bits
}
func (me *Player) SetPoint(pt fig.FloatPoint) {
	me.FloatPoint = pt
}

func (me *Player) HitRects() ([]fig.Rect) {
	x, y := int(me.X) + AdjustX, int(me.Y) + AdjustY
	return []fig.Rect{{x, y, x + Width, y + Height}}
}

func (me *Player) Hit(origin action.Object) {
	if me.Blackout || me.Beaten {
	} else {
		me.Beaten = true
		me.V.Jump(13)
		if origin.Point().X > me.X {
			me.V.Left.Accel(me.V.Left.MaxPower)
		} else {
			me.V.Right.Accel(me.V.Right.MaxPower)
		}
	}
/*
	if me.Beaten {
	} else {
		me.Beaten = true
		me.BlackoutTimer.Start(60)
		me.V.Jump(13)
		if origin.Point().X > me.X {
			me.V.Left.Accel(me.V.Left.MaxPower)
		} else {
			me.V.Right.Accel(me.V.Right.MaxPower)
		}
	}
	*/
}

func (me *Player) HitWall(origin action.Object) {
}

func (me *Player) Expel(hitWalls []fig.Rect) {
	pt, status := me.FallingRects.HitWall(me.FloatPoint.ToInt(), me.V.Power(), hitWalls)

	if (status & funcs.WallTop) != 0 {
		me.V.JumpCancel()
	}

	if (status & funcs.WallBottom) != 0 {
		me.CanJump = true
		me.V.JumpCancel()
	} else {
		me.CanJump = false
	}

	if (status & funcs.WallLeft) != 0 {
		me.V.Left.Reset()
	}

	if (status & funcs.WallRight) != 0 {
		me.V.Right.Reset()
	}

	me.FloatPoint = pt.ToFloat()
}

func (me *Player) Stack() (*script.Stack) {
	return nil
}

func New(pt fig.FloatPoint) (*Player) {
	hitWidth := 32
	hitHeight := 64
	hitAdjustX := -(hitWidth / 2)
	hitAdjustY := -64

	return &Player{
		FloatPoint: pt,
		V:          move.NewFallingInertia(radian.Right(), 0, 0.7, 7),
		Kbuffer:    &kcmd.Buffer{},
		Endurance:  100,
		Anime:      anime.NewFrames(8, 8),
		FallingRects: funcs.NewFallingRects(hitWidth, hitHeight, hitAdjustX, hitAdjustY),
		BlackoutTimer:  timer.NewFrame(0),
	}
}

//########################################
//# Objects
//########################################
type Interface interface {
	action.Object
	SetPoint(fig.FloatPoint)
	HitWall(action.Object)
	Expel([]fig.Rect)
	SetInput(ginput.InputBits)
}

type Objects struct {
	*sprites.Objects
}

func (me *Objects) Get(i int) (Interface) {
	return me.Objs[i].(Interface)
}

func (me *Objects) HitWall(i int, obj action.Object) {
	me.Get(i).HitWall(obj)
}

func (me *Objects) Expel(i int, hitWalls []fig.Rect) {
	me.Get(i).Expel(hitWalls)
}

func (me *Objects) SetInput(i int, bits ginput.InputBits) {
	me.Get(i).SetInput(bits)
}

func (me *Objects) SetPoint(i int, pt fig.FloatPoint) {
	me.Get(i).SetPoint(pt)
}

func NewObjects() (*Objects) {
	return &Objects {
		Objects: sprites.NewObjects(),
	}
}

