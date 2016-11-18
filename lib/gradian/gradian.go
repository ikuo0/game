
package gradian

import (
	"github.com/ikuo0/game/lib/fig"
	"github.com/ikuo0/game/lib/radian"
	"math"
	//"fmt"
)

/*
	角度を上下反転させる
*/
func NormalizeDeg(deg int) (int) {
	n := deg % 360
	n = 360 - n
	return radian.NormalizeDeg(n)
}

func ToIndex(deg int) (int) {
	deg = NormalizeDeg(deg)
	return radian.ToIndex(deg)
}

func DegreeToRadian(deg int) (radian.Radian) {
	return radian.DegArray[ToIndex(deg)]
}

func RadianToDegree(rad radian.Radian) (int) {
	return NormalizeDeg(int(radian.ToDeg(rad)))
}

func Up() (radian.Radian) {
	return radian.DegArray[ToIndex(90)]
}

func RightUp() (radian.Radian) {
	return radian.DegArray[ToIndex(315)]
}

func Right() (radian.Radian) {
	return radian.DegArray[ToIndex(0)]
}

func RightDown() (radian.Radian) {
	return radian.DegArray[ToIndex(45)]
}

func Down() (radian.Radian) {
	return radian.DegArray[ToIndex(90)]
}

func LeftDown() (radian.Radian) {
	return radian.DegArray[ToIndex(135)]
}

func Left() (radian.Radian) {
	return radian.DegArray[ToIndex(180)]
}

func LeftUp() (radian.Radian) {
	return radian.DegArray[ToIndex(225)]
}

func Aim(to, from fig.Point) (radian.Radian) {
	//return radian.Radian(math.Atan2(to.Y - from.Y, to.X - from.X))
	return radian.Radian(math.Atan2(float64(to.Y - from.Y), float64(to.X - from.X)))
}

//########################################
//# Degree
//########################################
type Degree struct {
	Deg int
}

func (me *Degree) TurnRight(n int) (*Degree) {
	me.Deg -= n
	return me
}

func (me *Degree) TurnLeft(n int) (*Degree) {
	me.Deg += n
	return me
}

func (me *Degree) Radian() (radian.Radian) {
	return radian.DegArray[ToIndex(me.Deg)]
}

