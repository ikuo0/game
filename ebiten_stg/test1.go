package main

import (
	"./bullet"
	"./enemy"
	"./explosion"
	"./player"
	"./instrument"
	"./event"
	"./result"
	"./sprites"
	"./script"
	"./world"
	"./lib/fig"
	"./lib/ginput"
	"./lib/orig"
	"./lib/radian"
	"./lib/sound"
	"./lib/ttpl"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image/color"
	"fmt"
	"os"
	"strings"
	//"github.com/hajimehoshi/ebiten/ebitenutil"
)

const PanelTemplate = `
 Frame #FrameCount#

 Objects #ObjectCount#

 Boss #BossEndurance#

 Player #PlayerEndurance#

 Score #Score#
`

//########################################
//# Source
//########################################
type Sounds struct {
	Explosion *sound.Wav
	Hit       *sound.Wav
}

func NewSounds() (*Sounds) {
	exp, _ := sound.NewWav("./resource/sound/se_maoudamashii_explosion06.wav")
	hit, _ := sound.NewWav("./resource/sound/se_maoudamashii_battle16.wav")
	return &Sounds {
		Explosion: exp,
		Hit:       hit,
	}
}

//########################################
//# Scene
//########################################
type Scene struct {
	Player          *player.Players
	PlayerImage     *ebiten.Image
	PlayerEntity    *player.Player

	Shot            *sprites.Objects
	ShotImage       *ebiten.Image

	Sheld           *player.Shields
	SheldImage      *ebiten.Image

	Heli1           *sprites.Objects
	Heli2           *sprites.Objects
	HeliImage       *ebiten.Image

	Aide1           *sprites.RotaObjects
	Aide2           *sprites.RotaObjects
	AideImage       *ebiten.Image

	Boss1           *sprites.RotaObjects
	Boss1Image      *ebiten.Image

	Bullet1         *sprites.Objects
	Bullet1Image    *ebiten.Image

	Explosion1      *sprites.Objects
	Explosion1Image *ebiten.Image

	Vanishing1      *sprites.Objects

	HitImage        *ebiten.Image

	Boss1Entity     *enemy.Boss1

	Template        *ttpl.Ttpl
	Instrument      *instrument.Instrument

	Result          *result.Result

	Score           int
	GameEnd         bool

	Sound           *Sounds

	Inner       fig.Rect
	Outer       fig.Rect

	Stack     script.Stack
	Source    script.Input

	Pushed    *ginput.Pushed
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

func (me *Scene) IsGameEnd() (bool) {
	return me.GameEnd
}

func (me *Scene) Update() {
	if me.IsGameEnd() {
		me.Pushed.Update()
		if me.Pushed.Check(ginput.Key1) {
			if me.Result.IsEnd() {
				os.Exit(0)
			} else {
				me.Result.Next()
			}
		}
	} else {
		script.Exec(me.Source, &me.Stack, me, me)

		bits := ginput.Standard()

		sprites.SetInput(bits, me.Player)
		sprites.Update(me, me.Heli1, me.Heli2, me.Aide1, me.Aide2, me.Boss1, me.Bullet1, me.Player, me.Shot, me.Sheld, me.Vanishing1, me.Explosion1)
		sprites.HitCheck(me.Shot, me.Heli1)
		sprites.HitCheck(me.Shot, me.Heli2)
		sprites.HitCheck(me.Shot, me.Aide1)
		sprites.HitCheck(me.Shot, me.Aide2)
		sprites.HitCheck(me.Shot, me.Boss1)
		sprites.HitCheck(me.Bullet1, me.Sheld)
		sprites.HitCheck(me.Player, me.Bullet1)
		sprites.HitCheck(me.Player, me.Heli1)
		sprites.HitCheck(me.Player, me.Heli2)
		sprites.CarryPress(me.Sheld, me.Shot)
		sprites.InScreen(me.Inner, me.Player)
		sprites.GoOutside(me.Outer, me.Heli1, me.Heli2, me.Aide1, me.Aide2, me.Boss1, me.Shot, me.Sheld, me.Bullet1)
		sprites.Clean(me.Heli1, me.Heli2, me.Aide1, me.Aide2, me.Boss1, me.Player, me.Shot, me.Sheld, me.Vanishing1, me.Explosion1, me.Bullet1)
	}

	sound.Update()
}

func (me *Scene) ObjectCount() (int) {
	res := 0
	res += me.Player.Len()
	res += me.Shot.Len()
	res += me.Sheld.Len()
	res += me.Heli1.Len()
	res += me.Heli2.Len()
	res += me.Aide1.Len()
	res += me.Aide2.Len()
	res += me.Boss1.Len()
	res += me.Bullet1.Len()
	res += me.Explosion1.Len()
	res += me.Vanishing1.Len()
	return res
}

func (me *Scene) Draw(screen *ebiten.Image) {
	sprites.ExDraw(screen, me.PlayerImage, me.Player)
	sprites.ExDraw(screen, me.AideImage, me.Aide1)
	sprites.ExDraw(screen, me.AideImage, me.Aide2)
	sprites.ExDraw(screen, me.Boss1Image, me.Boss1)
	screen.DrawImage(me.HeliImage, me.Heli1.Options())
	screen.DrawImage(me.HeliImage, me.Heli2.Options())
	screen.DrawImage(me.ShotImage, me.Shot.Options())
	screen.DrawImage(me.SheldImage, me.Sheld.Options())
	screen.DrawImage(me.Bullet1Image, me.Vanishing1.Options())
	screen.DrawImage(me.Explosion1Image, me.Explosion1.Options())
	screen.DrawImage(me.Bullet1Image, me.Bullet1.Options())


	me.Template.SetFloat("FrameCount", ebiten.CurrentFPS())
	me.Template.SetInt("ObjectCount", me.ObjectCount())
	if me.Boss1Entity != nil {
		me.Template.SetInt("BossEndurance", me.Boss1Entity.Endurance)
	} else {
		me.Template.SetInt("BossEndurance", 0)
	}
	if me.PlayerEntity != nil {
		me.Template.SetInt("PlayerEndurance", me.PlayerEntity.Endurance)
	} else {
		me.Template.SetInt("PlayerEndurance", 0)
	}

	me.Template.SetInt("Score", me.Score)

	me.Instrument.UpdateText(me.Template.Text())
	screen.DrawImage(me.Instrument.Image(), me.Instrument.Options())

	if me.IsGameEnd() {
		me.Result.Update(screen)
	}

	if true {
		hitObjs := sprites.NewHitObjects(me.Player, me.Sheld, me.Bullet1, me.Heli1)
		screen.DrawImage(me.HitImage, hitObjs.Options())
	}
}

func (me *Scene) EventTrigger(id event.Id, argument interface{}, origin orig.Interface) {
	switch id {
		case event.Player:
			me.PlayerEntity = player.NewPlayer(world.StartGetPoint())
			me.Player.Occure(me.PlayerEntity)

		case event.PlayerDied:
			cnt := world.GetPlayerCount()
			cnt -= 1
			world.SetPlayerCount(cnt)
			if cnt < 1 {
				msg := fmt.Sprintf(` GameOver
  Score: %d
  End`, me.Score)
				me.GameEnd = true
				me.Result = result.New(strings.Split(msg, "\n"))
			} else {
				me.EventTrigger(event.Player, nil, nil)
			}

		case event.StageClear:
			msg := fmt.Sprintf(` Game Clear
  Score: %d
  End`, me.Score)
			me.GameEnd = true
			me.Result = result.New(strings.Split(msg, "\n"))

		case event.Shot:
			//pt := argument.(fig.Point)
			me.Shot.Occure(player.NewShot(origin.GetPoint()))

		case event.Sheld:
			//pt := argument.(fig.Point)
			if me.Sheld.Len() < 3 {
				me.Sheld.Occure(player.NewSheld(origin.GetPoint()))
			}

		case event.Heli1:
			pt := argument.(fig.Point)
			me.Heli1.Occure(enemy.NewHeli1(pt))

		case event.Heli2:
			pt := argument.(fig.Point)
			me.Heli2.Occure(enemy.NewHeli2(pt))

		case event.Aide1:
			pt := argument.(fig.Point)
			me.Aide1.Occure(enemy.NewAide(pt))

		case event.Aide2:
			pt := argument.(fig.Point)
			me.Aide2.Occure(enemy.NewAide(pt))

		case event.Boss1:
			pt := argument.(fig.Point)
			me.Boss1Entity =enemy.NewBoss1(pt); 
			me.Boss1.Occure(me.Boss1Entity)

		case event.Bullet1:
			rad := origin.Direction() + argument.(radian.Radian)
			me.Bullet1.Occure(bullet.NewBullet1(origin.GetPoint(), rad))

		case event.Explosion1:
			me.Sound.Explosion.Play()
			pt := origin.GetPoint()
			if relative, ok := argument.(fig.Point); ok {
				pt.X += relative.X
				pt.Y += relative.Y
			}
			me.Explosion1.Occure(explosion.NewExplosion1(pt))

		case event.Vanishing1:
			pt := origin.GetPoint()
			if relative, ok := argument.(fig.Point); ok {
				pt.X += relative.X
				pt.Y += relative.Y
			}
			me.Vanishing1.Occure(explosion.NewVanishing1(pt))

		case event.Score:
			n := argument.(int)
			me.Score += n
	}
}

func NewScene() *Scene {
	hitImage, _ := ebiten.NewImage(1, 1, ebiten.FilterLinear)
	hitImage.Fill(color.RGBA{0xff, 0x00, 0x00, 0x77})

	panelRect := fig.Rect {500, 0, 800, 600}

	return &Scene{
		Player:          player.NewPlayers(),
		PlayerImage:     LoadImage("./resource/image/Player0102.png"),

		Shot:            sprites.NewObjects(),
		ShotImage:       LoadImage("./resource/image/PlayerSchott.PNG"),

		Sheld:           player.NewShelds(),
		SheldImage:      LoadImage("./resource/image/bomb.png"),

		Heli1:           sprites.NewObjects(),
		Heli2:           sprites.NewObjects(),
		HeliImage:       LoadImage("./resource/image/h01.png"),

		Aide1:           sprites.NewRotaObjects(),
		Aide2:           sprites.NewRotaObjects(),
		AideImage:       LoadImage("./resource/image/houdai01.PNG"),

		Boss1:           sprites.NewRotaObjects(),
		Boss1Image:      LoadImage("./resource/image/middleBoss01.png"),

		Bullet1:         sprites.NewObjects(),
		Bullet1Image:    LoadImage("./resource/image/tekidan01.PNG"),

		Explosion1:      sprites.NewObjects(),
		Explosion1Image: LoadImage("./resource/image/bakuhatsuM01.png"),

		HitImage:        hitImage,

		Vanishing1:      sprites.NewObjects(),

		Template:        ttpl.New(PanelTemplate),
		Instrument:      instrument.NewInstrument(panelRect, color.RGBA{0xb2, 0x9a, 0x8e, 0xff}),

		Pushed:          ginput.NewPushed(),

		Sound:       NewSounds(),

		Inner:       fig.Rect{0, 0, 500, 600},
		Outer:       fig.Rect{-64, -64, 564, 664},

		Source:      script.NewSource([]script.Proc {
			script.NewWaitProc(20),
			script.NewEventProc(event.Boss1, fig.Point{250, 100}),
			script.NewEventProc(event.Aide1, fig.Point{100, 64}),
			script.NewEventProc(event.Aide1, fig.Point{400, 64}),
			script.NewWaitProc(30),
			script.NewEventProc(event.Player, fig.Point{250, 400}),
			script.NewWaitProc(30),
			script.NewEventProc(event.Heli1, fig.Point{100, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(event.Heli1, fig.Point{110, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(event.Heli1, fig.Point{120, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(event.Heli1, fig.Point{130, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(event.Heli1, fig.Point{140, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(event.Heli1, fig.Point{150, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(event.Heli1, fig.Point{160, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(event.Heli1, fig.Point{170, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(event.Heli1, fig.Point{180, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(event.Heli1, fig.Point{190, 0}),
			script.NewWaitProc(30),
			script.NewEventProc(event.Heli1, fig.Point{100, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(event.Heli1, fig.Point{110, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(event.Heli1, fig.Point{120, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(event.Heli1, fig.Point{130, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(event.Heli1, fig.Point{140, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(event.Heli1, fig.Point{150, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(event.Heli1, fig.Point{160, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(event.Heli1, fig.Point{170, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(event.Heli1, fig.Point{180, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(event.Heli1, fig.Point{190, 0}),
			script.NewWaitProc(30),
			script.NewEventProc(event.Heli1, fig.Point{100, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(event.Heli1, fig.Point{110, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(event.Heli1, fig.Point{120, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(event.Heli1, fig.Point{130, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(event.Heli1, fig.Point{140, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(event.Heli1, fig.Point{150, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(event.Heli1, fig.Point{160, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(event.Heli1, fig.Point{170, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(event.Heli1, fig.Point{180, 0}),
			script.NewWaitProc(10),
			script.NewEventProc(event.Heli1, fig.Point{0, 0}),
			script.NewEventProc(event.Heli1, fig.Point{50, 0}),
			script.NewEventProc(event.Heli1, fig.Point{100, 0}),
			script.NewEventProc(event.Heli1, fig.Point{150, 0}),
			script.NewEventProc(event.Heli1, fig.Point{200, 0}),
			script.NewEventProc(event.Heli1, fig.Point{250, 0}),
			script.NewEventProc(event.Heli1, fig.Point{300, 0}),
			script.NewEventProc(event.Heli1, fig.Point{350, 0}),
			script.NewEventProc(event.Heli1, fig.Point{400, 0}),
			script.NewEventProc(event.Heli1, fig.Point{450, 0}),
			script.NewEventProc(event.Heli1, fig.Point{500, 0}),
			script.NewWaitProc(10),
			script.NewEventProc(event.Heli1, fig.Point{0, 0}),
			script.NewEventProc(event.Heli1, fig.Point{50, 0}),
			script.NewEventProc(event.Heli1, fig.Point{100, 0}),
			script.NewEventProc(event.Heli1, fig.Point{150, 0}),
			script.NewEventProc(event.Heli1, fig.Point{200, 0}),
			script.NewEventProc(event.Heli1, fig.Point{250, 0}),
			script.NewEventProc(event.Heli1, fig.Point{300, 0}),
			script.NewEventProc(event.Heli1, fig.Point{350, 0}),
			script.NewEventProc(event.Heli1, fig.Point{400, 0}),
			script.NewEventProc(event.Heli1, fig.Point{450, 0}),
			script.NewEventProc(event.Heli1, fig.Point{500, 0}),
			script.NewWaitProc(10),
			script.NewEventProc(event.Heli1, fig.Point{0, 0}),
			script.NewEventProc(event.Heli1, fig.Point{50, 0}),
			script.NewEventProc(event.Heli1, fig.Point{100, 0}),
			script.NewEventProc(event.Heli1, fig.Point{150, 0}),
			script.NewEventProc(event.Heli1, fig.Point{200, 0}),
			script.NewEventProc(event.Heli1, fig.Point{250, 0}),
			script.NewEventProc(event.Heli1, fig.Point{300, 0}),
			script.NewEventProc(event.Heli1, fig.Point{350, 0}),
			script.NewEventProc(event.Heli1, fig.Point{400, 0}),
			script.NewEventProc(event.Heli1, fig.Point{450, 0}),
			script.NewEventProc(event.Heli1, fig.Point{500, 0}),
			script.NewWaitProc(10),
			script.NewEventProc(event.Heli1, fig.Point{0, 0}),
			script.NewEventProc(event.Heli1, fig.Point{50, 0}),
			script.NewEventProc(event.Heli1, fig.Point{100, 0}),
			script.NewEventProc(event.Heli1, fig.Point{150, 0}),
			script.NewEventProc(event.Heli1, fig.Point{200, 0}),
			script.NewEventProc(event.Heli1, fig.Point{250, 0}),
			script.NewEventProc(event.Heli1, fig.Point{300, 0}),
			script.NewEventProc(event.Heli1, fig.Point{350, 0}),
			script.NewEventProc(event.Heli1, fig.Point{400, 0}),
			script.NewEventProc(event.Heli1, fig.Point{450, 0}),
			script.NewEventProc(event.Heli1, fig.Point{500, 0}),
			script.NewWaitProc(180),
			script.NewEventProc(event.Heli2, fig.Point{20, 0}),
			script.NewEventProc(event.Heli2, fig.Point{40, 0}),
			script.NewEventProc(event.Heli2, fig.Point{60, 0}),
			script.NewEventProc(event.Heli2, fig.Point{80, 0}),
			script.NewEventProc(event.Heli2, fig.Point{100, 0}),
			script.NewEventProc(event.Heli2, fig.Point{120, 0}),
			script.NewEventProc(event.Heli2, fig.Point{140, 0}),
			script.NewEventProc(event.Heli2, fig.Point{160, 0}),
			script.NewEventProc(event.Heli2, fig.Point{180, 0}),
			script.NewEventProc(event.Heli2, fig.Point{200, 0}),
			script.NewEventProc(event.Heli2, fig.Point{220, 0}),
			script.NewEventProc(event.Heli2, fig.Point{240, 0}),
			script.NewEventProc(event.Heli2, fig.Point{260, 0}),
			script.NewEventProc(event.Heli2, fig.Point{280, 0}),
			script.NewEventProc(event.Heli2, fig.Point{300, 0}),
			script.NewEventProc(event.Heli2, fig.Point{320, 0}),
			script.NewEventProc(event.Heli2, fig.Point{340, 0}),
			script.NewEventProc(event.Heli2, fig.Point{360, 0}),
			script.NewEventProc(event.Heli2, fig.Point{380, 0}),
			script.NewEventProc(event.Heli2, fig.Point{400, 0}),
			script.NewEventProc(event.Heli2, fig.Point{420, 0}),
			script.NewEventProc(event.Heli2, fig.Point{440, 0}),
			script.NewEventProc(event.Heli2, fig.Point{460, 0}),
			script.NewEventProc(event.Heli2, fig.Point{480, 0}),
			script.NewEventProc(event.Heli2, fig.Point{500, 0}),
			script.NewWaitProc(20),
			script.NewEventProc(event.Heli2, fig.Point{20, 0}),
			script.NewEventProc(event.Heli2, fig.Point{40, 0}),
			script.NewEventProc(event.Heli2, fig.Point{60, 0}),
			script.NewEventProc(event.Heli2, fig.Point{80, 0}),
			script.NewEventProc(event.Heli2, fig.Point{100, 0}),
			script.NewEventProc(event.Heli2, fig.Point{120, 0}),
			script.NewEventProc(event.Heli2, fig.Point{140, 0}),
			script.NewEventProc(event.Heli2, fig.Point{160, 0}),
			script.NewEventProc(event.Heli2, fig.Point{180, 0}),
			script.NewEventProc(event.Heli2, fig.Point{200, 0}),
			script.NewEventProc(event.Heli2, fig.Point{220, 0}),
			script.NewEventProc(event.Heli2, fig.Point{240, 0}),
			script.NewEventProc(event.Heli2, fig.Point{260, 0}),
			script.NewEventProc(event.Heli2, fig.Point{280, 0}),
			script.NewEventProc(event.Heli2, fig.Point{300, 0}),
			script.NewEventProc(event.Heli2, fig.Point{320, 0}),
			script.NewEventProc(event.Heli2, fig.Point{340, 0}),
			script.NewEventProc(event.Heli2, fig.Point{360, 0}),
			script.NewEventProc(event.Heli2, fig.Point{380, 0}),
			script.NewEventProc(event.Heli2, fig.Point{400, 0}),
			script.NewEventProc(event.Heli2, fig.Point{420, 0}),
			script.NewEventProc(event.Heli2, fig.Point{440, 0}),
			script.NewEventProc(event.Heli2, fig.Point{460, 0}),
			script.NewEventProc(event.Heli2, fig.Point{480, 0}),
			script.NewEventProc(event.Heli2, fig.Point{500, 0}),
			script.NewJumpProc(6),
		}),
	}
}

var scene *Scene

type MainFunc func(*ebiten.Image) (error)
var currentFn MainFunc

func initialize(screen *ebiten.Image) error {
	ginput.Initialize()
	sound.NewContext(44100)
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
