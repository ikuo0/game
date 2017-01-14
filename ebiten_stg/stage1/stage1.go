
package stage1

import (
	"github.com/ikuo0/game/ebiten_stg/action"
	"github.com/ikuo0/game/ebiten_stg/bullet"
	"github.com/ikuo0/game/ebiten_stg/enemy"
	"github.com/ikuo0/game/ebiten_stg/explosion"
	"github.com/ikuo0/game/ebiten_stg/eventid"
	"github.com/ikuo0/game/ebiten_stg/global"
	"github.com/ikuo0/game/ebiten_stg/player"
	"github.com/ikuo0/game/ebiten_stg/instrument"
	"github.com/ikuo0/game/ebiten_stg/result"
	"github.com/ikuo0/game/ebiten_stg/sheld"
	"github.com/ikuo0/game/ebiten_stg/shot"
	"github.com/ikuo0/game/ebiten_stg/world"
	"github.com/ikuo0/game/lib/event"
	"github.com/ikuo0/game/lib/fig"
	"github.com/ikuo0/game/lib/ginput"
	"github.com/ikuo0/game/lib/log"
	"github.com/ikuo0/game/lib/orig"
	"github.com/ikuo0/game/lib/radian"
	"github.com/ikuo0/game/lib/scene"
	"github.com/ikuo0/game/lib/script"
	"github.com/ikuo0/game/lib/sound"
	"github.com/ikuo0/game/lib/ttpl"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image/color"
	"fmt"
	"os"
	"strings"
	"time"
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
//# Sounds
//########################################
type Sounds struct {
	Explosion *sound.Wav
	Hit       *sound.Wav
	Bgm       *sound.Ogg
}

func (me *Sounds) Dispose() {
	me.Explosion.Dispose()
	me.Hit.Dispose()
	me.Bgm.Dispose()
}

func NewSounds() (*Sounds) {
	exp, _ := sound.NewWav("./resource/sound/se_maoudamashii_explosion06.wav")
	hit, _ := sound.NewWav("./resource/sound/se_maoudamashii_battle16.wav")
	bgm, _ := sound.NewOgg("./resource/sound/bgm.ogg")
	return &Sounds {
		Explosion: exp,
		Hit:       hit,
		Bgm:       bgm,
	}
}

//########################################
//# Stage1
//########################################
type Stage1 struct {
	KeyConfig       *global.KeyConfigSt

	Player          action.Objects
	PlayerImage     *ebiten.Image
	PlayerEntity    *player.Player

	Shot            action.Objects
	ShotImage       *ebiten.Image

	Sheld           action.Objects
	SheldImage      *ebiten.Image

	Heli1           action.Objects
	Heli2           action.Objects
	HeliImage       *ebiten.Image

	Aide1           action.Objects
	Aide2           action.Objects
	AideImage       *ebiten.Image

	Boss1           action.Objects
	Boss1Image      *ebiten.Image

	Bullet1         action.Objects
	Bullet1Image    *ebiten.Image

	Explosion1      action.Objects
	Explosion1Image *ebiten.Image

	Vanishing1      action.Objects

	HitImage        *ebiten.Image

	Boss1Entity     *enemy.Boss1

	Template        *ttpl.Ttpl
	Instrument      *instrument.Instrument

	Result          *result.Result

	Score           int
	GameEnd         bool
	SceneEnd        bool

	Sound           *Sounds

	Inner       fig.Rect
	Outer       fig.Rect

	Stack     script.Stack
	Source    script.Input

	Pushed    *ginput.Pushed
	DirectPushed *ginput.DirectPushed
	Debug        bool
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

func (me *Stage1) GetPoint() (fig.Point) {
	return fig.Point{0, 0}
}

func (me *Stage1) Direction() (radian.Radian) {
	return 0
}

func (me *Stage1) IsGameEnd() (bool) {
	return me.GameEnd
}

func (me *Stage1) Update() {
	if me.IsGameEnd() {
		me.Pushed.Update()
		if me.Pushed.Check(ginput.Key1) {
			if me.Result.IsEnd() {
				me.SceneEnd = true
			} else {
				me.Result.Next()
			}
		}
	} else {
		script.Exec(me.Source, &me.Stack, me, me)

		//bits := ginput.Standard()
		bits := ginput.Bits(ginput.Values(), me.KeyConfig.Maps)

		action.SetInput(bits, me.Player)

		action.Update(me, me.Heli1, me.Heli2, me.Aide1, me.Aide2, me.Boss1, me.Bullet1, me.Player, me.Shot, me.Sheld, me.Vanishing1, me.Explosion1)
		action.HitCheck(me.Shot, me.Heli1, me.Heli2, me.Aide1, me.Aide2, me.Boss1)
		action.HitCheck(me.Bullet1, me.Sheld)
		action.HitCheck(me.Player, me.Bullet1, me.Heli1, me.Heli2)
		action.CarryPress(me.Sheld, me.Shot)
		action.InScreen(me.Inner, me.Player)
		action.GoOutside(me.Outer, me.Heli1, me.Heli2, me.Aide1, me.Aide2, me.Boss1, me.Shot, me.Sheld, me.Bullet1)
		//action.Clean(me.Heli1, me.Heli2, me.Aide1, me.Aide2, me.Boss1, me.Player, me.Shot, me.Sheld, me.Vanishing1, me.Explosion1, me.Bullet1)
		me.Heli1 = action.Clean(me.Heli1)
		me.Heli2 = action.Clean(me.Heli2)
		me.Aide1 = action.Clean(me.Aide1)
		me.Aide2 = action.Clean(me.Aide2)
		me.Boss1 = action.Clean(me.Boss1)
		me.Player = action.Clean(me.Player)
		me.Shot = action.Clean(me.Shot)
		me.Sheld = action.Clean(me.Sheld)
		me.Vanishing1 = action.Clean(me.Vanishing1)
		me.Explosion1 = action.Clean(me.Explosion1)
		me.Bullet1 = action.Clean(me.Bullet1)
	}

	me.DirectPushed.Update()
	if me.DirectPushed.Check(ginput.PressValue(ebiten.KeyF1)) {
		me.Debug = !me.Debug
	}

	if !me.Sound.Bgm.IsPlaying() {
		me.Sound.Bgm.Play(time.Millisecond * 500)
	}
	//sound.Update()
}

func (me *Stage1) ObjectCount() (int) {
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

func (me *Stage1) Draw(screen *ebiten.Image) {
	action.DrawImageRota(screen, me.PlayerImage, me.Player)
	action.DrawImageRota(screen, me.AideImage, me.Aide1)
	action.DrawImageRota(screen, me.AideImage, me.Aide2)
	action.DrawImageRota(screen, me.Boss1Image, me.Boss1)
	screen.DrawImage(me.HeliImage, action.DrawOptions(me.Heli1))
	screen.DrawImage(me.HeliImage, action.DrawOptions(me.Heli2))
	screen.DrawImage(me.ShotImage, action.DrawOptions(me.Shot))
	screen.DrawImage(me.SheldImage, action.DrawOptions(me.Sheld))
	screen.DrawImage(me.Bullet1Image, action.DrawOptions(me.Vanishing1))
	screen.DrawImage(me.Explosion1Image, action.DrawOptions(me.Explosion1))
	screen.DrawImage(me.Bullet1Image, action.DrawOptions(me.Bullet1))


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

	if me.Debug {
		hitObjs := action.NewHitObjects(me.Player, me.Shot, me.Sheld, me.Bullet1, me.Heli1, me.Heli2)
		screen.DrawImage(me.HitImage, hitObjs.Options())
	}
}

func (me *Stage1) EventTrigger(id event.Id, argument interface{}, origin orig.Interface) {
	switch id {
		case eventid.Player:
			me.PlayerEntity = player.NewPlayer(world.StartGetPoint())
			me.Player = append(me.Player, me.PlayerEntity)

		case eventid.BgmPlay:
			//me.Sound.Bgm.Play(32)

		case eventid.PlayerDied:
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
				me.EventTrigger(eventid.Player, nil, nil)
			}

		case eventid.StageClear:
			msg := fmt.Sprintf(` Game Clear
  Score: %d
  End`, me.Score)
			me.GameEnd = true
			me.Result = result.New(strings.Split(msg, "\n"))

		case eventid.Shot:
			me.Shot = append(me.Shot, shot.NewShot(origin.GetPoint()))

		case eventid.Sheld:
			if me.Sheld.Len() < 3 {
				me.Sheld = append(me.Sheld, sheld.NewSheld(origin.GetPoint()))
			}

		case eventid.Heli1:
			pt := argument.(fig.Point)
			me.Heli1 = append(me.Heli1, enemy.NewHeli1(pt))

		case eventid.Heli2:
			pt := argument.(fig.Point)
			me.Heli2 = append(me.Heli2, enemy.NewHeli2(pt))

		case eventid.Aide1:
			pt := argument.(fig.Point)
			me.Aide1 = append(me.Aide1, enemy.NewAide(pt))

		case eventid.Aide2:
			pt := argument.(fig.Point)
			me.Aide2 = append(me.Aide2, enemy.NewAide(pt))

		case eventid.Boss1:
			pt := argument.(fig.Point)
			me.Boss1Entity = enemy.NewBoss1(pt); 
			me.Boss1 = append(me.Boss1, me.Boss1Entity)

		case eventid.Bullet1:
			rad := origin.Direction() + argument.(radian.Radian)
			me.Bullet1 = append(me.Bullet1, bullet.NewBullet1(origin.GetPoint(), rad))

		case eventid.Bullet2:
			rad := argument.(radian.Radian)
			me.Bullet1 = append(me.Bullet1, bullet.NewBullet1(origin.GetPoint(), rad))

		case eventid.Explosion1:
			me.Sound.Explosion.Play(0)
			pt := origin.GetPoint()
			if relative, ok := argument.(fig.Point); ok {
				pt.X += relative.X
				pt.Y += relative.Y
			}
			me.Explosion1 = append(me.Explosion1, explosion.NewExplosion1(pt))

		case eventid.Vanishing1:
			pt := origin.GetPoint()
			if relative, ok := argument.(fig.Point); ok {
				pt.X += relative.X
				pt.Y += relative.Y
			}
			me.Vanishing1 = append(me.Vanishing1, explosion.NewVanishing1(pt))

		case eventid.Score:
			n := argument.(int)
			me.Score += n
	}
}

func (me *Stage1) Main(screen *ebiten.Image) (bool) {
	me.Update()
	me.Draw(screen)
	return !me.SceneEnd
}

func (me *Stage1) Dispose() {
	me.Sound.Dispose()
}

func (me *Stage1) ReturnValue() (scene.Parameter) {
	return []string{"title"}
}

func New(args scene.Parameter) (scene.Interface) {
	world.Init()

	hitImage, _ := ebiten.NewImage(1, 1, ebiten.FilterLinear)
	hitImage.Fill(color.RGBA{0xff, 0x00, 0x00, 0x77})

	panelRect := fig.IntRect {500, 0, 800, 600}

	conf := global.KeyConfig()
	if e1 := conf.Load(global.Path().KeyConfig()); e1 != nil {
		log.Log("keyconfig#New#conf.Load error %s", e1.Error())
		conf.Set(ginput.DefaultKeyBoard())
	}

	return &Stage1{
		KeyConfig:       conf,

		Player:          action.NewObjects(),
		PlayerImage:     LoadImage("./resource/image/Player0102.png"),

		Shot:            action.NewObjects(),
		ShotImage:       LoadImage("./resource/image/PlayerSchott.PNG"),

		Sheld:           action.NewObjects(),
		SheldImage:      LoadImage("./resource/image/bomb.png"),

		Heli1:           action.NewObjects(),
		Heli2:           action.NewObjects(),
		HeliImage:       LoadImage("./resource/image/h01.png"),

		Aide1:           action.NewObjects(),
		Aide2:           action.NewObjects(),
		AideImage:       LoadImage("./resource/image/houdai01.PNG"),

		Boss1:           action.NewObjects(),
		Boss1Image:      LoadImage("./resource/image/middleBoss01.png"),

		Bullet1:         action.NewObjects(),
		Bullet1Image:    LoadImage("./resource/image/tekidan01.PNG"),

		Explosion1:      action.NewObjects(),
		Explosion1Image: LoadImage("./resource/image/bakuhatsuM01.png"),

		HitImage:        hitImage,

		Vanishing1:      action.NewObjects(),

		Template:        ttpl.New(PanelTemplate),
		Instrument:      instrument.NewInstrument(panelRect, color.RGBA{0xb2, 0x9a, 0x8e, 0xff}),

		Pushed:          ginput.NewPushed(),
		DirectPushed:    ginput.NewDirectPushed(),

		Sound:       NewSounds(),

		Inner:       fig.Rect{0, 0, 500, 600},
		Outer:       fig.Rect{-64, -64, 564, 664},

		Source:      script.NewSource([]script.Proc {
			/*
			script.NewEventProc(eventid.Heli1, fig.Point{200, 0}),
			script.NewWaitProc(1),
			script.NewJumpProc(1),
			*/

			script.NewWaitProc(20),
			script.NewEventProc(eventid.Boss1, fig.Point{250, 100}),
			script.NewEventProc(eventid.Aide1, fig.Point{100, 64}),
			script.NewEventProc(eventid.Aide1, fig.Point{400, 64}),
			script.NewWaitProc(30),
			script.NewEventProc(eventid.Player, fig.Point{250, 400}),
			script.NewEventProc(eventid.BgmPlay, fig.Point{250, 400}),
			script.NewWaitProc(30),
			//script.NewJumpProc(7),
			script.NewEventProc(eventid.Heli1, fig.Point{100, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(eventid.Heli1, fig.Point{110, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(eventid.Heli1, fig.Point{120, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(eventid.Heli1, fig.Point{130, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(eventid.Heli1, fig.Point{140, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(eventid.Heli1, fig.Point{150, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(eventid.Heli1, fig.Point{160, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(eventid.Heli1, fig.Point{170, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(eventid.Heli1, fig.Point{180, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(eventid.Heli1, fig.Point{190, 0}),
			script.NewWaitProc(30),
			script.NewEventProc(eventid.Heli1, fig.Point{100, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(eventid.Heli1, fig.Point{110, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(eventid.Heli1, fig.Point{120, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(eventid.Heli1, fig.Point{130, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(eventid.Heli1, fig.Point{140, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(eventid.Heli1, fig.Point{150, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(eventid.Heli1, fig.Point{160, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(eventid.Heli1, fig.Point{170, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(eventid.Heli1, fig.Point{180, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(eventid.Heli1, fig.Point{190, 0}),
			script.NewWaitProc(30),
			script.NewEventProc(eventid.Heli1, fig.Point{100, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(eventid.Heli1, fig.Point{110, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(eventid.Heli1, fig.Point{120, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(eventid.Heli1, fig.Point{130, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(eventid.Heli1, fig.Point{140, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(eventid.Heli1, fig.Point{150, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(eventid.Heli1, fig.Point{160, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(eventid.Heli1, fig.Point{170, 0}),
			script.NewWaitProc(3),
			script.NewEventProc(eventid.Heli1, fig.Point{180, 0}),
			script.NewWaitProc(10),
			script.NewEventProc(eventid.Heli1, fig.Point{0, 0}),
			script.NewEventProc(eventid.Heli1, fig.Point{50, 0}),
			script.NewEventProc(eventid.Heli1, fig.Point{100, 0}),
			script.NewEventProc(eventid.Heli1, fig.Point{150, 0}),
			script.NewEventProc(eventid.Heli1, fig.Point{200, 0}),
			script.NewEventProc(eventid.Heli1, fig.Point{250, 0}),
			script.NewEventProc(eventid.Heli1, fig.Point{300, 0}),
			script.NewEventProc(eventid.Heli1, fig.Point{350, 0}),
			script.NewEventProc(eventid.Heli1, fig.Point{400, 0}),
			script.NewEventProc(eventid.Heli1, fig.Point{450, 0}),
			script.NewEventProc(eventid.Heli1, fig.Point{500, 0}),
			script.NewWaitProc(10),
			script.NewEventProc(eventid.Heli1, fig.Point{0, 0}),
			script.NewEventProc(eventid.Heli1, fig.Point{50, 0}),
			script.NewEventProc(eventid.Heli1, fig.Point{100, 0}),
			script.NewEventProc(eventid.Heli1, fig.Point{150, 0}),
			script.NewEventProc(eventid.Heli1, fig.Point{200, 0}),
			script.NewEventProc(eventid.Heli1, fig.Point{250, 0}),
			script.NewEventProc(eventid.Heli1, fig.Point{300, 0}),
			script.NewEventProc(eventid.Heli1, fig.Point{350, 0}),
			script.NewEventProc(eventid.Heli1, fig.Point{400, 0}),
			script.NewEventProc(eventid.Heli1, fig.Point{450, 0}),
			script.NewEventProc(eventid.Heli1, fig.Point{500, 0}),
			script.NewWaitProc(10),
			script.NewEventProc(eventid.Heli1, fig.Point{0, 0}),
			script.NewEventProc(eventid.Heli1, fig.Point{50, 0}),
			script.NewEventProc(eventid.Heli1, fig.Point{100, 0}),
			script.NewEventProc(eventid.Heli1, fig.Point{150, 0}),
			script.NewEventProc(eventid.Heli1, fig.Point{200, 0}),
			script.NewEventProc(eventid.Heli1, fig.Point{250, 0}),
			script.NewEventProc(eventid.Heli1, fig.Point{300, 0}),
			script.NewEventProc(eventid.Heli1, fig.Point{350, 0}),
			script.NewEventProc(eventid.Heli1, fig.Point{400, 0}),
			script.NewEventProc(eventid.Heli1, fig.Point{450, 0}),
			script.NewEventProc(eventid.Heli1, fig.Point{500, 0}),
			script.NewWaitProc(10),
			script.NewEventProc(eventid.Heli1, fig.Point{0, 0}),
			script.NewEventProc(eventid.Heli1, fig.Point{50, 0}),
			script.NewEventProc(eventid.Heli1, fig.Point{100, 0}),
			script.NewEventProc(eventid.Heli1, fig.Point{150, 0}),
			script.NewEventProc(eventid.Heli1, fig.Point{200, 0}),
			script.NewEventProc(eventid.Heli1, fig.Point{250, 0}),
			script.NewEventProc(eventid.Heli1, fig.Point{300, 0}),
			script.NewEventProc(eventid.Heli1, fig.Point{350, 0}),
			script.NewEventProc(eventid.Heli1, fig.Point{400, 0}),
			script.NewEventProc(eventid.Heli1, fig.Point{450, 0}),
			script.NewEventProc(eventid.Heli1, fig.Point{500, 0}),
			script.NewWaitProc(180),
			script.NewEventProc(eventid.Heli2, fig.Point{20, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{40, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{60, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{80, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{100, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{120, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{140, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{160, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{180, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{200, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{220, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{240, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{260, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{280, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{300, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{320, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{340, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{360, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{380, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{400, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{420, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{440, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{460, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{480, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{500, 0}),
			script.NewWaitProc(20),
			script.NewEventProc(eventid.Heli2, fig.Point{20, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{40, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{60, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{80, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{100, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{120, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{140, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{160, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{180, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{200, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{220, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{240, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{260, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{280, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{300, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{320, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{340, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{360, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{380, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{400, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{420, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{440, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{460, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{480, 0}),
			script.NewEventProc(eventid.Heli2, fig.Point{500, 0}),
			script.NewJumpProc(7),

		}),
	}
}
