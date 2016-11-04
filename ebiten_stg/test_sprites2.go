package main

import (
	"./lib/fig"
	"./lib/move"
	"./lib/radian"
	"./lib/orig"
	"./lib/sprites"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"fmt"
	"math"
	"os"
	//"github.com/hajimehoshi/ebiten/ebitenutil"
)


type AnimeFrame struct {
	Duration int
	Counter int
}

func (me *AnimeFrame) Tick() (bool) {
	if me.Counter >= me.Duration {
		return true
	} else {
		me.Counter++
		return false
	}
}
func (me *AnimeFrame) Reset() {
	me.Counter = 0
}

type AnimeFrames struct {
	Frames []AnimeFrame
	Count  int
}
func (me *AnimeFrames) Update() {
	if me.Frames[me.Count].Tick() {
		me.Frames[me.Count].Reset()
		me.Count++
		if me.Count >= len(me.Frames) {
			me.Count = 0
		}
	}
}

func (me *AnimeFrames) Index() (int) {
	return me.Count
}

func NewAnimeFrame(frames ...int) (*AnimeFrames) {
	ary := []AnimeFrame{}
	for _, v := range frames {
		ary = append(ary, AnimeFrame {
			Duration: v,
		})
	}
	return &AnimeFrames {
		Frames: ary,
	}
}

type Enemy struct {
	fig.FloatPoint
	V *move.Vector
	Anime *AnimeFrames
	Vanished bool
}

func (me *Enemy) Point() (fig.FloatPoint) {
	return me.FloatPoint
}

func (me *Enemy) Direction() (radian.Radian) {
	return me.V.Radian
}

func (me *Enemy) Update() {
	me.V.Accel()
	p := me.V.Power()
	me.X += p.X
	me.Y += p.Y
	me.Anime.Update()

	//me.V.Radian += radian.FromDeg(1)
	me.V.Radian = me.V.Radian.TurnLeft(1)
	/*
	me.X += math.Cos(90 * math.Pi / 180)
	me.Y += 4

	//me.V.Radian += radian.FromDeg(1)
	*/
}

func (me *Enemy) Vanish() {
	me.Vanished = true
}

func (me *Enemy) IsVanish() (bool) {
	return me.Vanished
}

func (me *Enemy) AnimeIndex() (int) {
	return me.Anime.Index()
}

func NewEnemy(x, y float64) (*Enemy) {
	return &Enemy {
		FloatPoint: fig.FloatPoint{x, y},
		//V:          move.NewAccel(radian.Down(), 1, 0.4, 8),
		V:          move.NewAccel((90 * math.Pi) / 180, 1, 0.4, 8),
		Anime:      NewAnimeFrame(4, 4),
	}
}

type Enemies struct {
	Objs     []*Enemy
	SrcRects []fig.Rect
}

func (me *Enemies) Len() (int) {
	return len(me.Objs)
}
func (me *Enemies) Dst(i int) (x0, y0, x1, y1 int) {
	o := me.Objs[i]
	x, y := int(o.X) - 32, int(o.Y) - 24
	return x, y, x + 64, y + 48
}
func (me *Enemies) Src(i int) (x0, y0, x1, y1 int) {
	rect := me.SrcRects[me.Objs[i].AnimeIndex()]
	return rect.Left, rect.Top, rect.Right, rect.Bottom
}
func (me *Enemies) Update(i int) {
	me.Objs[i].Update()
}
func (me *Enemies) Origin(i int) (orig.Interface) {
	return me.Objs[i]
}
func (me *Enemies) Vanish(i int) {
	me.Objs[i].Vanish()
}
func (me *Enemies) Clean(i int) {
	newObjs := []*Enemy{}
	for _, v := range me.Objs {
		if !v.IsVanish() {
			newObjs = append(newObjs, v)
		}
	}
	me.Objs = newObjs
}
func (me *Enemies) Options() (*ebiten.DrawImageOptions) {
	return &ebiten.DrawImageOptions {
		ImageParts: me,
	}
}

func (me *Enemies) DrawOption(i int) (*ebiten.DrawImageOptions) {
	sx0, sy0, sx1, sy1 := me.Src(i)
	dx0, dy0, dx1, dy1 := me.Dst(i)
	opt := ebiten.DrawImageOptions {
		ImageParts: sprites.NewOneSprites(sx0, sy0, sx1, sy1, dx0, dy0, dx1, dy1),
	}

	o := me.Objs[i]
	opt.GeoM.Translate(-o.X, -o.Y)
	opt.GeoM.Rotate(float64(o.V.Radian))
	opt.GeoM.Translate(o.X, o.Y)
	return &opt
}

func (me *Enemies) Occure(x, y int) {
	me.Objs = append(me.Objs, NewEnemy(float64(x), float64(y)))
}

func NewEnemies() (*Enemies) {
	return &Enemies {
		SrcRects: []fig.Rect {
			{
				0,
				122,
				0 + 64,
				122 + 48,
			},
			{
				64,
				122,
				64 + 64,
				122 + 48,
			},
		},
	}
}

type Scene struct {
	Enemies     *Enemies
	Enemy1Image *ebiten.Image
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
		Enemies:     NewEnemies(),
		Enemy1Image: LoadImage("./resource/image/h01.png"),
		Outer:       fig.Rect{0, 0, 800, 600},
	}
}

func (me *Scene) Update() {
	if me.Counter == 10 {
		me.Enemies.Occure(100, 100)
		me.Counter = 0
	}
	me.Counter++

	sprites.Update(me.Enemies)
	sprites.OutScreen(me.Outer, me.Enemies)
	sprites.Clean(me.Enemies)
}

func (me *Scene) Draw(screen *ebiten.Image) {
	sprites.ExDraw(screen, me.Enemy1Image, me.Enemies)
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
