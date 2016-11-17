
package block

import (
	"github.com/ikuo0/game/ebiten_act/eventid"
	"github.com/ikuo0/game/ebiten_act/funcs"
	"github.com/ikuo0/game/lib/action"
	"github.com/ikuo0/game/lib/event"
	"github.com/ikuo0/game/lib/script"
	"github.com/ikuo0/game/lib/fig"
	"github.com/ikuo0/game/lib/radian"
	"github.com/ikuo0/game/lib/timer"
	//"fmt"
	"math/rand"
)

//########################################
//# Block
//########################################
var ImageSource = []fig.Rect {
	{
		0,
		0,
		32,
		32,
	},
}

type Block struct {
	fig.FloatPoint
}

func (me *Block) Point() (fig.FloatPoint) {
	return me.FloatPoint
}

func (me *Block) Direction() (radian.Radian) {
	return 0
}

func (me *Block) Update(trigger event.Trigger) {
}

func (me *Block) Vanish() {
}
func (me *Block) IsVanish() (bool) {
	return false
}
func (me *Block) Src() (x0, y0, x1, y1 int) {
	x := ImageSource[0]
	return x.Left, x.Top, x.Right, x.Bottom
}
func (me *Block) Dst() (x0, y0, x1, y1 int) {
	x, y := int(me.X) - 16, int(me.Y) - 16
	return x, y, x + 32, y + 32
}
func (me *Block) HitRects() ([]fig.Rect) {
	x, y := int(me.X) - 16, int(me.Y) - 16
	return []fig.Rect{{x, y, x + 32, y + 32}}
}

func (me *Block) Hit(obj action.Object) {
}

func (me *Block) Stack() (*script.Stack) {
	return nil
}

func NewBlock(pt fig.FloatPoint) (*Block) {
	return &Block {
		FloatPoint: pt,
	}
}

//########################################
//# OccureBlock
//########################################
type OccureDirection int
const (
	OccureLeft      OccureDirection = iota + 1
	OccureRight
	OccureRand
)

type Config struct {
	Point           fig.FloatPoint
	Span            int
	OccureDirection OccureDirection
}

type OccureBlock struct {
	*Block
	Config Config
	Timer   timer.Frame
}

func (me *OccureBlock) Update(trigger event.Trigger) {
	if me.Timer.Up() {
		defer func() {
			me.Timer.Start(me.Config.Span)
		}()

		d := funcs.FaceDirection(0)
		if me.Config.OccureDirection == OccureLeft {
			d = funcs.FaceLeft
		} else if me.Config.OccureDirection == OccureRight {
			d = funcs.FaceRight
		} else {
			if rand.Intn(2) == 0 {
				d = funcs.FaceLeft
			} else {
				d = funcs.FaceRight
			}
		}
		setting := funcs.EnemyConfig {
			Point:     me.FloatPoint,
			Direction: d,
		}
		trigger.EventTrigger(eventid.Enemy, setting, me)
	}
}

func NewOccureBlock(config Config) (*OccureBlock) {
	return &OccureBlock {
		Block:   NewBlock(config.Point),
		Config:  config,
		//Timer:   *timer.NewFrame(config.Span),
		Timer:   *timer.NewFrame(0),
	}
}
