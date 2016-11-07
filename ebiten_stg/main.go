package main

import (
	"github.com/ikuo0/game/ebiten_stg/changer"
	"github.com/ikuo0/game/ebiten_stg/keyconfig"
	"github.com/ikuo0/game/ebiten_stg/title"
	"github.com/ikuo0/game/ebiten_stg/stage1"
	"github.com/ikuo0/game/ebiten_stg/global"
	"github.com/ikuo0/game/lib/log"
	"github.com/ikuo0/game/lib/sound"
	"github.com/hajimehoshi/ebiten"
	"os"
	"path"
)

var sceneChanger = changer.NewChanger(changer.SceneList {
	"title":     title.New,
	"keyconfig": keyconfig.New,
	"stage1":    stage1.New,
})

func update(screen *ebiten.Image) error {
	sound.Update()
	return sceneChanger.Play(screen)
}

func main() {
	global.SetRootPath(path.Dir(os.Args[0]))
	log.Start(global.Path().File(`log.txt`))
	sound.Initialize(global.SampleRate, 20, 20)
	sceneChanger.Apply([]string{"title"})
	ebiten.Run(update, 800, 600, 1, "EBITEN STG")
}
