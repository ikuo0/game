
package stage

import (
	"github.com/ikuo0/game/ebiten_act/block"
	"github.com/ikuo0/game/ebiten_act/enemy"
	"github.com/ikuo0/game/ebiten_act/eventid"
	"github.com/ikuo0/game/ebiten_act/explosion"
	"github.com/ikuo0/game/ebiten_act/funcs"
	"github.com/ikuo0/game/ebiten_act/global"
	"github.com/ikuo0/game/ebiten_act/instrument"
	"github.com/ikuo0/game/ebiten_act/player"
	"github.com/ikuo0/game/ebiten_act/result"
	"github.com/ikuo0/game/ebiten_act/shot"
	"github.com/ikuo0/game/ebiten_act/vortex"
	"github.com/ikuo0/game/lib/action"
	"github.com/ikuo0/game/lib/event"
	"github.com/ikuo0/game/lib/fig"
	"github.com/ikuo0/game/lib/ginput"
	"github.com/ikuo0/game/lib/log"
	"github.com/ikuo0/game/lib/orig"
	"github.com/ikuo0/game/lib/radian"
	"github.com/ikuo0/game/lib/scene"
	"github.com/ikuo0/game/lib/script"
	"github.com/ikuo0/game/lib/sprites"
	"github.com/ikuo0/game/lib/ttpl"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image/color"
	"fmt"
	"strconv"
	"strings"
)

const LastStage = 3
const PanelTemplate = `Frame #Frame# Total #TotalFrame# Objects #ObjectCount# #FrameCount#fps`

//########################################
//# Stage1
//########################################
type Stage1 struct {
	KeyConfig        *global.KeyConfigSt

	GameStatus       global.GameStatus
	Frame            int

	Player           *player.Objects
	PlayerImage      *ebiten.Image
	PlayerEntity     *player.Player

	Shot             *sprites.Objects
	ShotImage        *ebiten.Image

	Block            *sprites.Objects
	BlockImage       *ebiten.Image

	OccureBlock      *sprites.Objects
	OccureBlockImage *ebiten.Image

	Enemy            *enemy.Objects
	EnemyImage       *ebiten.Image

	Explosion1        *sprites.Objects
	Explosion1Image   *ebiten.Image

	Vortex           *sprites.Objects
	VortexImage      *ebiten.Image

	HitImage         *ebiten.Image

	Template         *ttpl.Ttpl
	Instrument       *instrument.Instrument

	Result           *result.Result

	Restart          bool
	GameEnd          bool
	SceneEnd         bool

	Inner            fig.Rect
	Outer            fig.Rect

	Stack            script.Stack
	Source           script.Input

	Pushed           *ginput.Pushed
	DirectPushed     *ginput.DirectPushed
	Debug            bool
}

func LoadImage(fileName string) *ebiten.Image {
	if img, _, err := ebitenutil.NewImageFromFile(fileName, ebiten.FilterNearest); err != nil {
		log.Exit("画像読み込みエラー: %s", err.Error())
		return nil
	} else {
		return img
	}
}

func (me *Stage1) Point() (fig.FloatPoint) {
	return fig.FloatPoint{0, 0}
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
		me.Frame++
		script.Exec(me.Source, &me.Stack, me, me)

		bits := ginput.Bits(ginput.Values(), me.KeyConfig.Maps)
		if bits.And(ginput.Key3) {
			me.Restart = true
			me.SceneEnd = true
		}

		action.SetInput(bits, me.Player)

		action.Update(me, me.Enemy, me.Player, me.Shot, me.Explosion1, me.OccureBlock)
		action.HitCheck(me.Shot, me.Enemy)
		action.UniHitCheck(me.Shot, me.Block, me.OccureBlock)
		action.UniHitCheck(me.Player, me.Enemy)
		action.UniHitCheck(me.Vortex, me.Player)
		action.HitWall(me.Player, me.Block, me.OccureBlock)
		action.HitWall(me.Enemy, me.Block, me.OccureBlock)
		action.InScreen(me.Inner, me.Player)
		action.GoOutside(me.Outer, me.Player, me.Shot, me.Enemy)
		action.Clean(me.Player, me.Shot, me.Enemy, me.Explosion1, me.Vortex)

		if me.Vortex.Len() == 0 {
			me.EventTrigger(eventid.StageClear, nil, nil)
		}
	}

	me.DirectPushed.Update()
	if me.DirectPushed.Check(ginput.PressValue(ebiten.KeyF1)) {
		me.Debug = !me.Debug
	}
	if me.DirectPushed.Check(ginput.PressValue(ebiten.Key1)) {
		global.FlagOn()
	} else {
		global.FlagOff()
	}
}

