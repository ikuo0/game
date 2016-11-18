
package wall

import (
	"github.com/ikuo0/game/ebiten_race/eventid"
	"github.com/ikuo0/game/lib/action"
	"github.com/ikuo0/game/lib/event"
	"github.com/ikuo0/game/lib/fig"
	"github.com/ikuo0/game/lib/radian"
)

//########################################
//# Wall
//########################################
const Width = 24
const Height = 24
const AdjustX = -12
const AdjustY = -12
var ImageSrc = fig.IntRect {0, 0, Width, Height}

type Wall struct {
	fig.Point
	Hitme bool
}

func (me *Wall) GetPoint() (fig.Point) {
	return me.Point
}

func (me *Wall) Direction() (radian.Radian) {
	return 0
}

func (me *Wall) Update(trigger event.Trigger) {
	if me.Hitme {
		trigger.EventTrigger(eventid.CollisionWall, me.Point, me)
		me.Hitme = false
	}
}

func (me *Wall) Vanish() {
}
func (me *Wall) IsVanish() (bool) {
	return false
}
func (me *Wall) Src() (x0, y0, x1, y1 int) {
	return ImageSrc.Left, ImageSrc.Top, ImageSrc.Right, ImageSrc.Bottom
}
func (me *Wall) Dst() (x0, y0, x1, y1 int) {
	x, y := int(me.X) + AdjustX, int(me.Y) + AdjustY
	return x, y, x + Width, y + Height
}
func (me *Wall) HitRects() ([]fig.Rect) {
	x, y := me.X + AdjustX, me.Y + AdjustY
	return []fig.Rect{{x, y, x + Width, y + Height}}
}

func (me *Wall) Hit(obj action.Object) {
	me.Hitme = true
}

func New(pt fig.Point) (*Wall) {
	return &Wall {
		Point: pt,
	}
}

