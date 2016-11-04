
package radian

import (
	"math"
)

type Radian float64
func (me Radian) TurnRight(deg int) (Radian) {
	//Radian(i) * math.Pi / 180
	res := me + (Radian(deg) * math.Pi / 180)
	if res > math.Pi {
		diff := res - math.Pi
		res = - math.Pi + diff
	} else if res < -math.Pi {
		diff := res - math.Pi
		res = math.Pi + diff
	}
	return res
}
func (me Radian) TurnLeft(deg int) (Radian) {
	res := me - (Radian(deg) * math.Pi / 180)
	if res > math.Pi {
		diff := res - math.Pi
		res = - math.Pi + diff
	} else if res < -math.Pi {
		diff := res - math.Pi
		res = math.Pi + diff
	}
	return res
}

type Deg []Radian
func NewDeg() (Deg) {
	res := make([]Radian, 360)
	for i := 0; i < 360; i++ {
		res[i] = Radian(i) * math.Pi / 180
	}
	return res
}

func (me Deg) Up() (Radian) {return me[270];}
func (me Deg) RightUp() (Radian) {return me[315];}
func (me Deg) Right() (Radian) {return me[0];}
func (me Deg) RightDown() (Radian) {return me[45];}
func (me Deg) Down() (Radian) {return me[90];}
func (me Deg) LeftDown() (Radian) {return me[135];}
func (me Deg) Left() (Radian) {return me[180];}
func (me Deg) LeftUp() (Radian) {return me[225];}

func FromDeg(deg int) (Radian) {
	if deg < 0 {
		deg = 360 + deg
	}
	return degArray[deg % 360]
}

func Up() (Radian) {return degArray[270];}
func RightUp() (Radian) {return degArray[315];}
func Right() (Radian) {return degArray[0];}
func RightDown() (Radian) {return degArray[45];}
func Down() (Radian) {return degArray[90];}
func LeftDown() (Radian) {return degArray[135];}
func Left() (Radian) {return degArray[180];}
func LeftUp() (Radian) {return degArray[225];}
