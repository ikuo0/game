
package enemy

import (
	"github.com/ikuo0/game/ebiten_race/eventid"
	"github.com/ikuo0/game/ebiten_race/effect"
	"github.com/ikuo0/game/ebiten_race/world"
	"github.com/ikuo0/game/lib/action"
	"github.com/ikuo0/game/lib/event"
	"github.com/ikuo0/game/lib/script"
	"github.com/ikuo0/game/lib/anime"
	"github.com/ikuo0/game/lib/fig"
	"github.com/ikuo0/game/lib/move"
	"github.com/ikuo0/game/lib/radian"
	"github.com/ikuo0/game/lib/timer"
	"math"
	"math/rand"
)

//########################################
//# Heli0
//########################################
var heli0Source = []fig.Rect {
	{
		0,
		64,
		0 + 48,
		64 + 48,
	},
	{
		0,
		128,
		0 + 48,
		128 + 48,
	},
}

type Heli0 struct {
	fig.FloatPoint
	V         *move.Inertia
	Anime     *anime.Frames
	Vanished  bool
	Endurance int
	MyStack   script.Stack
	Timer     *timer.Frame
}

func (me *Heli0) Point() (fig.FloatPoint) {
	return me.FloatPoint
}

func (me *Heli0) Direction() (radian.Radian) {
	return me.V.Radian
}

func (me *Heli0) Update(trigger event.Trigger) {
	if me.Endurance < 1 {
		me.Vanish()
		trigger.EventTrigger(eventid.Explosion1, nil, me)
	} else {
		me.V.Accel()
		p := me.V.Power()
		me.X += p.X
		me.Y += p.Y
		me.Anime.Update()
		me.V.Radian = me.V.Radian.TurnLeft(1)
	}
}

func (me *Heli0) SuperUpdate(trigger event.Trigger) (bool) {
	if me.Endurance < 1 {
		me.Vanish()
		trigger.EventTrigger(eventid.Explosion1, nil, me)
		trigger.EventTrigger(eventid.Score, int(10), me)
		return false
	} else {
		if me.Timer.Up() {
			aim := world.GetPlayer().Point()
			aimRad := radian.Radian(math.Atan2(aim.Y - me.Y, aim.X - me.X))
			trigger.EventTrigger(eventid.Bullet2, aimRad, me)
			me.Timer.Start(10000)
		}
		return true
	}
}

func (me *Heli0) Vanish() {
	me.Vanished = true
}
func (me *Heli0) IsVanish() (bool) {
	return me.Vanished
}
func (me *Heli0) Src() (x0, y0, x1, y1 int) {
	x := heli0Source[me.Anime.Index()]
	return x.Left, x.Top, x.Right, x.Bottom
}
func (me *Heli0) Dst() (x0, y0, x1, y1 int) {
	x, y := int(me.X) - 24, int(me.Y) - 24
	return x, y, x + 48, y + 48
}
func (me *Heli0) HitRects() ([]fig.Rect) {
	x, y := int(me.X) - 24, int(me.Y) - 24
	return []fig.Rect{{x, y, x + 48, y + 48}}
}

func (me *Heli0) Hit(obj action.Object) {
	me.Endurance--
}
func (me *Heli0) Stack() (*script.Stack) {
	return &me.MyStack
}

func NewHeli0(pt fig.FloatPoint) (*Heli0) {
	return &Heli0 {
		FloatPoint: pt,
		V:          move.NewInertia(radian.Down(), 1, 0.4, 8),
		Anime:      anime.NewFrames(4, 4),
		Endurance:  10,
	}
}

//########################################
//# Heli1
//########################################
type Heli1 struct {
	Heli0
}

func (me *Heli1) Update(trigger event.Trigger) {
	if me.SuperUpdate(trigger) {
		aim := world.GetPlayer().Point()
		aimRad := radian.Radian(math.Atan2(me.Y - aim.Y, me.X - aim.X))

		a := me.V.Radian - aimRad
		lr := a <= math.Pi && a >= -math.Pi
		if a < 0 {
			lr = !lr
		}

		if lr {
			me.V.Radian = me.V.Radian.TurnRight(1)
		} else {
			me.V.Radian = me.V.Radian.TurnLeft(1)
		}

		me.V.Accel()
		p := me.V.Power()
		me.X += p.X
		me.Y += p.Y
		me.Anime.Update()
	}
}

func NewHeli1(pt fig.FloatPoint) (*Heli1) {
	return &Heli1 {
		Heli0: Heli0 {
			FloatPoint: pt,
			V:          move.NewInertia(radian.Down(), 1, 0.4, 8),
			Anime:      anime.NewFrames(4, 4),
			Endurance:  1,
			Timer:      timer.NewFrame(rand.Intn(30) + 15),
		},
	}
}

//########################################
//# Heli2
//########################################
type Heli2 struct {
	Heli0
}

func (me *Heli2) Update(trigger event.Trigger) {
	if me.SuperUpdate(trigger) {
		me.V.Accel()
		p := me.V.Power()
		me.X += p.X
		me.Y += p.Y
		me.Anime.Update()
	}
}

func NewHeli2(pt fig.FloatPoint) (*Heli2) {
	return &Heli2 {
		Heli0: Heli0 {
			FloatPoint: pt,
			V:          move.NewInertia(radian.Down(), 0, 0.3, 8),
			Anime:      anime.NewFrames(4, 4),
			Endurance:  2,
			Timer:      timer.NewFrame(rand.Intn(30) + 15),
		},
	}
}

