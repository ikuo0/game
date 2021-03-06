
package player

import (
	"github.com/ikuo0/game/ebiten_race/action"
	//"github.com/ikuo0/game/ebiten_race/eventid"
	"github.com/ikuo0/game/ebiten_race/world"
	"github.com/ikuo0/game/lib/event"
	"github.com/ikuo0/game/lib/fig"
	"github.com/ikuo0/game/lib/ginput"
	//"github.com/ikuo0/game/lib/gradian"
	"github.com/ikuo0/game/lib/move"
	"github.com/ikuo0/game/lib/radian"
	"github.com/ikuo0/game/lib/script"
	//"math"
	//"fmt"
)

var ShotCommand    = []ginput.InputBits {ginput.Nkey1, ginput.Key1}
var SheldCommand   = []ginput.InputBits {ginput.Nkey2, ginput.Key2}
const Width = 48
const Height = 24
const AdjustX = -24
const AdjustY = -12
const CollisionWidth = 24
const CollisionHeight = 24
const CollisionAdjustX = -12
const CollisionAdjustY = -12
var ImageSrc = fig.IntRect {0, 0, Width, Height}

type Player struct {
	action.Object
	PrePoint    fig.Point
	V           *move.Vector
	XYcomponent *move.XYcomponent
	InputBits  ginput.InputBits
}

func (me *Player) Direction() (radian.Radian) {
	return me.V.Radian()
}

func (me *Player) Update(trigger event.Trigger) {
	world.SetPlayer(me)
	me.PrePoint = me.Point

	bits := me.InputBits

	if bits.And(ginput.Left) {
		me.V.TurnLeft(4)
	} else if bits.And(ginput.Right) {
		me.V.TurnRight(4)
	}

	if bits.And(ginput.Key1) {
		me.V.Accel(0.2)
		me.XYcomponent.Set(me.V.X(), me.V.Y())
	} else {
		me.V.Frictional(0.1)
	}

	me.XYcomponent.Update()

	p := me.XYcomponent.Power()
	me.X += p.X
	me.Y += p.Y

	/*
	me.X += math.Cos(float64(me.V.Radian)) * 4
	me.Y += math.Sin(float64(me.V.Radian)) * 4
	*/
}
func (me *Player) Src() (x0, y0, x1, y1 int) {
	return ImageSrc.Left, ImageSrc.Top, ImageSrc.Right, ImageSrc.Bottom
}
func (me *Player) Dst() (x0, y0, x1, y1 int) {
	x, y := int(me.X) + AdjustX, int(me.Y) + AdjustY
	return x, y, x + Width, y + Height
}
func (me *Player) SetInput(bits ginput.InputBits) {
	me.InputBits = bits
}
func (me *Player) GetRect() (fig.Rect) {
	x, y := me.X + CollisionAdjustX, me.Y + CollisionAdjustY
	return fig.Rect{x, y, x + CollisionWidth, y + CollisionHeight}
}

func (me *Player) HitRect(obj action.Interface) {
}

func (me *Player) GetLine() (fig.Line) {
	return fig.Line{me.PrePoint, me.Point}
}

func (me *Player) HitLine(obj action.Interface) {
/*
	if rects := obj.HitRects(); len(rects) != 1 {
		return
	} else {
		rect := rects[0]
		ptFrom := me.PrePoint
		ptTo := me.Point
		myVector := fig.PointToLine(ptFrom, ptTo)
		if ptFrom.Equal(ptTo) {
			//fmt.Println("Equal")
			//return
		} else {
			//fmt.Println(myVector)
		}

		myRad := gradian.Aim(me.PrePoint, me.Point)

		var wallRad radian.Radian
		wallDeg := int(0)
		if rect.LeftLine().Hit(&myVector) {
			fmt.Println("Left Line")
			wallRad = gradian.DegreeToRadian(90)
			wallDeg = 90
		} else if rect.TopLine().Hit(&myVector) {
			fmt.Println("Top Line")
			wallRad = gradian.DegreeToRadian(0)
			wallDeg = 0
		} else if rect.RightLine().Hit(&myVector) {
			fmt.Println("Right Line")
			wallRad = gradian.DegreeToRadian(90)
			wallDeg = 90
		} else if rect.BottomLine().Hit(&myVector) {
			fmt.Println("Bottom Line")
			wallRad = gradian.DegreeToRadian(0)
			wallDeg = 0
		} else {
			wallRad = myRad + (math.Pi / 2)
		}

		newDeg := (me.V.Degree.Deg * -1) + wallDeg

		//fmt.Println("wallRad", wallRad, radian.ToDeg(wallRad))
		//me.V.Degree.Deg = gradian.RadianToDegree(myRad - wallRad)
		me.V.Degree.Deg = newDeg
		return
		fmt.Println(wallRad)
*/

/*
		reflectRad := myRad - wallRad
		fmt.Println(radian.ToDeg(myRad), radian.ToDeg(wallRad), radian.ToDeg(reflectRad))
*/
		/*
		aim := me.Point
		from := me.PrePoint
		aimRad := radian.Radian(math.Atan2(aim.Y - from.Y, aim.X - from.X))
		*/
//	}


	/*
	//度 = ラジアン × 180 ÷ 円周率
	deg1 := me.V.Radian * 180 / math.Pi
	deg2 := reaction * 180 / math.Pi
	fmt.Println(deg1, deg2)
	*/

	
/*
	aim := me.Point
	orig := obj.GetPoint()
	aimRad := radian.Radian(math.Atan2(aim.Y - orig.Y, aim.X - orig.X))
	*/
	//me.V.Radian = me.V.Radian + math.Pi


/*
	c := math.Cos(float64(aimRad))
	if c > 0 {
		me.V.Right.Accel(c)
	} else if c < 0 {
		me.V.Left.Accel(math.Abs(c))
	}

	s := math.Sin(float64(aimRad))
	if s > 0 {
		me.V.Down.Accel(s)
	} else if s < 0 {
		me.V.Up.Accel(math.Abs(s))
	}
	*/


/*
	c := math.Cos(float64(me.Radian))
	s := math.Sin(float64(me.Radian))
*/

/*
			aim := world.GetPlayer().GetPoint()
			aimRad := radian.Radian(math.Atan2(aim.Y - me.Y, aim.X - me.X))
			trigger.EventTrigger(eventid.Bullet2, aimRad, me)
			me.Timer.Start(10000)
*/

}

func (me *Player) Stack() (*script.Stack) {
	return nil
}

const MaxSpeed float64 = 6
func NewPlayer(pt fig.Point) (*Player) {
	return &Player{
		Object: action.Object {
			Point:  pt,
		},
		V:           move.NewVector(90, MaxSpeed),
		XYcomponent: move.NewXYcomponent(0.05),
	}
}
