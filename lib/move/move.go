
package move

import (
	"github.com/ikuo0/game/lib/gradian"
	"github.com/ikuo0/game/lib/radian"
	"math"
	//"fmt"
)

//########################################
//# Force
//########################################
type Force struct {
	Max  float64
	Rate float64
}

func (me *Force) Reset() {
	me.Rate = 0
}

func (me *Force) Accel(rating float64) {
	if me.Rate += rating; me.Rate > me.Max {
		me.Rate = me.Max
	} else if me.Rate < -me.Max {
		me.Rate = -me.Max
	}
}

func (me *Force) Frictional(friction float64) {
	if me.Rate > 0 {
		if me.Rate -= friction; me.Rate < 0 {
			me.Rate = 0
		}
	} else if me.Rate < 0 {
		if me.Rate += friction; me.Rate > 0 {
			me.Rate = 0
		}
	}
}

func (me *Force) Value() (float64) {
	return me.Rate
}

func NewForce(max float64) (*Force) {
	return &Force {max, 0}
}

//########################################
//# XYpower
//########################################
type XYpower struct {
	X float64
	Y float64
}

//########################################
//# Vector
//########################################
type Vector struct {
	Force
	gradian.Degree
}

func (me *Vector) X() (float64) {
	return math.Cos(float64(me.Radian())) * me.Value()
}

func (me *Vector) Y() (float64) {
	return math.Sin(float64(me.Radian())) * me.Value()
}

func NewVector(deg int, max float64) (*Vector) {
	return &Vector {
		Force:  Force{max, 0},
		Degree: gradian.Degree {deg},
	}
}

//########################################
//# FixedVector
//########################################
type FixedVector struct {
	Force
	C float64
	S float64
}

func (me *FixedVector) X() (float64) {
	return me.C * me.Value()
}

func (me *FixedVector) Y() (float64) {
	return me.S * me.Value()
}

func NewFixedVector(rad radian.Radian, max float64) (*FixedVector) {
	return &FixedVector {
		Force:  Force{max, 0},
		C:      math.Cos(float64(rad)),
		S:      math.Sin(float64(rad)),
	}
}

//########################################
//# Inertia
//########################################
type Inertia struct {
	Frictional float64
	Advance  Force
	Backward Force
}

func (me *Inertia) Update() {
	me.Advance.Frictional(me.Frictional)
	me.Backward.Frictional(me.Frictional)
}

func (me *Inertia) Set(force float64) {
/*
	if force < 0 {
		me.Backward.Rate = math.Abs(force)
	} else if force > 0 {
		me.Advance.Rate = force
	}
*/

	if force < -me.Backward.Rate {
		me.Backward.Rate = math.Abs(force)
	} else if force > me.Advance.Rate {
		me.Advance.Rate = force
	}
}

func (me *Inertia) Accel(force float64) {
	if force < 0 {
		v := math.Abs(force)
		me.Backward.Rate += v * me.Frictional
		if me.Backward.Rate > v {
			me.Backward.Rate = v
		}
	} else if force > 0 {
		me.Advance.Rate += force * me.Frictional
		if me.Advance.Rate > force {
			me.Advance.Rate = force
		}
	}
}

func (me *Inertia) Value() (float64) {
	return me.Advance.Value() - me.Backward.Value()
}

func NewInertia(frictional float64) (*Inertia) {
	return &Inertia {
		Advance:    Force {0, 0},
		Backward:   Force {0, 0},
		Frictional: frictional,
	}
}

//########################################
//# XYcomponent
//########################################
type XYcomponent struct {
	Xinertia Inertia
	Yinertia Inertia
}

func (me *XYcomponent) Update() {
	me.Xinertia.Update()
	me.Yinertia.Update()
}

func (me *XYcomponent) Set(x, y float64) {
	me.Xinertia.Set(x)
	me.Yinertia.Set(y)
}

func (me *XYcomponent) Power() (XYpower) {
	return XYpower {
		X: me.Xinertia.Value(),
		Y: me.Yinertia.Value(),
	}
}

func NewXYcomponent(frictional float64) (*XYcomponent) {
	return &XYcomponent {
		Xinertia: *NewInertia(frictional),
		Yinertia: *NewInertia(frictional),
	}
}

//########################################
//# GravityJump
//########################################
type Gravity struct {
	Gravity      *Force
	GravityAccel float64
}
func (me *Gravity) Update() {
	me.Gravity.Accel(me.GravityAccel)
}
func (me *Gravity) Jump(power float64) {
	me.Gravity.Rate = -power
}

func (me *Gravity) JumpCancel() {
	me.Gravity.Reset()
}

func (me *Gravity) Landing() {
	me.Gravity.Reset()
}

func (me *Gravity) Value() (float64) {
	return me.Gravity.Value()
}

func NewGravity(max, accel float64) (*Gravity) {
	return &Gravity {
		Gravity:      NewForce(max),
		GravityAccel: accel,
	}
}

//########################################
//# FallingInertia
//########################################

/*
type FallingInertia struct {
	Inertia 
}

func (me *FallingInertia) Reset() {
	me.Left.Reset()
	me.Right.Reset()
	me.Up.Reset()
	me.Down.Reset()
}

func (me *FallingInertia) Jump(power float64) {
	//me.Up.Accel(me.Up.MaxPower)
	me.Up.Accel(power)
	me.Down.Reset()
}

func (me *FallingInertia) JumpCancel() {
	me.Up.Reset()
	me.Down.Reset()
}

func (me *FallingInertia) Accel() {
	c := math.Cos(float64(me.Radian))
	if c > 0 {
		me.Right.Accel(c)
	} else if c < 0 {
		me.Left.Accel(math.Abs(c))
	}
}

func (me *FallingInertia) Chafe(v float64) {
	me.Left.Chafe(v)
	me.Right.Chafe(v)
}

func (me *FallingInertia) Fall() {
	if me.Up.Power > 0 {
		me.Up.Brake(1)
	} else {
		me.Down.Accel(1)
	}
}

*/

/*
var MaxGravity = float64(16)
var MaxJumpPower = float64(16)

func NewFallingInertia(rad radian.Radian, p, a, m float64) (*FallingInertia) {
	l := RateScalar(0, a, m)
	r := RateScalar(0, a, m)
	u := RateScalar(0, 1, MaxJumpPower)
	d := RateScalar(0, 1, MaxGravity)

	c := math.Cos(float64(rad))
	if c > 0 {
		r.Accel(c * p)
	} else if c < 0 {
		l.Accel(math.Abs(c * p))
	}

	s := math.Sin(float64(rad))
	if s > 0 {
		d.Accel(s * p)
	} else if s < 0 {
		u.Accel(math.Abs(s * p))
	}

	return &FallingInertia {
		Inertia: Inertia {
			Radian: rad,
			Left:   l,
			Right:  r,
			Up:     u,
			Down:   d,
		},
	}
}
*/
