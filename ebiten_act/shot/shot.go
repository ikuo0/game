
package shot

import (
	"github.com/ikuo0/game/ebiten_act/action"
	"github.com/ikuo0/game/lib/event"
	"github.com/ikuo0/game/lib/script"
	"github.com/ikuo0/game/lib/fig"
	"github.com/ikuo0/game/lib/move"
	"github.com/ikuo0/game/lib/radian"
)

//########################################
//# Shot
//########################################
const Width = 20
const Height = 20
const AdjustX = -10
const AdjustY = -10
var ImageSrc = fig.IntRect {0, 0, Width, Height}
type Shot struct {
	action.Object
	V          *move.FixedVector
	Endurance  int
}

func (me *Shot) Update(trigger event.Trigger) {
	me.V.Accel(16)
	me.X += me.V.X()
	me.Y += me.V.Y()
}

func (me *Shot) Src() (x0, y0, x1, y1 int) {
	return ImageSrc.Left, ImageSrc.Top, ImageSrc.Right, ImageSrc.Bottom
}
func (me *Shot) Dst() (x0, y0, x1, y1 int) {
	x, y := int(me.X) + AdjustX, int(me.Y) + AdjustY
	return x, y, x + Width, y + Height
}
func (me *Shot) HitRects() ([]fig.Rect) {
	if me.Endurance <= 0 {
		return nil
	} else {
		x, y := me.X + AdjustX, me.Y + AdjustY
		return []fig.Rect{{x, y, x + Width, y + Height}}
	}
}

func (me *Shot) Hit(obj action.Interface) {
	me.Endurance--
	me.Vanish()
}

func (me *Shot) Stack() (*script.Stack) {
	return nil
}

func New(pt fig.Point, rad radian.Radian) (*Shot) {
	return &Shot{
		Object: action.Object {
			Point: pt,
			Radian:     rad,
		},
		V:          move.NewFixedVector(rad, 16),
		Endurance:  1,
	}
}

