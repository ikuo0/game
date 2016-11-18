
package shot

import (
	"github.com/ikuo0/game/lib/action"
	"github.com/ikuo0/game/lib/event"
	"github.com/ikuo0/game/lib/fig"
	"github.com/ikuo0/game/lib/gradian"
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
var SrcShot = fig.IntRect {0, 66, 0 + 60, 66 + 66}
type Shot struct {
	fig.Point
	Vanished   bool
	V          *move.FixedVector
	Endurance  int
}

func (me *Shot) GetPoint() (fig.Point) {
	return me.Point
}

func (me *Shot) Direction() (radian.Radian) {
	return gradian.Up()
}

func (me *Shot) Update(trigger event.Trigger) {
	me.V.Accel(32)
	me.X += me.V.X()
	me.Y += me.V.Y()
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
func (me *Shot) SetPoint(pt fig.Point) {
	me.Point = pt
}
func (me *Shot) HitRects() ([]fig.Rect) {
	if me.Endurance <= 0 {
		return nil
	} else {
		x, y := me.X - 24, me.Y - 30
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

func NewShot(pt fig.Point) (*Shot) {
	return &Shot{
		Point: pt,
		V:          move.NewFixedVector(gradian.Up(), 32),
		Endurance:  6,
	}
}
