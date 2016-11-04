
package sprites

import (
	"github.com/ikuo0/game/lib/event"
	"github.com/ikuo0/game/lib/fig"
	"github.com/ikuo0/game/lib/ginput"
	"github.com/ikuo0/game/lib/orig"
	"github.com/ikuo0/game/lib/radian"
	"github.com/ikuo0/game/lib/script"
	"github.com/hajimehoshi/ebiten"
)

type Interface interface {
	Len() (int)
	Dst(i int) (x0, y0, x1, y1 int)
	Src(i int) (x0, y0, x1, y1 int)
	Update(i int, trigger event.Trigger)
	Clean(i int)
	Options() (*ebiten.DrawImageOptions)
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
	Hit(int)
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
					subjective.Hit(a)
					objs.Hit(b)
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
	HitWall(int, []fig.Rect)
	Expel(int)
}
func HitWall(subjective LawsOfPhisics, allWalls ...CanHit) {
	for a := 0; a < subjective.Len(); a++ {
		for _, walls := range allWalls {
			for i := 0; i < walls.Len(); i++ {
				if IsHit(subjective.HitRects(a), walls.HitRects(i)) {
					subjective.HitWall(a, walls.HitRects(i))
					walls.Hit(i)
				}
			}
		}
	}
	for i := 0; i < subjective.Len(); i++ {
		subjective.Expel(i)
	}
}

//########################################
//# CarryPress
//########################################
type CanPress interface {
	Len() (int)
	HitRects(int) ([]fig.Rect)
	Pushed(int)
}

func CarryPress(subjective CanPress, objective CanHit) {
	for a := 0; a < subjective.Len(); a++ {
		for b := 0; b < objective.Len(); b++ {
			if IsHit(subjective.HitRects(a), objective.HitRects(b)) {
				subjective.Pushed(a)
				objective.Hit(b)
			}
		}
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
	Hit()
}

type Objects struct {
	Objs []Object
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

func (me *Objects) Hit(i int) {
	me.Objs[i].Hit()
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
	newObjs := []Object{}
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

func (me *Objects) Occure(objIf Object) {
	me.Objs = append(me.Objs, objIf)
}

func NewObjects() (*Objects) {
	return &Objects {}
}

//########################################
//# RotaObjects
//########################################
type Rotates interface {
	Point() (fig.FloatPoint)
	Direction() (radian.Radian)
	Update(trigger event.Trigger)
	Vanish()
	IsVanish() (bool)
	Src() (x0, y0, x1, y1 int)
	Dst() (x0, y0, x1, y1 int)
	HitRects() ([]fig.Rect)
	Hit()
	Stack() (*script.Stack)
	//Origin() (orig.Interface)
}

type RotaObjects struct {
	Objs     []Rotates
}

func (me *RotaObjects) Len() (int) {
	return len(me.Objs)
}
func (me *RotaObjects) Src(i int) (x0, y0, x1, y1 int) {
	return me.Objs[i].Src()
}
func (me *RotaObjects) Dst(i int) (x0, y0, x1, y1 int) {
	return me.Objs[i].Dst()
}
func (me *RotaObjects) Update(i int, trigger event.Trigger) {
	me.Objs[i].Update(trigger)
}
func (me *RotaObjects) Origin(i int) (orig.Interface) {
	return me.Objs[i]
}
func (me *RotaObjects) Vanish(i int) {
	me.Objs[i].Vanish()
}
func (me *RotaObjects) Clean(i int) {
	newObjs := []Rotates{}
	for _, v := range me.Objs {
		if !v.IsVanish() {
			newObjs = append(newObjs, v)
		}
	}
	me.Objs = newObjs
}

func (me *RotaObjects) HitRects(i int) ([]fig.Rect) {
	return me.Objs[i].HitRects()
}

func (me *RotaObjects) Hit(i int) {
	me.Objs[i].Hit()
}
func (me *RotaObjects) Stack(i int) (*script.Stack) {
	return me.Objs[i].Stack()
}
func (me *RotaObjects) Options() (*ebiten.DrawImageOptions) {
	return &ebiten.DrawImageOptions {
		ImageParts: me,
	}
}

func (me *RotaObjects) DrawOption(i int) (*ebiten.DrawImageOptions) {
	sx0, sy0, sx1, sy1 := me.Src(i)
	dx0, dy0, dx1, dy1 := me.Dst(i)
	opt := ebiten.DrawImageOptions {
		ImageParts: NewOneSprites(sx0, sy0, sx1, sy1, dx0, dy0, dx1, dy1),
	}

	o := me.Objs[i]
	pt := o.Point()
	opt.GeoM.Translate(-pt.X, -pt.Y)
	opt.GeoM.Rotate(float64(o.Direction()))
	opt.GeoM.Translate(pt.X, pt.Y)
	return &opt
}
func (me *RotaObjects) Occure(objIf Rotates) {
	me.Objs = append(me.Objs, objIf)
}

func NewRotaObjects() (*RotaObjects) {
	return &RotaObjects {}
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
	return r.Left, r.Top, r.Right, r.Bottom
}
func (me *HitObjects) Options() (*ebiten.DrawImageOptions) {
	return &ebiten.DrawImageOptions {
		ImageParts: me,
	}
}
func NewHitObjects(who ...CanHit) (*HitObjects) {
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
