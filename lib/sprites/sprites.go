
package sprites

import (
	"github.com/ikuo0/game/lib/action"
	"github.com/ikuo0/game/lib/event"
	"github.com/ikuo0/game/lib/fig"
	"github.com/ikuo0/game/lib/orig"
	"github.com/hajimehoshi/ebiten"
)

//########################################
//# OneSprites
//########################################
type OneSprites struct {
	SX0, SY0, SX1, SY1 int
	DX0, DY0, DX1, DY1 int
}

func (me *OneSprites) Len() (int) {
	return 1
}

func (me *OneSprites) Dst(i int) (x0, y0, x1, y1 int) {
	return me.DX0, me.DY0, me.DX1, me.DY1
}

func (me *OneSprites) Src(i int) (x0, y0, x1, y1 int) {
	return me.SX0, me.SY0, me.SX1, me.SY1
}

func NewOneSprites(sx0, sy0, sx1, sy1, dx0, dy0, dx1, dy1 int) (*OneSprites) {
	res := OneSprites{
		sx0, sy0, sx1, sy1,
		dx0, dy0, dx1, dy1,
	}
	return &res;
}

//########################################
//# Objects
//########################################
type Objects struct {
	Objs []action.Object
}

func (me *Objects) GetObject(i int) (action.Object) {
	return me.Objs[i]
}

func (me *Objects) Len() (int) {
	return len(me.Objs)
}

func (me *Objects) Src(i int) (x0, y0, x1, y1 int) {
	return me.Objs[i].Src()
}

func (me *Objects) Dst(i int) (x0, y0, x1, y1 int) {
	return me.Objs[i].Dst()
}

func (me *Objects) HitRects(i int) ([]fig.Rect) {
	return me.Objs[i].HitRects()
}

func (me *Objects) Hit(i int, obj action.Object) {
	me.Objs[i].Hit(obj)
}

func (me *Objects) Update(i int, trigger event.Trigger) {
	me.Objs[i].Update(trigger)
}

func (me *Objects) Origin(i int) (orig.Interface) {
	return me.Objs[i]
}

func (me *Objects) Vanish(i int) {
	me.Objs[i].Vanish()
}

func (me *Objects) Clean(i int) {
	newObjs := []action.Object{}
	for _, v := range me.Objs {
		if !v.IsVanish() {
			newObjs = append(newObjs, v)
		}
	}
	me.Objs = newObjs
}

func (me *Objects) Options() (*ebiten.DrawImageOptions) {
	return &ebiten.DrawImageOptions {
		ImageParts: me,
	}
}

func (me *Objects) Occure(objIf action.Object) {
	me.Objs = append(me.Objs, objIf)
}

func NewObjects() (*Objects) {
	return &Objects {}
}

//########################################
//# RotaObjects
//########################################
type RotaObjects struct {
	*Objects
}

func (me *RotaObjects) DrawOption(i int) (*ebiten.DrawImageOptions) {
	sx0, sy0, sx1, sy1 := me.Src(i)
	dx0, dy0, dx1, dy1 := me.Dst(i)
	opt := ebiten.DrawImageOptions {
		ImageParts: NewOneSprites(sx0, sy0, sx1, sy1, dx0, dy0, dx1, dy1),
	}

	o := me.Objs[i]
	pt := o.GetPoint()
	opt.GeoM.Translate(float64(-pt.X), float64(-pt.Y))
	opt.GeoM.Rotate(float64(o.Direction()))
	opt.GeoM.Translate(float64(pt.X), float64(pt.Y))
	return &opt
}

func NewRotaObjects() (*RotaObjects) {
	return &RotaObjects {
		Objects: NewObjects(),
	}
}


//########################################
//# DrawHitRects
//########################################
type HitObjects struct {
	Objs []fig.Rect
}
func (me *HitObjects) Len() (int) {
	return len(me.Objs)
}
func (me *HitObjects) Src(i int) (x0, y0, x1, y1 int) {
	return 0, 0, 1, 1
}
func (me *HitObjects) Dst(i int) (x0, y0, x1, y1 int) {
	r := me.Objs[i]
	return int(r.Left), int(r.Top), int(r.Right), int(r.Bottom)
}
func (me *HitObjects) Options() (*ebiten.DrawImageOptions) {
	return &ebiten.DrawImageOptions {
		ImageParts: me,
	}
}
func NewHitObjects(who ...action.CanHit) (*HitObjects) {
	rects := []fig.Rect{}
	for _, x := range who {
		for i := 0; i < x.Len(); i++ {
			rects = append(rects, x.HitRects(i)...)
		}
	}
	return &HitObjects {
		Objs: rects,
	}
}
