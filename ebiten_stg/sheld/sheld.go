
package sheld

import (
	"github.com/ikuo0/game/lib/action"
	"github.com/ikuo0/game/lib/anime"
	"github.com/ikuo0/game/lib/event"
	"github.com/ikuo0/game/lib/fig"
	"github.com/ikuo0/game/lib/gradian"
	"github.com/ikuo0/game/lib/move"
	"github.com/ikuo0/game/lib/radian"
	"github.com/ikuo0/game/lib/script"
	"github.com/ikuo0/game/lib/sprites"
	//"github.com/hajimehoshi/ebiten"
	//"math"
	//"fmt"
)

//########################################
//# Shield
//########################################
var SrcSheld = fig.Rect {64, 64, 64 + 320, 64 + 320}
type Sheld struct {
	fig.FloatPoint
	Vanished   bool
	Anime      *anime.Frames
	V          *move.FixedVector
}

func (me *Sheld) Point() (fig.FloatPoint) {
	return me.FloatPoint
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
func (me *Sheld) SetPoint(pt fig.FloatPoint) {
	me.FloatPoint = pt
}
func (me *Sheld) HitRects() ([]fig.Rect) {
	if me.IsVanish() {
		return nil
	} else {
		x, y := int(me.X) - 48, int(me.Y) - 48
		return []fig.Rect{{x, y, x + 96, y + 96}}
	}
}

func (me *Sheld) Hit(obj action.Object) {
}

func (me *Sheld) Pushed() {
	me.V.Accel(0.2)
}

func (me *Sheld) Stack() (*script.Stack) {
	return nil
}

func NewSheld(pt fig.FloatPoint) (*Sheld) {
	v := move.NewFixedVector(gradian.Up(), 0.7)
	v.Rate = 0.2
	return &Sheld{
		FloatPoint: pt,
		Anime:      anime.NewFrames(8, 7, 3, 2, 8),
		V:          v,
	}
}

//########################################
//# Objects
//########################################
type Interface interface {
	action.Object
	Pushed()
}

type Objects struct {
	*sprites.Objects
}
func (me *Objects) Get(i int) (Interface) {
	return me.Objs[i].(Interface)
}
func (me *Objects) Pushed(i int) {
	me.Get(i).Pushed()
}

func NewObjects() (*Objects) {
	return &Objects {
		Objects: sprites.NewObjects(),
	}
}
