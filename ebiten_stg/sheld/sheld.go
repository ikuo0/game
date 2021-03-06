
package sheld

import (
	"github.com/ikuo0/game/ebiten_stg/action"
	"github.com/ikuo0/game/lib/anime"
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
//# Shield
//########################################
var SrcSheld = fig.IntRect {64, 64, 64 + 320, 64 + 320}
type Sheld struct {
	action.Object
	fig.Point
	Vanished   bool
	Anime      *anime.Frames
	V          *move.FixedVector
}

func (me *Sheld) GetPoint() (fig.Point) {
	return me.Point
}

func (me *Sheld) Direction() (radian.Radian) {
	return gradian.Up()
}

func (me *Sheld) Update(trigger event.Trigger) {
	//me.V.Accel(0.2)
	me.X += me.V.X()
	me.Y += me.V.Y()
	me.Anime.Update()
}

func (me *Sheld) Vanish() {
	me.Vanished = true
}
func (me *Sheld) IsVanish() (bool) {
	return me.Vanished
}
func (me *Sheld) Src() (x0, y0, x1, y1 int) {
	return SrcSheld.Left, SrcSheld.Top, SrcSheld.Right, SrcSheld.Bottom
}
func (me *Sheld) Dst() (x0, y0, x1, y1 int) {
// 96
// 16 + 80
	width := me.Anime.Index() * 15 + 36
	adjust := width / 2
	x, y := int(me.X) - adjust, int(me.Y) - adjust
	return x, y, x + width, y + width
}
func (me *Sheld) SetPoint(pt fig.Point) {
	me.Point = pt
}
func (me *Sheld) HitRects() ([]fig.Rect) {
	if me.IsVanish() {
		return nil
	} else {
		x, y := me.X - 48, me.Y - 48
		return []fig.Rect{{x, y, x + 96, y + 96}}
	}
}

func (me *Sheld) Hit(obj action.Interface) {
}

func (me *Sheld) Pushed() {
	me.V.Accel(0.2)
}

func (me *Sheld) Stack() (*script.Stack) {
	return nil
}

func NewSheld(pt fig.Point) (*Sheld) {
	v := move.NewFixedVector(gradian.Up(), 0.7)
	v.Rate = 0.2
	return &Sheld{
		Point: pt,
		Anime:      anime.NewFrames(8, 7, 3, 2, 8),
		V:          v,
	}
}
