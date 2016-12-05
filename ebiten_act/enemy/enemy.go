package enemy

// 漢字

import (
	"github.com/ikuo0/game/ebiten_act/action"
	"github.com/ikuo0/game/ebiten_act/eventid"
	"github.com/ikuo0/game/ebiten_act/funcs"
	"github.com/ikuo0/game/lib/anime"
	"github.com/ikuo0/game/lib/event"
	"github.com/ikuo0/game/lib/fig"
	"github.com/ikuo0/game/lib/move"
	"github.com/ikuo0/game/lib/script"
	"github.com/ikuo0/game/lib/timer"
	//"fmt"
)

const Width = 64
const Height = 43
const AdjustX = -32
const AdjustY = -43

//128, 43
var ImageSources []fig.IntRect = []fig.IntRect {
	{0, 0, 64, 43},
	{64, 0, 64 + 64, 43},
	{64, 0, 0, 43},
	{64 + 64, 0, 64, 43},
}

type Enemy struct {
	action.Object
	Config        funcs.EnemyConfig
	Ready         bool
	ReadyTimer    *timer.Frame
	FaceDirection funcs.FaceDirection
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
		trigger.EventTrigger(eventid.Explosion1, me.Point, nil)
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
	x, y := me.X + AdjustX, me.Y + AdjustY
	return []fig.Rect{{x, y, x + Width, y + Height}}
}

func (me *Enemy) Hit(obj action.Interface) {
	me.Endurance--
	if me.Endurance <= 0 {
		me.Dead = true
	}
}

func (me *Enemy) HitWall(origin action.Interface) {
}

func (me *Enemy) Expel(hitWalls []fig.Rect) {
	pt, status := me.FallingRects.HitFloor(me.Point, true, hitWalls)

	if (status & funcs.FloorTop) != 0 {
		me.Gravity.JumpCancel()
	}

	if (status & funcs.FloorBottom) != 0 {
		me.CanJump = true
		me.Gravity.Landing()
	} else {
		me.CanJump = false
	}

	if (status & funcs.FloorLeft) != 0 {
		//me.Xinertia.Backward.Reset()
		me.FacingRight()
	}

	if (status & funcs.FloorRight) != 0 {
		//me.Xinertia.Advance.Reset()
		me.FacingLeft()
	}

	if !me.Ready {
		me.Point.Y -= 0.5
	} else {
		me.Point = pt
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
		Object: action.Object {
			Point: config.Point,
		},
		Config:        config,
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

