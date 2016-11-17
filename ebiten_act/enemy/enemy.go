package enemy

// 漢字

import (
	"github.com/ikuo0/game/ebiten_act/eventid"
	"github.com/ikuo0/game/ebiten_act/funcs"
	"github.com/ikuo0/game/lib/action"
	"github.com/ikuo0/game/lib/anime"
	"github.com/ikuo0/game/lib/event"
	"github.com/ikuo0/game/lib/fig"
	"github.com/ikuo0/game/lib/move"
	"github.com/ikuo0/game/lib/radian"
	"github.com/ikuo0/game/lib/script"
	"github.com/ikuo0/game/lib/sprites"
	"github.com/ikuo0/game/lib/timer"
	//"fmt"
)

const Width = 64
const Height = 43
const AdjustX = -32
const AdjustY = -43

//128, 43
var ImageSources []fig.Rect = []fig.Rect {
	{0, 0, 64, 43},
	{64, 0, 64 + 64, 43},
	{64, 0, 0, 43},
	{64 + 64, 0, 64, 43},
}

type Enemy struct {
	fig.FloatPoint
	Config        funcs.EnemyConfig
	Ready         bool
	ReadyTimer    *timer.Frame
	FaceDirection funcs.FaceDirection
	Vanished      bool
	Gravity       *funcs.Gravity
	V             *move.Vector
	Xinertia      *move.Inertia
	Endurance     int
	Dead          bool
	Anime         *anime.Frames
	FallingRects  *funcs.FallingRects
	CanJump       bool
	DeadTimer     *timer.Frame
}

func (me *Enemy) Point() (fig.FloatPoint) {
	return me.FloatPoint
}

func (me *Enemy) Direction() (radian.Radian) {
	return 0
}

func (me *Enemy) FacingLeft() {
	me.V.Degree.Deg = 180
	me.FaceDirection = funcs.FaceLeft
}

func (me *Enemy) FacingRight() {
	me.V.Degree.Deg = 0
	me.FaceDirection = funcs.FaceRight
}

func (me *Enemy) Update(trigger event.Trigger) {
	if me.Dead {
		trigger.EventTrigger(eventid.Explosion1, me.FloatPoint, nil)
		me.Vanish()
	} else if !me.Ready {
		if me.ReadyTimer.Up() {
			me.Ready = true
		}
	} else {
		if me.DeadTimer.Up() {
			me.Dead = true
		}

		if me.CanJump {
			me.V.Accel(0.4)
			me.Xinertia.Set(me.V.X())
		}

		me.Anime.Update()
		me.Gravity.Update()
		me.Xinertia.Update()

		me.X += me.Xinertia.Value()
		me.Y += me.Gravity.Value()
	}
}

func (me *Enemy) Vanish() {
	me.Vanished = true
}
func (me *Enemy) IsVanish() (bool) {
	return me.Vanished
}
func (me *Enemy) Src() (x0, y0, x1, y1 int) {
	idx := me.Anime.Index()
	if me.FaceDirection == funcs.FaceLeft {
		idx += 2
	}
	r := ImageSources[idx]
	return r.Left, r.Top, r.Right, r.Bottom
}
func (me *Enemy) Dst() (x0, y0, x1, y1 int) {
	x, y := int(me.X) + AdjustX, int(me.Y) + AdjustY
	return x, y, x + Width, y + Height
}
func (me *Enemy) HitRects() ([]fig.Rect) {
	x, y := int(me.X) + AdjustX, int(me.Y) + AdjustY
	return []fig.Rect{{x, y, x + Width, y + Height}}
}

func (me *Enemy) Hit(origin action.Object) {
	me.Endurance--
	if me.Endurance <= 0 {
		me.Dead = true
	}
}

func (me *Enemy) HitWall(origin action.Object) {
}

func (me *Enemy) Expel(hitWalls []fig.Rect) {
	pt, status := me.FallingRects.HitWall(me.FloatPoint, true, hitWalls)

	if (status & funcs.WallTop) != 0 {
		me.Gravity.JumpCancel()
	}

	if (status & funcs.WallBottom) != 0 {
		me.CanJump = true
		me.Gravity.Landing()
	} else {
		me.CanJump = false
	}

	if (status & funcs.WallLeft) != 0 {
		//me.Xinertia.Backward.Reset()
		me.FacingRight()
	}

	if (status & funcs.WallRight) != 0 {
		//me.Xinertia.Advance.Reset()
		me.FacingLeft()
	}

	if !me.Ready {
		me.FloatPoint.Y -= 0.5
	} else {
		me.FloatPoint = pt.ToFloat()
	}
}

func (me *Enemy) Stack() (*script.Stack) {
	return nil
}

func New(config funcs.EnemyConfig) (*Enemy) {
	deg := int(0)
	fd := funcs.FaceDirection(0)
	if config.Direction == funcs.FaceLeft {
		deg = 180
		fd = funcs.FaceLeft
	} else {
		deg = 0
		fd = funcs.FaceRight
	}

	return &Enemy{
		Config:        config,
		FloatPoint:    config.Point,
		Gravity:       funcs.NewGravity(),
		V:             move.NewVector(deg, 5),
		Xinertia:      move.NewInertia(0.4),
		FaceDirection: fd,
		Endurance:     1,
		Anime:         anime.NewFrames(7, 7),
		FallingRects:  funcs.NewFallingRects(Width, Height, AdjustX, AdjustY),
		DeadTimer:     timer.NewFrame(1200),
		ReadyTimer:    timer.NewFrame(32),
	}
}

//########################################
//# Objects
//########################################
type Interface interface {
	action.Object
	HitWall(action.Object)
	Expel([]fig.Rect)
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

func NewObjects() (*Objects) {
	return &Objects {
		Objects: sprites.NewObjects(),
	}
}

