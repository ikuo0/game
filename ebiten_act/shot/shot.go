
package shot

import (
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
var ImageSrc = fig.Rect {0, 0, Width, Height}
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
	return ImageSrc.Left, ImageSrc.Top, ImageSrc.Right, ImageSrc.Bottom
}
func (me *Shot) Dst() (x0, y0, x1, y1 int) {
	x, y := int(me.X) + AdjustX, int(me.Y) + AdjustY
	return x, y, x + Width, y + Height
}
func (me *Shot) SetPoint(pt fig.FloatPoint) {
	me.FloatPoint = pt
}
func (me *Shot) HitRects() ([]fig.Rect) {
	if me.Endurance <= 0 {
		return nil
	} else {
		x, y := int(me.X) + AdjustX, int(me.Y) + AdjustY
		return []fig.Rect{{x, y, x + Width, y + Height}}
	}
}

func (me *Shot) Hit() {
	me.Endurance--
	me.Vanish()
}

func (me *Shot) Stack() (*script.Stack) {
	return nil
}

func New(pt fig.FloatPoint, rad radian.Radian) (*Shot) {
	return &Shot{
		FloatPoint: pt,
		V:          move.NewConstant(rad, 16),
		Endurance:  1,
	}
}

