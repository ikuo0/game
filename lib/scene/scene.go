
package scene

import (
	"github.com/hajimehoshi/ebiten"
)

type Parameter []string

type Interface interface {
	Main(*ebiten.Image) (bool)
	Dispose()
	ReturnValue() (Parameter)
}

type New func(Parameter) (Interface)


