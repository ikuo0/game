
package player

import (
	//"github.com/ikuo0/game/ebiten_race/eventid"
	"github.com/ikuo0/game/ebiten_race/world"
	"github.com/ikuo0/game/lib/action"
	"github.com/ikuo0/game/lib/event"
	"github.com/ikuo0/game/lib/fig"
	"github.com/ikuo0/game/lib/ginput"
	"github.com/ikuo0/game/lib/gradian"
	"github.com/ikuo0/game/lib/move"
	"github.com/ikuo0/game/lib/radian"
	"github.com/ikuo0/game/lib/script"
	"github.com/ikuo0/game/lib/sprites"
	"math"
	"fmt"
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
var ImageSrc = fig.Rect {0, 0, Width, Height}

type Player struct {
	fig.FloatPoint
	PrePoint   fig.FloatPoint
	V           *move.Vector
	XYcomponent *move.XYcomponent
	InputBits  ginput.InputBits
}

func (me *Player) Point() (fig.FloatPoint) {
	return me.FloatPoint
}

func (me *Player) Direction() (radian.Radian) {
	return me.V.Radian()
}

func (me *Player) Update(trigger event.Trigger) {
	world.SetPlayer(me)
	me.PrePoint = me.FloatPoint

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

func (me *Player) Vanish() {
}
func (me *Player) IsVanish() (bool) {
	return false
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
func (me *Player) HitRects() ([]fig.Rect) {
	x, y := int(me.X) + CollisionAdjustX, int(me.Y) + CollisionAdjustY
	return []fig.Rect{{x, y, x + CollisionWidth, y + CollisionHeight}}
}

func (me *Player) Hit(obj action.Object) {
	if rects := obj.HitRects(); len(rects) != 1 {
		return
	} else {
		return
		rect := rects[0]
		ptFrom := me.PrePoint.ToInt()
		ptTo := me.FloatPoint.ToInt()
		myVector := fig.PointToLine(ptFrom, ptTo)
		if ptFrom.Equal(ptTo) {
			fmt.Println("Equal")
			return
		} else {
			//fmt.Println(myVector)
		}

		myRad := gradian.Aim(me.PrePoint, me.FloatPoint)

		var wallRad radian.Radian
		if rect.LeftLine().Hit(&myVector) {
			fmt.Println("Left Line")
			wallRad = gradian.DegreeToRadian(90)
		} else if rect.TopLine().Hit(&myVector) {
			fmt.Println("Top Line")
			wallRad = gradian.DegreeToRadian(0)
		} else if rect.RightLine().Hit(&myVector) {
			fmt.Println("Right Line")
			wallRad = gradian.DegreeToRadian(90)
		} else if rect.BottomLine().Hit(&myVector) {
			fmt.Println("Bottom Line")
			wallRad = gradian.DegreeToRadian(0)
		} else {
			wallRad = myRad + (math.Pi / 2)
		}

		fmt.Println("wallRad", wallRad, radian.ToDeg(wallRad))

/*
		reflectRad := myRad - wallRad
		fmt.Println(radian.ToDeg(myRad), radian.ToDeg(wallRad), radian.ToDeg(reflectRad))
*/
		/*
		aim := me.FloatPoint
		from := me.PrePoint
		aimRad := radian.Radian(math.Atan2(aim.Y - from.Y, aim.X - from.X))
		*/
	}


	/*
	//度 = ラジアン × 180 ÷ 円周率
	deg1 := me.V.Radian * 180 / math.Pi
	deg2 := reaction * 180 / math.Pi
	fmt.Println(deg1, deg2)
	*/

	
/*
	aim := me.FloatPoint
	orig := obj.Point()
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
			aim := world.GetPlayer().Point()
			aimRad := radian.Radian(math.Atan2(aim.Y - me.Y, aim.X - me.X))
			trigger.EventTrigger(eventid.Bullet2, aimRad, me)
			me.Timer.Start(10000)
*/
}
func (me *Player) Stack() (*script.Stack) {
	return nil
}

const MaxSpeed float64 = 6
func NewPlayer(pt fig.FloatPoint) (*Player) {
	return &Player{
		FloatPoint:  pt,
		V:           move.NewVector(90, MaxSpeed),
		XYcomponent: move.NewXYcomponent(0.05),
	}
}


//########################################
//# Objects
//########################################
type Interface interface {
	action.Object
	SetInput(ginput.InputBits)
}

type Objects struct {
	*sprites.RotaObjects
}
func (me *Objects) Get(i int) (Interface) {
	return me.Objs[i].(Interface)
}
func (me *Objects) SetInput(i int, bits ginput.InputBits) {
	me.Get(i).SetInput(bits)
}
func NewObjects() (*Objects) {
	return &Objects {
		RotaObjects: sprites.NewRotaObjects(),
	}
}

