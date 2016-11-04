package main

import (
	"github.com/ikuo0/game/ebiten_act/changer"
	"github.com/ikuo0/game/ebiten_act/keyconfig"
	"github.com/ikuo0/game/ebiten_act/title"
	"github.com/ikuo0/game/ebiten_act/stage1"
	"github.com/ikuo0/game/ebiten_act/global"
	"github.com/ikuo0/game/lib/log"
	"github.com/ikuo0/game/lib/sound"
	"github.com/hajimehoshi/ebiten"
	//"github.com/hajimehoshi/ebiten/ebitenutil"
	"os"
	"path"
)

var sceneChanger = changer.NewChanger(changer.SceneList {
	"title":     title.New,
	"keyconfig": keyconfig.New,
	"stage1":    stage1.New,
})

func update(screen *ebiten.Image) error {
	//ebitenutil.DebugPrint(screen, "Hello world!")
	sound.Update()
	return sceneChanger.Play(screen)
}

func main() {
	global.SetRootPath(path.Dir(os.Args[0]))
	log.Start(global.Path().File(`log.txt`))
	sound.NewContext(global.SampleRate)
	sceneChanger.Apply([]string{"stage1"})
	ebiten.Run(update, 800, 600, 1, "Hello world!")
}
