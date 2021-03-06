
package world

import (
	"github.com/ikuo0/game/lib/fig"
	"github.com/ikuo0/game/lib/radian"
	"github.com/ikuo0/game/lib/orig"
)

type DummyPlayer struct {
}
func (me *DummyPlayer) GetPoint() (fig.Point) {
	return fig.Point{0,0}
}
func (me *DummyPlayer) Direction() (radian.Radian) {
	return 0
}

var Player orig.Interface = &DummyPlayer{}
var PlayerCount = 3

func GetPlayer() (orig.Interface) {
	return Player
}

func SetPlayer(oIf orig.Interface) {
	Player = oIf
}

func GetPlayerCount() (int) {
	return PlayerCount
}

func SetPlayerCount(count int) {
	PlayerCount = count
}

func StartGetPoint() (fig.Point) {
	return fig.Point{250, 600}
}

func Init() {
	SetPlayer(&DummyPlayer{})
	SetPlayerCount(3)
}

func Dispose() {
	Init()
}

