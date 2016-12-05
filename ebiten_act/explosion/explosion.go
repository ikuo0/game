
package explosion

import (
	"github.com/ikuo0/game/ebiten_act/action"
	"github.com/ikuo0/game/lib/event"
	"github.com/ikuo0/game/lib/anime"
	"github.com/ikuo0/game/lib/fig"
)

const Width = 128
const Height = 128
const AdjustX = -64
const AdjustY = -64

func SrcCalc() ([]fig.IntRect) {
	res := []fig.IntRect{}
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			left, top := x * Width, y * Height
			right, bottom := left + Width, top + Height
			res = append(res, fig.IntRect{left, top, right, bottom})
		}
	}
	return res
}

//########################################
//# Explosion
//########################################
var ImageSrc = SrcCalc()

type Explosion struct {
	action.Object
	Anime *anime.Frames
}

func (me *Explosion) HitRects() ([]fig.Rect) {
	return nil
}

func (me *Explosion) Hit(obj action.Interface) {
}

func (me *Explosion) Update(trigger event.Trigger) {
	me.Anime.Update()
	if me.Anime.Arounded() {
		me.Vanish()
	}
}

func (me *Explosion) Src() (x0, y0, x1, y1 int) {
	x := ImageSrc[me.Anime.Index()]
	return x.Left, x.Top, x.Right, x.Bottom
}
func (me *Explosion) Dst() (x0, y0, x1, y1 int) {
	x, y := int(me.X) + AdjustX, int(me.Y) + AdjustY
	return x, y, x + Width, y + Height
}
func New(pt fig.Point, frames *anime.Frames) (*Explosion) {
	return &Explosion {
		Object: action.Object {
			Point: pt,
		},
		Anime:      frames,
	}
}

//########################################
//# Explosion1
//########################################
type Explosion1 struct {
	*Explosion
}
func NewExplosion1(pt fig.Point) (*Explosion1) {
	return &Explosion1 {
		Explosion: New(pt, anime.NewFrames(
			2, 2, 2, 2, 2, 
			2, 2, 2, 2, 2, 
			1, 1, 1, 1, 1, 
			1, 1, 1, 1, 1, 
			2, 2,
		)),
	}
}

//########################################
//# Explosion2
//########################################
type Explosion2 struct {
	*Explosion
}
func NewExplosion2(pt fig.Point) (*Explosion2) {
	return &Explosion2 {
		Explosion: New(pt, anime.NewFrames(
			4, 4, 4, 4, 4, 
			4, 4, 4, 4, 4, 
			4, 4, 4, 4, 4, 
			4, 4, 4, 4, 4, 
			4, 4,
		)),
	}
}

