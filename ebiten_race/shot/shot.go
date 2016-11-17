
package shot

import (
	"github.com/ikuo0/game/lib/action"
	"github.com/ikuo0/game/lib/event"
	"github.com/ikuo0/game/lib/fig"
	"github.com/ikuo0/game/lib/move"
	"github.com/ikuo0/game/lib/radian"
	"github.com/ikuo0/game/lib/script"
	//"github.com/hajimehoshi/ebiten"
	//"math"
	//"fmt"
)

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

func (me *Shot) Hit(obj action.Object) {
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