func (me *Stage1) ObjectCount() (int) {
	res := 0
	res += me.Player.Len()
	res += me.Shot.Len()
	res += me.Block.Len()
	res += me.Enemy.Len()
	res += me.Explosion1.Len()
	return res
}

func (me *Stage1) Draw(screen *ebiten.Image) {
	screen.DrawImage(me.EnemyImage, me.Enemy.Options())
	screen.DrawImage(me.BlockImage, me.Block.Options())
	screen.DrawImage(me.VortexImage, me.Vortex.Options())
	screen.DrawImage(me.OccureBlockImage, me.OccureBlock.Options())
	screen.DrawImage(me.PlayerImage, me.Player.Options())
	screen.DrawImage(me.ShotImage, me.Shot.Options())
	screen.DrawImage(me.Explosion1Image, me.Explosion1.Options())

	me.Template.SetFloat("FrameCount", ebiten.CurrentFPS())
	me.Template.SetInt("ObjectCount", me.ObjectCount())
	if me.PlayerEntity != nil {
		me.Template.SetInt("PlayerEndurance", me.PlayerEntity.Endurance)
	} else {
		me.Template.SetInt("PlayerEndurance", 0)
	}

	me.Template.SetInt("Frame", me.Frame)
	me.Template.SetInt("TotalFrame", me.GameStatus.Frame + me.Frame)

	me.Instrument.UpdateText(me.Template.Text())
	screen.DrawImage(me.Instrument.Image(), me.Instrument.Options())

	if me.IsGameEnd() {
		me.Result.Update(screen)
	}

	if me.Debug {
		hitObjs := sprites.NewHitObjects(me.Block, me.OccureBlock, me.Shot, me.Enemy, global.RectDebug)
		screen.DrawImage(me.HitImage, hitObjs.Options())
	}
}

func (me *Stage1) EventTrigger(id event.Id, argument interface{}, origin orig.Interface) {
	switch id {
		case eventid.Block:
			me.Block.Occure(block.NewBlock(argument.(fig.FloatPoint)))

		case eventid.OccureBlock:
			me.OccureBlock.Occure(block.NewOccureBlock(argument.(block.Config)))

		case eventid.Player:
			me.PlayerEntity = player.New(argument.(fig.FloatPoint))
			me.Player.Occure(me.PlayerEntity)

		case eventid.Shot:
			o := argument.(orig.Interface)
			me.Shot.Occure(shot.New(o.Point(), o.Direction()))

		case eventid.Enemy:
			setting := argument.(funcs.EnemyConfig)
			me.Enemy.Occure(enemy.New(setting))

		case eventid.Explosion1:
			pt := argument.(fig.FloatPoint)
			me.Explosion1.Occure(explosion.NewExplosion1(pt))

		case eventid.Vortex:
			pt := argument.(fig.FloatPoint)
			me.Vortex.Occure(vortex.New(pt))

/*
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
			*/

		case eventid.StageClear:
			msg := ""
			if me.GameStatus.Stage == LastStage {
				msg = fmt.Sprintf(` All Clear
  Stage %d
  Frame: %d
  TotalFrame: %d
  End`, me.GameStatus.Stage, me.Frame, me.GameStatus.Frame + me.Frame)
			} else {
				msg = fmt.Sprintf(` Clear
  Stage %d
  Frame: %d
  TotalFrame: %d
  Next`, me.GameStatus.Stage, me.Frame, me.GameStatus.Frame + me.Frame)
			}

			me.GameEnd = true
			me.GameStatus.Stage += 1
			me.GameStatus.Frame += me.Frame
			global.SetGameStatus(me.GameStatus)
			me.Result = result.New(strings.Split(msg, "\n"))

	}
}

func (me *Stage1) Main(screen *ebiten.Image) (bool) {
	me.Update()
	me.Draw(screen)
	return !me.SceneEnd
}

func (me *Stage1) Dispose() {
}

func (me *Stage1) CreateReturnValue() (scene.Parameter) {
	if me.Restart || me.GameStatus.Stage > LastStage || me.GameStatus.Once {
		return []string{"title"}
	} else {
		return []string{"stage"}
	}
}

func (me *Stage1) ReturnValue() (scene.Parameter) {
	return me.CreateReturnValue()
}

