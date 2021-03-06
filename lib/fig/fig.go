// 漢字
package fig

import (
	"fmt"
	//"../mathlib"
)

type IntPoint struct {
	X int
	Y int
}

//########################################
//# Point
//########################################
type Point struct {
	X float64
	Y float64
}

func (me Point) Equal(pt Point) (bool) {
	return me.X == pt.X && me.Y == pt.Y
}



func (me Point) Diff(pt Point) (Point) {
	return Point {
		X: me.X - pt.X,
		Y: me.Y - pt.Y,
	}
}

func (me Point) ToInt() (IntPoint) {
	return IntPoint {
		X: int(me.X),
		Y: int(me.Y),
	}
}


//########################################
//# Line
//########################################
type Line struct {
	Start Point
	End   Point
}

func (me Line) Hit(l Line) (bool) {
	ta := (l.Start.X - l.End.X) * (me.Start.Y - l.Start.Y) + (l.Start.Y - l.End.Y) * (l.Start.X - me.Start.X);
	tb := (l.Start.X - l.End.X) * (me.End.Y - l.Start.Y) + (l.Start.Y - l.End.Y) * (l.Start.X - me.End.X);
	tc := (me.Start.X - me.End.X) * (l.Start.Y - me.Start.Y) + (me.Start.Y - me.End.Y) * (me.Start.X - l.Start.X);
	td := (me.Start.X - me.End.X) * (l.End.Y - me.Start.Y) + (me.Start.Y - me.End.Y) * (me.Start.X - l.End.X);
	return tc * td < 0 && ta * tb < 0;
}

func (me Line) Relative(pt Point) (Line) {
	return Line {
		Point{me.Start.X + pt.X, me.Start.Y + pt.Y},
		Point{me.End.X + pt.X, me.End.Y + pt.Y},
	}
}

//########################################
//# IntRect
//########################################
type IntRect struct {
	Left   int
	Top    int
	Right  int
	Bottom int
}

func (me IntRect) Width() (int) {
	return me.Right - me.Left
}

func (me IntRect) Height() (int) {
	return me.Bottom - me.Top
}

//########################################
//# Rect
//########################################
type Rect struct {
	Left   float64
	Top    float64
	Right  float64
	Bottom float64
}

func (me Rect) String() (string) {
	return fmt.Sprintf("%d, %d, %d, %d", me.Left, me.Top, me.Right, me.Bottom)
}

func (me Rect) Width() (float64) {
	return me.Right - me.Left
}

func (me Rect) Height() (float64) {
	return me.Bottom - me.Top
}

func (me Rect) Center() (Point) {
	return Point{
		X: me.Left + me.Width() / 2,
		Y: me.Top + me.Height() / 2,
	}
}

func (me Rect) Relative(pt Point) (Rect) {
	me.Left   += pt.X
	me.Right  += pt.X
	me.Top    += pt.Y
	me.Bottom += pt.Y
	return me
}

func (me *Rect) Add(x, y float64) {
	me.Left   += x
	me.Right  += x
	me.Top    += y
	me.Bottom += y
}

func (me Rect) Hit(r Rect) (bool) {
	return me.Left < r.Right &&
		me.Right > r.Left &&
		me.Top < r.Bottom &&
		me.Bottom > r.Top
}

func (me Rect) In(pt Point) (bool) {
	return me.Left < pt.X &&
		me.Right > pt.X &&
		me.Top < pt.Y &&
		me.Bottom > pt.Y
}

func (me Rect) LeftLine() (Line) {
	return Line {
		Start: Point {me.Left, me.Top},
		End:   Point {me.Left, me.Bottom},
	}
}

func (me Rect) TopLine() (Line) {
	return Line {
		Start: Point {me.Left, me.Top},
		End:   Point {me.Right, me.Top},
	}
}

func (me Rect) RightLine() (Line) {
	return Line {
		Start: Point {me.Right, me.Top},
		End:   Point {me.Right, me.Bottom},
	}
}

func (me Rect) BottomLine() (Line) {
	return Line {
		Start: Point {me.Left, me.Bottom},
		End:   Point {me.Right, me.Bottom},
	}
}

//########################################
//# Funcs
//########################################
func PointToLine(from, to Point) (Line) {
	return Line {from, to}
}
