
package funcs

import (
	"github.com/ikuo0/game/ebiten_act/global"
	"github.com/ikuo0/game/lib/fig"
	"github.com/ikuo0/game/lib/move"
	//"math"
	//"fmt"
)

//########################################
//# FaceDirection
//########################################
type FaceDirection int
const (
	FaceLeft FaceDirection = iota + 1
	FaceRight
)

func (me FaceDirection) String() (string) {
	switch(me) {
		case FaceLeft: return "FaceLeft"
		case FaceRight: return "FaceRight"
	}
	return "FaceUnknown"
}

//########################################
//# Standard Gravity
//########################################
const GravityAccel float64 = 0.6
type Gravity struct {
	*move.Gravity
}

/*
func (me *Gravity) Value() (float64) {
// 重力が小数点時、整数になるまで落下判定で着地できなくなるため+１未満は１にする
	if res := me.Gravity.Value(); res > 0 && res < 1 {
		return 1
	} else {
		return res
	}
}
*/

func NewGravity() (*Gravity) {
	return &Gravity {
		Gravity: move.NewGravity(12, 0.6),
	}
}

//########################################
//# Hit Check
//########################################
type WallStatus int

const (
	WallNone WallStatus = 0x00
	WallLeft            = 0x01
	WallRight           = 0x02
	WallTop             = 0x04
	WallBottom          = 0x08
)

func (me WallStatus) IsHit() (bool) {
	return (me != 0)
}

type FallingRects struct {
	Head fig.Rect
	Body fig.Rect
	Foot fig.Rect
}

func (me *FallingRects) HitWall(pt fig.Point, descend bool, walls []fig.Rect) (fig.Point, WallStatus) {
	//pt := fpt.ToInt()
	status := WallNone

	global.RectDebug.Clear()

	if !descend {// 上昇中
		for _, w := range walls {
			head := me.Head.Relative(pt)
			//global.RectDebug.Append(head)
			if head.Hit(&w) {
				status |= WallTop
				pt.Y += w.Bottom - head.Top
			}
		}
	}

	for _, w := range walls {
		body := me.Body.Relative(pt)
		global.RectDebug.Append(body)
		if body.Hit(&w) {
			if body.Center().X > w.Center().X {
				status |= WallLeft
				pt.X += w.Right - body.Left
			} else {
				status |= WallRight
				pt.X -= body.Right - w.Left
			}
		}
	}

	if descend {// 下降中
		for _, w := range walls {
			foot := me.Foot.Relative(pt)
			global.RectDebug.Append(foot)
			//global.RectDebug.Append(foot)
			if foot.Hit(&w) {
				status |= WallBottom
				pt.Y -= (foot.Bottom - w.Top)
			}
		}
	}

	return pt, status
}

func NewFallingRects(width, height, adjustX, adjustY float64) (*FallingRects) {
	w1 := width
	w2 := width / 2
	w4 := width / 4
	h1 := height
	h2 := height / 2
	h4 := height / 4

	head := fig.Rect{w4, 0, w4 + w2, 0}.Relative(fig.Point{adjustX, adjustY})
	body := fig.Rect{0, h4, w1, h4 + h2}.Relative(fig.Point{adjustX, adjustY})
	foot := fig.Rect{0, h1 - 1, w1, h1}.Relative(fig.Point{adjustX, adjustY})

	return &FallingRects {
		Head: head,
		Body: body,
		Foot: foot,
	}
}


//########################################
//# Enemy Occure
//########################################
type EnemyConfig struct {
	Point     fig.Point
	Direction FaceDirection
}
