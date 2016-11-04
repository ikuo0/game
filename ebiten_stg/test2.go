package main

import (
	"./changer"
	"./keyconfig"
	"./title"
	"./stage1"
	"./global"
	"./lib/log"
	"./lib/sound"
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
	sceneChanger.Apply([]string{"title"})
	ebiten.Run(update, 800, 600, 1, "Hello world!")
}
