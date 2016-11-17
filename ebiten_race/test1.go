package main

import (
	"github.com/ikuo0/game/ebiten_race/changer"
	"github.com/ikuo0/game/ebiten_race/global"
	"github.com/ikuo0/game/ebiten_race/keyconfig"
	"github.com/ikuo0/game/ebiten_race/title"
	"github.com/ikuo0/game/ebiten_race/stage"

	"github.com/ikuo0/game/lib/ginput"
	"github.com/ikuo0/game/lib/log"
	"github.com/ikuo0/game/lib/scene"
	"github.com/ikuo0/game/lib/sound"
	"github.com/hajimehoshi/ebiten"
	"os"
	"path"
)

type Initialize struct {
}

func (me *Initialize) Main(screen *ebiten.Image) (bool) {
	ginput.Initialize()
	return false
}

func (me *Initialize) Dispose() {
}

func (me *Initialize) ReturnValue() (scene.Parameter) {
	return []string{"stage"}
}

func NewInitialize(args scene.Parameter) (scene.Interface) {
	return &Initialize {}
}

var sceneChanger = changer.NewChanger(changer.SceneList {
	"initialize": NewInitialize,
	"title":      title.New,
	"keyconfig":  keyconfig.New,
	"stage":      stage.New,
})

func update(screen *ebiten.Image) error {
	sound.Update()
	return sceneChanger.Play(screen)
}

func main() {
	global.SetRootPath(path.Dir(os.Args[0]))
	log.Start(global.Path().File(`log.txt`))
	sound.Initialize(global.SampleRate, 20, 20)
	sceneChanger.Apply([]string{"initialize"})
	ebiten.Run(update, 800, 600, 1, "EBITEN RACE")
}
