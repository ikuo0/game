
package wall

import (
	"github.com/ikuo0/game/ebiten_race/action"
	"github.com/ikuo0/game/ebiten_race/eventid"
	"github.com/ikuo0/game/lib/event"
	"github.com/ikuo0/game/lib/fig"
	"github.com/ikuo0/game/lib/radian"
)

//########################################
//# Wall
//########################################
const Width = 32
const Height = 32
const AdjustX = -16
const AdjustY = -16

type Kind int
const (
	Right       Kind = 1
	RightTop    Kind = 2
	Top         Kind = 3
	LeftTop     Kind = 4
	Left        Kind = 5
	LeftBottom  Kind = 6
	Bottom      Kind = 7
	RightBottom Kind = 8
	Space       Kind = 9
)

func (me Kind) Src() (int, int, int, int) {
	switch(me) {
		case Right:       return 32, 0, 63, 32
		case RightTop:    return 64, 0, 95, 32
		case Top:         return 96, 0, 127, 32
		case LeftTop:     return 128, 0, 159, 32
		case Left:        return 160, 0, 191, 32
		case LeftBottom:  return 192, 0, 223, 32
		case Bottom:      return 224, 0, 255, 32
		case RightBottom: return 256, 0, 287, 32
	}
	return 0, 0, 31, 32
}

func (me Kind) Line() (fig.Line) {
	switch(me) {
		case Right:       return fig.Line{fig.Point{0, 0}, fig.Point{0, Height}}
		case RightTop:    return fig.Line{fig.Point{0, 0}, fig.Point{Width, Height}}
		case Top:         return fig.Line{fig.Point{0, Height}, fig.Point{Width, Height}}
		case LeftTop:     return fig.Line{fig.Point{0, Height}, fig.Point{Width, 0}}
		case Left:        return fig.Line{fig.Point{Width, 0}, fig.Point{Width, Height}}
		case LeftBottom:  return fig.Line{fig.Point{0, 0}, fig.Point{Width, Height}}
		case Bottom:      return fig.Line{fig.Point{0, 0}, fig.Point{Width, 0}}
		case RightBottom: return fig.Line{fig.Point{0, Height}, fig.Point{Width, 0}}
	}
	return fig.Line{fig.Point{0, 0}, fig.Point{0, 0}}
}

type Parameter struct {
	fig.Point
	Kind Kind
}

type Wall struct {
	action.Object
	Kind Kind
	Hitme bool
}

func (me *Wall) GetPoint() (fig.Point) {
	return me.Point
}

func (me *Wall) Direction() (radian.Radian) {
	return 0
}

func (me *Wall) Line() (fig.Line) {
	return me.Kind.Line().Relative(me.Point)
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
	return me.Kind.Src()
}
func (me *Wall) Dst() (x0, y0, x1, y1 int) {
	x, y := int(me.X) + AdjustX, int(me.Y) + AdjustY
	return x, y, x + Width, y + Height
}
func (me *Wall) HitRects() ([]fig.Rect) {
	x, y := me.X + AdjustX, me.Y + AdjustY
	return []fig.Rect{{x, y, x + Width, y + Height}}
}

func (me *Wall) Hit(obj action.Interface) {
	me.Hitme = true
}

func New(pa Parameter) (*Wall) {
	return &Wall {
		Object: action.Object {
			Point: pa.Point,
		},
		Kind:  pa.Kind,
	}
}

