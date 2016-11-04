package main

import (
	"./enemy"
	"./lib/fig"
	"./lib/sprites"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"fmt"
	"os"
	//"github.com/hajimehoshi/ebiten/ebitenutil"
)

type Scene struct {
	Heli1     *enemy.Enemies
	HeliImage *ebiten.Image

	Aide      *enemy.Enemies
	AideImage *ebiten.Image

	Counter     int
	Outer       fig.Rect
}

func LoadImage(fileName string) *ebiten.Image {
	if img, _, err := ebitenutil.NewImageFromFile(fileName, ebiten.FilterNearest); err != nil {
		fmt.Println(err)
		os.Exit(1)
		return nil
	} else {
		return img
	}
}

func NewScene() *Scene {
	return &Scene{
		Heli1:     enemy.NewEnemies(),
		HeliImage: LoadImage("./resource/image/h01.png"),

		Aide:      enemy.NewEnemies(),
		AideImage: LoadImage("./resource/image/houdai01.PNG"),

		Outer:     fig.Rect{0, 0, 800, 600},
	}
}

func (me *Scene) Update() {
	if me.Counter == 0 {
		me.Aide.Occure(enemy.NewAide(100, 80))
		me.Aide.Occure(enemy.NewAide(300, 80))
	}

	if (me.Counter % 20) == 0 {
		me.Heli1.Occure(enemy.NewHeli1(400, 0))
	}
	me.Counter++

	sprites.Update(me.Heli1, me.Aide)
	sprites.OutScreen(me.Outer, me.Heli1, me.Aide)
	sprites.Clean(me.Heli1, me.Aide)
}

func (me *Scene) Draw(screen *ebiten.Image) {
	sprites.ExDraw(screen, me.HeliImage, me.Heli1)
	sprites.ExDraw(screen, me.AideImage, me.Aide)
}

var scene = NewScene()

type MainFunc func(*ebiten.Image) (error)
var currentFn MainFunc

func initialize(screen *ebiten.Image) error {
	currentFn = mainLoop
	return nil
}

func mainLoop(screen *ebiten.Image) error {
	scene.Update()
	scene.Draw(screen)
	return nil
}

func update(screen *ebiten.Image) error {
	currentFn(screen)
	return nil
}

func main() {
	//ginput.Initialize()
	currentFn = initialize
	ebiten.Run(update, 800, 600, 1, "Hello world!")
}
