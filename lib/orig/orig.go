
package orig

import (
	"github.com/ikuo0/game/lib/fig"
	"github.com/ikuo0/game/lib/radian"
)

type Interface interface {
	GetPoint() (fig.Point)
	Direction() (radian.Radian)
}
