
package dbgrect

import (
	"github.com/ikuo0/game/ebiten_act/action"
	"github.com/ikuo0/game/lib/fig"
	//"fmt"
)

//########################################
//# DebugRect
//########################################
var ImageSource = []fig.IntRect {
	{
		0,
		0,
		32,
		32,
	},
}

type DebugRect struct {
	action.Object
	Rects []fig.Rect
}

func (me *DebugRect) Append(rects ...fig.Rect) {
	for _, v := range rects {
		me.Rects = append(me.Rects, v)
	}
}

func (me *DebugRect) Clear() {
	me.Rects = nil
}

func (me *DebugRect) HitRects() ([]fig.Rect) {
	return me.Rects
}

func New(fig.Rect) (*DebugRect) {
	return &DebugRect {
	}
}
