package main

import (
	"./player"
	"github.com/ikuo0/game/lib/event"
	"github.com/ikuo0/game/lib/fig"
	"github.com/ikuo0/game/lib/ginput"
	"github.com/ikuo0/game/lib/orig"
	"github.com/ikuo0/game/lib/radian"
	"github.com/ikuo0/game/lib/sprites"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"fmt"
	"os"
	//"github.com/hajimehoshi/ebiten/ebitenutil"
)

//########################################
//# Scene
//########################################
type Scene struct {
	Player      *player.Players
	PlayerImage *ebiten.Image

	Inner       fig.Rect
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

func (me *Scene) GetPoint() (fig.Point) {
	return fig.Point{0, 0}
}

func (me *Scene) Direction() (radian.Radian) {
	return 0
}

func (me *Scene) EventTrigger(id event.Id, argument interface{}, origin orig.Interface) {
}

func (me *Scene) Update() {
	bits := ginput.Standard()

	sprites.SetInput(bits, me.Player)
	sprites.Update(me, me.Player)
	sprites.InScreen(me.Inner, me.Player)
}

func (me *Scene) Draw(screen *ebiten.Image) {
	sprites.ExDraw(screen, me.PlayerImage, me.Player)
}

func NewScene() *Scene {
	x := player.NewPlayers()
	x.Occure(player.NewPlayer(fig.Point{100, 100}))
	return &Scene{
		Player:          x,
		PlayerImage:     LoadImage("./resource/image/Player0102.png"),
		Inner:       fig.Rect{0, 0, 500, 600},
	}
}

var scene *Scene

type MainFunc func(*ebiten.Image) (error)
var currentFn MainFunc

func initialize(screen *ebiten.Image) error {
	ginput.Initialize()
	scene = NewScene()

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
	currentFn = initialize
	ebiten.Run(update, 800, 600, 1, "Hello world!")
}
