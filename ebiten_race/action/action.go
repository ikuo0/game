
package action

import (
	"github.com/ikuo0/game/lib/event"
	"github.com/ikuo0/game/lib/fig"
	"github.com/ikuo0/game/lib/ginput"
	"github.com/ikuo0/game/lib/radian"
	"github.com/ikuo0/game/lib/script"
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
	GetRect() (fig.Rect)
	HitRect(Interface)
	GetLine() (fig.Line)
	HitLine(Interface)
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

func (me *Object) GetRect() (fig.Rect) {
	return fig.Rect{}
}

func (me *Object) HitRect(Interface) {
}

func (me *Object) GetLine() (fig.Line) {
	return fig.Line{}
}

func (me *Object) HitLine(Interface) {
}

//########################################
//# Objects
//########################################
type Objects struct {
	Objs []Interface
}

func (me *Objects) SetPoint(i int, pt fig.Point) {
	me.Objs[i].SetPoint(pt)
}

func (me *Objects) GetObject(i int) (Interface) {
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

func (me *Objects) SetInput(i int, bits ginput.InputBits) {
	me.Objs[i].SetInput(bits)
}

func (me *Objects) GetRect(i int) (fig.Rect) {
	return me.Objs[i].GetRect()
}

func (me *Objects) HitRect(i int, obj Interface) {
	me.Objs[i].HitRect(obj)
}

func (me *Objects) GetLine(i int) (fig.Line) {
	return me.Objs[i].GetLine()
}

func (me *Objects) HitLine(i int, obj Interface) {
	me.Objs[i].HitLine(obj)
}

func (me *Objects) Update(i int, trigger event.Trigger) {
	me.Objs[i].Update(trigger)
}

func (me *Objects) Vanish(i int) {
	me.Objs[i].Vanish()
}

func (me *Objects) Clean(i int) {
	newObjs := []Interface{}
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

func (me *Objects) Occure(objIf Interface) {
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
func NewHitObjects(who ...CanRectHit) (*HitObjects) {
	rects := []fig.Rect{}
	for _, x := range who {
		for i := 0; i < x.Len(); i++ {
			rects = append(rects, x.GetRect(i))
		}
	}
	return &HitObjects {
		Objs: rects,
	}
}



//########################################
//# SetInput
//########################################
type CanOperate interface {
	Len() (int)
	SetInput(int, ginput.InputBits)
}

func SetInput(bits ginput.InputBits, who ...CanOperate) {
	for _, x := range who {
		for i := 0; i < x.Len(); i++ {
			x.SetInput(i, bits)
		}
	}
}

//########################################
//# Update
//########################################
type Updatable interface {
	Len() (int)
	Update(i int, trigger event.Trigger)
}

func Update(trigger event.Trigger, who ...Updatable) {
	for _, x := range who {
		for i := 0; i < x.Len(); i++ {
			x.Update(i, trigger)
		}
	}
}

//########################################
//# Script
//########################################
type HasScript interface {
	Len() (int)
	Stack(i  int) (*script.Stack)
	GetObject(i int) (Interface)
}

func Script(input script.Input, output event.Trigger, who ...HasScript) {
	for _, x := range who {
		for i := 0; i < x.Len(); i++ {
			script.Exec(input, x.Stack(i), x.GetObject(i), output)
		}
	}
}

//########################################
//# HitRect
//########################################
type CanRectHit interface {
	Len() (int)
	GetRect(int) (fig.Rect)
	HitRect(int, Interface)
	GetObject(int) (Interface)
}

func HitCheckRect(subjective CanRectHit, objective ...CanRectHit) {
	for a := 0; a < subjective.Len(); a++ {
		sRect := subjective.GetRect(a)
		for _, objs := range objective{
			for b := 0; b < objs.Len(); b++ {
				if sRect.Hit(objs.GetRect(b)) {
					subjective.HitRect(a, objs.GetObject(b))
					objs.HitRect(a, subjective.GetObject(a))
				}
			}
		}
	}
}


//########################################
//# HitRect
//########################################
type CanLineHit interface {
	Len() (int)
	GetLine(int) (fig.Line)
	HitLine(int, Interface)
	GetObject(int) (Interface)
}

func IsHit(a, b fig.Line) (bool) {
	return false
}

func HitCheckLine(subjective CanLineHit, objective ...CanLineHit) {
	for a := 0; a < subjective.Len(); a++ {
		sLine := subjective.GetLine(a)
		for _, objs := range objective{
			for b := 0; b < objs.Len(); b++ {
				if sLine.Hit(objs.GetLine(b)) {
					subjective.HitLine(a, objs.GetObject(b))
					objs.HitLine(b, subjective.GetObject(a))
				}
			}
		}
	}
}

//########################################
//# OutScreen
//########################################
type InTheScreen interface {
	Len() (int)
	GetObject(i int) (Interface)
	SetPoint(int, fig.Point)
}
func InScreen(inner fig.Rect, who ...InTheScreen) {
	for _, v := range who {
		for i := 0; i < v.Len(); i++ {
			pt := v.GetObject(i).GetPoint()
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
			v.SetPoint(i, pt)
		}
	}
}

//########################################
//# OutScreen
//########################################
type InTheWorld interface {
	Len() (int)
	GetObject(int) (Interface)
	Vanish(int)
}
func GoOutside(outer fig.Rect, who ...InTheWorld) {
	for _, x := range who {
		for i := 0; i < x.Len(); i++ {
			if !outer.In(x.GetObject(i).GetPoint()) {
				x.Vanish(i)
			}
		}
	}
}

//########################################
//# Clean
//########################################
type Disposer interface {
	Len() (int)
	Clean(i int)
}
func Clean(who ...Disposer) {
	for _, x := range who {
		for i := 0; i < x.Len(); i++ {
			x.Clean(i)
		}
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

type ExDrawer interface {
	Len() (int)
	DrawOption(i int) (*ebiten.DrawImageOptions)
}
func ExDraw(screen, source *ebiten.Image, who ...ExDrawer) (error) {
	for _, x := range who {
		for i := 0; i < x.Len(); i++ {
			if err := screen.DrawImage(source, x.DrawOption(i)); err != nil {
				return err
			}
		}
	}
	return nil
}
