
package action

import (
	"github.com/ikuo0/game/lib/action"
	"github.com/ikuo0/game/lib/fig"
)

//########################################
//# CarryPress
//########################################
type CanPress interface {
	Len() (int)
	HitRects(int) ([]fig.Rect)
	Pushed(int)
	GetObject(int) (action.Object)
}

func CarryPress(subjective CanPress, objective action.CanHit) {
	for a := 0; a < subjective.Len(); a++ {
		for b := 0; b < objective.Len(); b++ {
			if action.IsHit(subjective.HitRects(a), objective.HitRects(b)) {
				subjective.Pushed(a)
				objective.Hit(b, subjective.GetObject(a))
			}
		}
	}
}