//########################################
//# Aide
//########################################
var AideScript = script.NewSource([]script.Proc {
	script.NewEventProc(eventid.Bullet1, radian.FromDeg(-8)),
	script.NewEventProc(eventid.Bullet1, radian.FromDeg(0)),
	script.NewEventProc(eventid.Bullet1, radian.FromDeg(8)),
	script.NewWaitProc(10),
	script.NewEventProc(eventid.Bullet1, radian.FromDeg(-12)),
	script.NewEventProc(eventid.Bullet1, radian.FromDeg(-4)),
	script.NewEventProc(eventid.Bullet1, radian.FromDeg(4)),
	script.NewEventProc(eventid.Bullet1, radian.FromDeg(12)),
	script.NewWaitProc(120),
	script.NewJumpProc(0),
})

var SrcAide = []fig.Rect {
	{
		64,
		0,
		64 + 64,
		0 + 64,
	},
}

type Aide struct {
	fig.FloatPoint
	MyStack      script.Stack
	Radian       radian.Radian
	Vanished     bool
	Dead         bool
	Endurance    int
	ReamExplosion *effect.ReamExplosion
}

func (me *Aide) Point() (fig.FloatPoint) {
	return me.FloatPoint
}

func (me *Aide) Direction() (radian.Radian) {
	return me.Radian
}

func (me *Aide) Update(trigger event.Trigger) {
	if me.Dead {
		if me.ReamExplosion.Update(trigger) {
			me.Vanish()
			trigger.EventTrigger(eventid.Score, int(1000), me)
		}
	} else {
		script.Exec(AideScript, &me.MyStack, me, trigger)

		aim := world.GetPlayer().Point()
		aimRad := radian.Radian(math.Atan2(me.Y - aim.Y, me.X - aim.X))

		a := me.Radian - aimRad
		lr := a <= math.Pi && a >= -math.Pi
		if a < 0 {
			lr = !lr
		}

		if lr {
			me.Radian = me.Radian.TurnRight(1)
		} else {
			me.Radian = me.Radian.TurnLeft(1)
		}
	}
}

func (me *Aide) Vanish() {
	me.Vanished = true
}
func (me *Aide) IsVanish() (bool) {
	return me.Vanished
}
func (me *Aide) Src() (x0, y0, x1, y1 int) {
	x := SrcAide[0]
	return x.Left, x.Top, x.Right, x.Bottom
}
func (me *Aide) Dst() (x0, y0, x1, y1 int) {
	x, y := int(me.X) - 32, int(me.Y) - 32
	return x, y, x + 64, y + 64
}
func (me *Aide) HitRects() ([]fig.Rect) {
	if me.Dead {
		return nil
	} else {
		x, y := int(me.X) - 32, int(me.Y) - 32
		return []fig.Rect{{x, y, x + 48, y + 60}}
	}
}

func (me *Aide) Hit(obj action.Object) {
	me.Endurance--
	if me.Endurance <= 0 {
		me.Dead = true
		me.ReamExplosion = effect.NewReamExplosion(64, 60, me.FloatPoint)
	}
}
func (me *Aide) Stack() (*script.Stack) {
	return &me.MyStack
}
func NewAide(pt fig.FloatPoint) (*Aide) {
	return &Aide {
		FloatPoint: pt,
		Endurance:  40,
	}
}

//########################################
//# Boss1
//########################################
var SrcBoss1 = []fig.Rect {
	{
		1,
		605,
		1 + 180,
		605 + 190,
	},
}

type Boss1 struct {
	*Aide
}

func (me *Boss1) Update(trigger event.Trigger) {
	me.Aide.Update(trigger)
	if me.IsVanish() {
		trigger.EventTrigger(eventid.StageClear, nil, nil)
	}
/*
	if me.Dead {
		if me.ReamExplosion.Update(trigger) {
			me.Vanish()
			trigger.EventTrigger(event.Score, int(1000), me)
		}
	} else {
		script.Exec(AideScript, &me.MyStack, me, trigger)

		aim := world.GetPlayer().Point()
		aimRad := radian.Radian(math.Atan2(me.Y - aim.Y, me.X - aim.X))

		a := me.Radian - aimRad
		lr := a <= math.Pi && a >= -math.Pi
		if a < 0 {
			lr = !lr
		}

		if lr {
			me.Radian = me.Radian.TurnRight(1)
		} else {
			me.Radian = me.Radian.TurnLeft(1)
		}
	}
	*/
}

func (me *Boss1) Src() (x0, y0, x1, y1 int) {
	x := SrcBoss1[0]
	return x.Left, x.Top, x.Right, x.Bottom
}
func (me *Boss1) Dst() (x0, y0, x1, y1 int) {
	x, y := int(me.X) - 90, int(me.Y) - 95
	return x, y, x + 180, y + 190
}
func (me *Boss1) HitRects() ([]fig.Rect) {
	if me.Dead {
		return nil
	} else {
		x, y := int(me.X) - 100, int(me.Y) - 100
		return []fig.Rect{{x, y, x + 200, y + 200}}
	}
}

func NewBoss1(pt fig.FloatPoint) (*Boss1) {
	return &Boss1 {
		Aide: &Aide {
			FloatPoint: pt,
			Endurance:  400,
		},
	}
}
