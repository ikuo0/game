
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
	HitRects() ([]fig.Rect)
	Hit(Interface)
	HitWall(obj Interface)
	Expel([]fig.Rect)
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

func (me *Object) HitWall(obj Interface) {
}

func (me *Object) Expel([]fig.Rect) {
}

//########################################
//# Objects
//########################################
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



//########################################
//# SetInput
//########################################
func SetInput(bits ginput.InputBits, who ...Objects) {
	for _, x := range who {
		for i := 0; i < x.Len(); i++ {
			x[i].SetInput(bits)
		}
	}
}

//########################################
//# Update
//########################################
func Update(trigger event.Trigger, who ...Objects) {
	for _, x := range who {
		for i := 0; i < x.Len(); i++ {
			x[i].Update(trigger)
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

func UniHitCheck(subjective Objects, objective ...Objects) {
	for a := 0; a < subjective.Len(); a++ {
		for _, objs := range objective{
			for b := 0; b < objs.Len(); b++ {
				if IsHit(subjective[a].HitRects(), objs[b].HitRects()) {
					subjective[a].Hit(objs[b])
				}
			}
		}
	}
}

//########################################
//# HitWall
//########################################
type LawsOfPhisics interface {
	Len() (int)
	HitRects(int) ([]fig.Rect)
	HitWall(int, Interface)
	//Hit(int, Object)
	GetObject(int) (Interface)
	Expel(int, []fig.Rect)
}
func HitWall(subjective Objects, allWalls ...Objects) {
	hitWalls := []fig.Rect{}
	for a := 0; a < subjective.Len(); a++ {
		for _, walls := range allWalls {
			for b := 0; b < walls.Len(); b++ {
				if IsHit(subjective[a].HitRects(), walls[b].HitRects()) {
					subjective[a].HitWall(walls[b])
					walls[b].Hit(subjective[a])
					hitWalls = append(hitWalls, walls[b].HitRects()...)
				}
			}
		}
	}
	for i := 0; i < subjective.Len(); i++ {
		subjective[i].Expel(hitWalls)
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

