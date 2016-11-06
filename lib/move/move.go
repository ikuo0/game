
package move

import (
	"github.com/ikuo0/game/lib/radian"
	"math"
	//"fmt"
)

//########################################
//# Scalar
//########################################
type Scalar struct {
	Power    float64
	Rate     float64
	MaxPower float64
}

func (me *Scalar) Reset() {
	me.Power = 0
}

func (me *Scalar) Accel(v float64) {
	me.Power += v * me.Rate
	if me.Power > me.MaxPower {
		me.Power = me.MaxPower
	}
}

func (me *Scalar) Brake(v float64) {
	me.Power -= v * me.Rate
	if me.Power < 0 {
		me.Power = 0
	}
}

func (me *Scalar) Chafe(v float64) {
	if me.Power -= v; me.Power < 0 {
		me.Power = 0
	}
}

func ConstantScalar(p float64) (*Scalar) {
	return &Scalar {
		Power:    0,
		Rate:     p,
		MaxPower: p,
	}
}

func RateScalar(p, r, m float64) (*Scalar) {// 0, 0.2, 4
	return &Scalar {
		Power:    p,
		Rate:     r,
		MaxPower: m,
	}
}


//########################################
//# Inertia
//########################################
type Power struct {
	X float64
	Y float64
}

type Inertia struct {
	Radian radian.Radian
	Left   *Scalar
	Right  *Scalar
	Up     *Scalar
	Down   *Scalar
}

func (me *Inertia) Accel() {
	c := math.Cos(float64(me.Radian))
	if c > 0 {
		me.Right.Accel(c)
	} else if c < 0 {
		me.Left.Accel(math.Abs(c))
	}

	s := math.Sin(float64(me.Radian))
	if s > 0 {
		me.Down.Accel(s)
	} else if s < 0 {
		me.Up.Accel(math.Abs(s))
	}
}

func (me *Inertia) Chafe(v float64) {
	me.Left.Chafe(v)
	me.Right.Chafe(v)
	me.Up.Chafe(v)
	me.Down.Chafe(v)
}

func (me *Inertia) Power() (Power) {
	return Power {
		X: me.Right.Power - me.Left.Power,
		Y: me.Down.Power - me.Up.Power,
	}
}

func NewInertia(rad radian.Radian, p, a, m float64) (*Inertia) {
	l := RateScalar(0, a, m)
	r := RateScalar(0, a, m)
	u := RateScalar(0, a, m)
	d := RateScalar(0, a, m)

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

	return &Inertia {
		Radian: rad,
		Left:   l,
		Right:  r,
		Up:     u,
		Down:   d,
	}
}

//########################################
//# FallingInertia
//########################################
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

//########################################
//# Constant
//########################################
type Constant struct {
	Radian radian.Radian
	Xspeed  float64
	Yspeed  float64
}
func NewConstant(rad radian.Radian, speed float64) (*Constant) {
	c := math.Cos(float64(rad))
	s := math.Sin(float64(rad))
	return &Constant {
		Radian: rad,
		Xspeed: speed * c,
		Yspeed: speed * s,
	}
}

func (me *Constant) Power() (Power) {
	return Power {
		X: me.Xspeed,
		Y: me.Yspeed,
	}
}

//########################################
//# Accel
//########################################
type Accel struct {
	Radian radian.Radian
	Speed  *Scalar
	Xspeed  float64
	Yspeed  float64
}

func (me *Accel) Accel() {
	me.Speed.Accel(1)
}

func (me *Accel) Power() (Power) {
	c := math.Cos(float64(me.Radian))
	s := math.Sin(float64(me.Radian))
	return Power {
		X: me.Speed.Power * c,
		Y: me.Speed.Power * s,
	}
}

func NewAccel(rad radian.Radian, p, a, m float64) (*Accel) {
	c := math.Cos(float64(rad))
	s := math.Sin(float64(rad))
	return &Accel {
		Speed:  RateScalar(p, a, m),
		Radian: rad,
		Xspeed: p * c,
		Yspeed: p * s,
	}
}
