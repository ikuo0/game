
package radian

import (
	"math"
	//"fmt"
)

type Radian float64
/*
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
*/
func (me Radian) Deg() (float64) {
	return float64(me * 180 / math.Pi)
}

func ToDeg(r Radian) (float64) {
	return float64(r * 180 / math.Pi)
}

func ToRad(deg int) (Radian) {
	return Radian(float64(deg) * math.Pi / 180)
}

func NormalizeDeg(deg int) (int) {
	n := deg % 360
	if n < 0 {
		n += 360
	}
	if n > 179 {
		return n - 360
	} else {
		return n
	}
}

func ToIndex(deg int) (int) {
	deg = NormalizeDeg(deg)
	if deg < 0 {
		return deg + 360
	} else {
		return deg
	}
}

