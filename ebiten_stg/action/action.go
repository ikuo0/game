
package action

import (
	"github.com/ikuo0/game/lib/event"
	"github.com/ikuo0/game/lib/fig"
	"github.com/ikuo0/game/lib/ginput"
	"github.com/ikuo0/game/lib/radian"
	"github.com/hajimehoshi/ebiten"
)


//########################################
//# Interface
//########################################
type Interface interface {
	GetPoint() (fig.Point)
	SetPoint(fig.Point)
	Direction() (radian.Radian)
	Update(trigger event.Trigger)
	Vanish()
	IsVanish() (bool)
	Src() (x0, y0, x1, y1 int)
	Dst() (x0, y0, x1, y1 int)
	SetInput(ginput.InputBits)
	HitRects() ([]fig.Rect)
	Hit(Interface)
	Pushed()
}

//########################################
//# Object
//########################################
type Object struct {
	fig.Point
	Radian    radian.Radian
	Vanished  bool
}

func (me *Object) SetInput(bits ginput.InputBits) {
}

func (me *Object) SetPoint(pt fig.Point) () {
	me.Point = pt
}

func (me *Object) GetPoint() (fig.Point) {
	return me.Point
}

func (me *Object) Direction() (radian.Radian) {
	return me.Radian
}

func (me *Object) Update(trigger event.Trigger) {
}

func (me *Object) Vanish() {
	me.Vanished = true
}

func (me *Object) IsVanish() (bool) {
	return me.Vanished
}

func (me *Object) Src() (x0, y0, x1, y1 int) {
	return 0, 0, 0, 0
}

func (me *Object) Dst() (x0, y0, x1, y1 int) {
	return 0, 0, 0, 0
}

func (me *Object) HitRects() ([]fig.Rect) {
	return nil
}

func (me *Object) Hit(Interface) {
}

func (me *Object) Pushed() {
}

type Objects []Interface

func (me Objects) Len() (int) {
	return len(me)
}

func (me Objects) Src(i int) (x0, y0, x1, y1 int) {
	return me[i].Src()
}

func (me Objects) Dst(i int) (x0, y0, x1, y1 int) {
	return me[i].Dst()
}

func NewObjects() (Objects) {
	return nil
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
func NewHitObjects(who ...Objects) (*HitObjects) {
	rects := []fig.Rect{}
	for _, x := range who {
		for i := 0; i < x.Len(); i++ {
			rects = append(rects, x[i].HitRects()...)
		}
	}
	return &HitObjects {
		Objs: rects,
	}
}

func SetInput(bits ginput.InputBits, who ...Objects) {
	for _, x := range who {
		for i := 0; i < x.Len(); i++ {
			x[i].SetInput(bits)
		}
	}
}

func Update(trigger event.Trigger, who ...Objects) {
	for _, x := range who {
		for i := 0; i < x.Len(); i++ {
			x[i].Update(trigger)
		}
	}
}

//########################################
//# HitRect
//########################################
func IsHit(a, b []fig.Rect) (bool) {
	for i, _ := range a {
		for j, _ := range b {
			if a[i].Hit(b[j]) {
				return true
			}
		}
	}
	return false
}

func HitCheck(subjective Objects, objective ...Objects) {
	for a := 0; a < subjective.Len(); a++ {
		for _, objs := range objective{
			for b := 0; b < objs.Len(); b++ {
				if IsHit(subjective[a].HitRects(), objs[b].HitRects()) {
					subjective[a].Hit(objs[b])
					objs[b].Hit(subjective[a])
				}
			}
		}
	}
}


//########################################
//# OutScreen
//########################################
func InScreen(inner fig.Rect, who ...Objects) {
	for _, v := range who {
		for i := 0; i < v.Len(); i++ {
			pt := v[i].GetPoint()
			x := pt.X
			y := pt.Y
			if x < inner.Left {
				pt.X = inner.Left
			}
			if x > inner.Right {
				pt.X = inner.Right
			}
			if y < inner.Top {
				pt.Y = inner.Top
			}
			if y > inner.Bottom {
				pt.Y = inner.Bottom
			}
			v[i].SetPoint(pt)
		}
	}
}

//########################################
//# OutScreen
//########################################
func GoOutside(outer fig.Rect, who ...Objects) {
	for _, x := range who {
		for i := 0; i < x.Len(); i++ {
			if !outer.In(x[i].GetPoint()) {
				x[i].Vanish()
			}
		}
	}
}

//########################################
//# Clean
//########################################
func Clean(objs Objects) (Objects) {
	res := Objects{}
	for _, o := range objs {
		if !o.IsVanish() {
			res = append(res, o)
		}
	}
	return res
}

func DrawOptions(objs Objects) (*ebiten.DrawImageOptions) {
	return &ebiten.DrawImageOptions {
		ImageParts: objs,
	}
}

//########################################
//# ExtraDraw
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

func RotaOptions(o Interface) (*ebiten.DrawImageOptions) {
	sx0, sy0, sx1, sy1 := o.Src()
	dx0, dy0, dx1, dy1 := o.Dst()
	opt := ebiten.DrawImageOptions {
		ImageParts: NewOneSprites(sx0, sy0, sx1, sy1, dx0, dy0, dx1, dy1),
	}

	pt := o.GetPoint()
	opt.GeoM.Translate(float64(-pt.X), float64(-pt.Y))
	opt.GeoM.Rotate(float64(o.Direction()))
	opt.GeoM.Translate(float64(pt.X), float64(pt.Y))
	return &opt
}

func DrawImageRota(screen, source *ebiten.Image, who ...Objects) (error) {
	for _, x := range who {
		for i := 0; i < x.Len(); i++ {
			if err := screen.DrawImage(source, RotaOptions(x[i])); err != nil {
				return err
			}
		}
	}
	return nil
}



//########################################
//# CarryPress
//########################################
func CarryPress(subjective Objects, objective Objects) {
	for a := 0; a < subjective.Len(); a++ {
		for b := 0; b < objective.Len(); b++ {
			if IsHit(subjective[a].HitRects(), objective[b].HitRects()) {
				subjective[a].Pushed()
				objective[b].Hit(subjective[a])
			}
		}
	}
}
