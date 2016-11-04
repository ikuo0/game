
package explosion

import (
	"github.com/ikuo0/game/lib/event"
	"github.com/ikuo0/game/lib/anime"
	"github.com/ikuo0/game/lib/fig"
	"github.com/ikuo0/game/lib/radian"
)

//########################################
//# Explosion1
//########################################
var SrcExplosion1 = []fig.Rect {
	{
		0,
		0,
		0 + 64,
		0 + 64,
	},
	{
		0,
		64,
		0 + 64,
		64 + 64,
	},
	{
		0,
		128,
		0 + 64,
		128 + 64,
	},
	{
		0,
		196,
		0 + 64,
		192 + 64,
	},

	{
		64,
		0,
		64 + 64,
		0 + 64,
	},
	{
		64,
		64,
		64 + 64,
		64 + 64,
	},
	{
		64,
		128,
		64 + 64,
		128 + 64,
	},
	{
		64,
		196,
		64 + 64,
		192 + 64,
	},

	{
		128,
		0,
		128 + 64,
		0 + 64,
	},
	{
		128,
		64,
		128 + 64,
		64 + 64,
	},
	{
		128,
		128,
		128 + 64,
		128 + 64,
	},
	{
		128,
		196,
		128 + 64,
		192 + 64,
	},
}

type Explosion1 struct {
	fig.FloatPoint
	Anime *anime.Frames
	Vanished bool
}

func (me *Explosion1) Point() (fig.FloatPoint) {
	return me.FloatPoint
}

func (me *Explosion1) Direction() (radian.Radian) {
	return 0
}

func (me *Explosion1) Update(trigger event.Trigger) {
	if me.Anime.Arounded() {
		me.Vanish()
	} else {
		me.Anime.Update()
	}
}

func (me *Explosion1) Vanish() {
	me.Vanished = true
}
func (me *Explosion1) IsVanish() (bool) {
	return me.Vanished
}
func (me *Explosion1) Src() (x0, y0, x1, y1 int) {
	x := SrcExplosion1[me.Anime.Index()]
	return x.Left, x.Top, x.Right, x.Bottom
}
func (me *Explosion1) Dst() (x0, y0, x1, y1 int) {
	x, y := int(me.X) - 32, int(me.Y) - 32
	return x, y, x + 64, y + 64
}
func (me *Explosion1) HitRects() ([]fig.Rect) {
	return nil
}

func (me *Explosion1) Hit() {
}

func NewExplosion1(pt fig.FloatPoint) (*Explosion1) {
	return &Explosion1 {
		FloatPoint: pt,
		Anime:      anime.NewFrames(
			2, 2, 2, 2,
			2, 2, 2, 2,
			2, 2, 2, 2,
			),
	}
}

//########################################
//# Vanishing1
//########################################
var SrcVanishing = []fig.Rect {
	{
		192,
		0,
		192 + 16,
		0 + 16,
	},
	{
		192,
		16,
		192 + 16,
		16 + 16,
	},
	{
		192,
		32,
		192 + 16,
		32 + 16,
	},
	{
		192,
		48,
		192 + 16,
		48 + 16,
	},
}

type Vanishing1 struct {
	*Explosion1
}

func (me *Vanishing1) Src() (x0, y0, x1, y1 int) {
	x := SrcVanishing[me.Anime.Index()]
	return x.Left, x.Top, x.Right, x.Bottom
}
func (me *Vanishing1) Dst() (x0, y0, x1, y1 int) {
	x, y := int(me.X) - 8, int(me.Y) - 8
	return x, y, x + 16, y + 16
}

func NewVanishing1(pt fig.FloatPoint) (*Vanishing1) {
	return &Vanishing1 {
		Explosion1: &Explosion1 {
			FloatPoint: pt,
			Anime:      anime.NewFrames(
			8, 8, 8, 8,
			),
		},
	}
}