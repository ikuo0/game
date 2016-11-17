
package action

import (
	"github.com/ikuo0/game/lib/event"
	"github.com/ikuo0/game/lib/fig"
	"github.com/ikuo0/game/lib/ginput"
	"github.com/ikuo0/game/lib/orig"
	"github.com/ikuo0/game/lib/radian"
	"github.com/ikuo0/game/lib/script"
	"github.com/hajimehoshi/ebiten"
)

//########################################
//# Objects
//########################################
type Object interface {
	Point() (fig.FloatPoint)
	Direction() (radian.Radian)
	Update(trigger event.Trigger)
	Vanish()
	IsVanish() (bool)
	Src() (x0, y0, x1, y1 int)
	Dst() (x0, y0, x1, y1 int)
	HitRects() ([]fig.Rect)
	Hit(Object)
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
	Origin(i int) (orig.Interface)
}

func Script(input script.Input, output event.Trigger, who ...HasScript) {
	for _, x := range who {
		for i := 0; i < x.Len(); i++ {
			script.Exec(input, x.Stack(i), x.Origin(i), output)
		}
	}
}

//########################################
//# HitCheck
//########################################
type CanHit interface {
	Len() (int)
	HitRects(int) ([]fig.Rect)
	Hit(int, Object)
	GetObject(int) (Object)
}

func IsHit(a, b []fig.Rect) (bool) {
	for i, _ := range a {
		for j, _ := range b {
			if a[i].Hit(&b[j]) {
				return true
			}
		}
	}
	return false
}

func HitCheck(subjective CanHit, objective ...CanHit) {
	for a := 0; a < subjective.Len(); a++ {
		for _, objs := range objective{
			for b := 0; b < objs.Len(); b++ {
				if IsHit(subjective.HitRects(a), objs.HitRects(b)) {
					subjective.Hit(a, objs.GetObject(b))
					objs.Hit(b, subjective.GetObject(a))
				}
			}
		}
	}
}

func UniHitCheck(subjective CanHit, objective ...CanHit) {
	for a := 0; a < subjective.Len(); a++ {
		for _, objs := range objective{
			for b := 0; b < objs.Len(); b++ {
				if IsHit(subjective.HitRects(a), objs.HitRects(b)) {
					subjective.Hit(a, objs.GetObject(b))
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
	HitWall(int, Object)
	//Hit(int, Object)
	GetObject(int) (Object)
	Expel(int, []fig.Rect)
}
func HitWall(subjective LawsOfPhisics, allWalls ...CanHit) {
	hitWalls := []fig.Rect{}
	for a := 0; a < subjective.Len(); a++ {
		for _, walls := range allWalls {
			for b := 0; b < walls.Len(); b++ {
				if IsHit(subjective.HitRects(a), walls.HitRects(b)) {
					subjective.HitWall(a, walls.GetObject(b))
					walls.Hit(b, subjective.GetObject(a))
					hitWalls = append(hitWalls, walls.HitRects(b)...)
				}
			}
		}
	}
	for i := 0; i < subjective.Len(); i++ {
		subjective.Expel(i, hitWalls)
	}
}

//########################################
//# OutScreen
//########################################
type InTheScreen interface {
	Len() (int)
	Origin(i int) (orig.Interface)
	SetPoint(int, fig.FloatPoint)
}
func InScreen(inner fig.Rect, who ...InTheScreen) {
	for _, v := range who {
		for i := 0; i < v.Len(); i++ {
			pt := v.Origin(i).Point()
			x := int(pt.X)
			y := int(pt.Y)
			if x < inner.Left {
				pt.X = float64(inner.Left)
			}
			if x > inner.Right {
				pt.X = float64(inner.Right)
			}
			if y < inner.Top {
				pt.Y = float64(inner.Top)
			}
			if y > inner.Bottom {
				pt.Y = float64(inner.Bottom)
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
	Origin(int) (orig.Interface)
	Vanish(int)
}
func GoOutside(outer fig.Rect, who ...InTheWorld) {
	for _, x := range who {
		for i := 0; i < x.Len(); i++ {
			if !outer.InF(x.Origin(i).Point()) {
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