func CreateStageScript(src MapData) ([]script.Proc) {
	res := []script.Proc{}
	for y, line := range src {
		for x, v := range line {
			x, y := float64(x * 32), float64(y * 32)
			if v == 1 {
				pt := fig.FloatPoint{x, y}
				res = append(res, script.NewEventProc(eventid.Block, pt))
			} else if v == 2 {
				config := block.Config {
					Point:           fig.FloatPoint{x, y},
					Span:            240,
					OccureDirection: block.OccureLeft,
				}
				res = append(res, script.NewEventProc(eventid.OccureBlock, config))
			} else if v == 3 {
				config := block.Config {
					Point:           fig.FloatPoint{x, y},
					Span:            240,
					OccureDirection: block.OccureRight,
				}
				res = append(res, script.NewEventProc(eventid.OccureBlock, config))
			} else if v == 4 {
				config := block.Config {
					Point:           fig.FloatPoint{x, y},
					Span:            240,
					OccureDirection: block.OccureRand,
				}
				res = append(res, script.NewEventProc(eventid.OccureBlock, config))
			} else if v == 8 {
				res = append(res, script.NewEventProc(eventid.Vortex, fig.FloatPoint{x, y}))
			} else if v == 9 {
				res = append(res, script.NewEventProc(eventid.Player, fig.FloatPoint{x, y}))
			}
		}
	}
	//res = append(res, script.NewEventProc(eventid.Player, fig.FloatPoint{100, 100}))

	//res = append(res, script.NewEventProc(eventid.Enemy, fig.FloatPoint{740, 0}))

	res = append(res, script.NewWaitProc(1))
	res = append(res, script.NewJumpProc(len(res)))

	return res
}

func GetMapSource(mapNo int64) (MapData) {
	switch mapNo {
		case 1: return Map1Source
		case 2: return Map2Source
		case 3: return Map3Source
	}
	return nil
}

func New(args scene.Parameter) (scene.Interface) {
	gameStatus := global.GetGameStatus()
	mapSource := MapData{}

	if len(args) > 0 {
		stageNo, _ := strconv.ParseInt(args[0], 10, 64)
		gameStatus.InitOnce(int(stageNo))
	} else if gameStatus.Stage <= 1 {
		gameStatus.Init()
	}

	if src := GetMapSource(int64(gameStatus.Stage)); src == nil {
		log.Exit("ステージマップ読み込みエラー: %d", gameStatus.Stage)
	} else {
		mapSource = src
	}

	scriptSource := CreateStageScript(mapSource)

	hitImage, _ := ebiten.NewImage(1, 1, ebiten.FilterLinear)
	hitImage.Fill(color.RGBA{0xff, 0x00, 0x00, 0x77})

	panelRect := fig.Rect {0, 0, 800, 32}

	conf := global.KeyConfig()
	if e1 := conf.Load(global.Path().KeyConfig()); e1 != nil {
		log.Log("keyconfig#New#conf.Load error %s", e1.Error())
		conf.Set(ginput.DefaultKeyBoard())
	}

	return &Stage1{
		GameStatus:       gameStatus,
		KeyConfig:        conf,

		Player:           player.NewObjects(),
		PlayerImage:      LoadImage("./resource/image/player.png"),

		Shot:             sprites.NewObjects(),
		ShotImage:        LoadImage("./resource/image/shot.png"),

		Enemy:            enemy.NewObjects(),
		EnemyImage:       LoadImage("./resource/image/enemy.png"),

		Explosion1:        sprites.NewObjects(),
		Explosion1Image:   LoadImage("./resource/image/explosion.png"),

		Vortex:           sprites.NewObjects(),
		VortexImage:      LoadImage("./resource/image/vortex.png"),

		Block:            sprites.NewObjects(),
		BlockImage:       LoadImage("./resource/image/wall.png"),

		OccureBlock:      sprites.NewObjects(),
		OccureBlockImage: LoadImage("./resource/image/enemy_occure.png"),

		HitImage:         hitImage,

		Template:         ttpl.New(PanelTemplate),
		Instrument:       instrument.NewInstrument(panelRect, 20, color.RGBA{0xb2, 0x9a, 0x8e, 0xff}),

		Pushed:           ginput.NewPushed(),
		DirectPushed:     ginput.NewDirectPushed(),

		Inner:            fig.Rect{0, -64, 800, 600},
		Outer:            fig.Rect{-64, -64, 800 + 64, 600 + 64},

		Source:           script.NewSource(scriptSource),
		Debug:            false,
	}
}
