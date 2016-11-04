package main

import (
	"./enemy"
	"./explosion"
	"./player"
	"github.com/ikuo0/game/lib/event"
	"github.com/ikuo0/game/lib/fig"
	"github.com/ikuo0/game/lib/ginput"
	"github.com/ikuo0/game/lib/orig"
	"github.com/ikuo0/game/lib/radian"
	"github.com/ikuo0/game/lib/script"
	"github.com/ikuo0/game/lib/sprites"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"fmt"
	"os"
	//"github.com/hajimehoshi/ebiten/ebitenutil"
)

//########################################
//# Source
//########################################
type Source struct {
	Procs []script.Proc
	Index int
}

func (me *Source) Read(idx int) (script.Proc) {
	if idx >= len(me.Procs) {
		me.Index = 0
	}
	res := me.Procs[idx]
	me.Index++
	return res
}

func NewSource() (*Source) {
	return &Source {
		Procs: []script.Proc {
			script.NewWaitProc(20),
			script.NewEventProc(event.Enemy2, fig.FloatPoint{100, 64}),
			script.NewEventProc(event.Enemy2, fig.FloatPoint{400, 64}),
			script.NewWaitProc(30),
			script.NewEventProc(event.Player, fig.FloatPoint{250, 400}),
			script.NewWaitProc(30),
			script.NewEventProc(event.Enemy1, fig.FloatPoint{100, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(event.Enemy1, fig.FloatPoint{110, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(event.Enemy1, fig.FloatPoint{120, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(event.Enemy1, fig.FloatPoint{130, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(event.Enemy1, fig.FloatPoint{140, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(event.Enemy1, fig.FloatPoint{150, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(event.Enemy1, fig.FloatPoint{160, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(event.Enemy1, fig.FloatPoint{170, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(event.Enemy1, fig.FloatPoint{180, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(event.Enemy1, fig.FloatPoint{190, 0}),
			script.NewWaitProc(180),
			script.NewJumpProc(5),
		},
	}
}

type Scene struct {
	Player         *player.Players
	PlayerImage    *ebiten.Image

	Shot           *sprites.RotaObjects
	ShotImage      *ebiten.Image

	Heli1          *sprites.RotaObjects
	HeliImage      *ebiten.Image

	Aide           *sprites.RotaObjects
	AideImage      *ebiten.Image

	Explosion1     *sprites.Objects
	Explosion1Image *ebiten.Image

	Inner       fig.Rect
	Outer       fig.Rect

	Stack     script.Stack
	Source    script.Input
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

func (me *Scene) Point() (fig.FloatPoint) {
	return fig.FloatPoint{0, 0}
}

func (me *Scene) Direction() (radian.Radian) {
	return 0
}

func (me *Scene) Update() {
	script.Exec(me.Source, &me.Stack, me, me)

	bits := ginput.Standard()

	sprites.SetInput(bits, me.Player)
	sprites.Update(me, me.Heli1, me.Aide, me.Player, me.Shot, me.Explosion1)
	sprites.HitCheck(me.Shot, me.Heli1)
	sprites.InScreen(me.Inner, me.Player)
	sprites.GoOutside(me.Outer, me.Heli1, me.Aide, me.Shot)
	sprites.Clean(me.Heli1, me.Aide, me.Player, me.Shot, me.Explosion1)
}

func (me *Scene) Draw(screen *ebiten.Image) {
	sprites.ExDraw(screen, me.PlayerImage, me.Player)
	sprites.ExDraw(screen, me.HeliImage, me.Heli1)
	sprites.ExDraw(screen, me.AideImage, me.Aide)
	sprites.ExDraw(screen, me.ShotImage, me.Shot)
	screen.DrawImage(me.Explosion1Image, me.Explosion1.Options())
}

func (me *Scene) EventTrigger(id event.Id, argument interface{}, origin orig.Interface) {
	switch id {
		case event.Player:
			pt := argument.(fig.FloatPoint)
			me.Player.Occure(player.NewPlayer(pt))

		case event.Shot:
			//pt := argument.(fig.FloatPoint)
			me.Shot.Occure(player.NewShot(origin.Point()))

		case event.Enemy1:
			pt := argument.(fig.FloatPoint)
			me.Heli1.Occure(enemy.NewHeli1(pt))

		case event.Enemy2:
			pt := argument.(fig.FloatPoint)
			me.Aide.Occure(enemy.NewAide(pt))

		case event.Explosion1:
			pt := origin.Point()
			me.Explosion1.Occure(explosion.NewExplosion1(pt))
	}
}

func NewScene() *Scene {
	return &Scene{
		Player:      player.NewPlayers(),
		PlayerImage: LoadImage("./resource/image/Player0102.png"),

		Shot:        sprites.NewRotaObjects(),
		ShotImage:   LoadImage("./resource/image/PlayerSchott.PNG"),

		Heli1:       sprites.NewRotaObjects(),
		HeliImage:   LoadImage("./resource/image/h01.png"),

		Aide:        sprites.NewRotaObjects(),
		AideImage:   LoadImage("./resource/image/houdai01.PNG"),

		Explosion1:  sprites.NewObjects(),
		Explosion1Image: LoadImage("./resource/image/bakuhatsuM01.png"),

		Inner:       fig.Rect{0, 0, 500, 600},
		Outer:       fig.Rect{-64, -64, 564, 664},

		Source:      NewSource(),
	}
}

var scene = NewScene()

type MainFunc func(*ebiten.Image) (error)
var currentFn MainFunc

func initialize(screen *ebiten.Image) error {
	ginput.Initialize()

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
