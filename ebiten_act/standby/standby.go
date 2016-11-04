
package standby

import (
	"fmt"
	"github.com/ikuo0/game/lib/scene"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type Standby struct {
	Counter int
}

func New(args []string) (scene.Interface) {
	return &Standby {
		Counter: 0,
	}
}

func (me *Standby) Main(screen *ebiten.Image) (bool) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Standby#Main#%d", me.Counter))
	me.Counter++
	return true
}

func (me *Standby) Dispose() {
	me.Counter = 0
}

func (me *Standby) ReturnValue() (scene.Parameter) {
	return []string{}
}

