
package title

import (
	//"fmt"
	"os"
	"github.com/ikuo0/game/ebiten_act/global"
	"github.com/ikuo0/game/lib/fontmap"
	"github.com/ikuo0/game/lib/ginput"
	"github.com/ikuo0/game/lib/log"
	"github.com/ikuo0/game/lib/menu"
	"github.com/ikuo0/game/lib/scene"
	"github.com/hajimehoshi/ebiten"
	//"github.com/hajimehoshi/ebiten/ebitenutil"
)

type MenuId int

const (
	Start MenuId = 1 + iota
	KeyConfig
	Exit
)

type Item struct {
	Id   MenuId
	Text string
}

type Items struct {
	Data []Item
}

func (me *Items) Len() (int) {
	return len(me.Data)
}

func (me *Items) Value(i int) (int) {
	return int(me.Data[i].Id)
}

func (me *Items) Text(i int) (string) {
	return me.Data[i].Text
}

//########################################
//# Title
//########################################
type Title struct {
	Counter    int
	MenuFont   *fontmap.FontMap
	Menu       *menu.Menu
	Pushed     *ginput.Pushed
	SelectedId MenuId
	Soundset   *global.SoundsetSt
}

func New(args scene.Parameter) (scene.Interface) {
	ginput.Initialize()

	items := Items {
		Data: []Item {
			{Start, "Start"},
			{KeyConfig, "KeyConfig"},
			{Exit, "Exit"},
		},
	}


	if menuFont, e1 := fontmap.New(global.Path().Resource().File(`font/ipa/ipag.ttf`), 24); e1 != nil {
		log.Log(`title#New#fontmap.New: %s`, e1.Error())
		os.Exit(1)
		return nil
	} else if soundSet, e2 := global.Soundset(); e2 != nil {
		log.Log(`keyconfig#New#global.Soundset: %s`, e2.Error())
		os.Exit(1)
		return nil
	} else {
		return &Title {
			Counter:  0,
			MenuFont: menuFont,
			Menu:     menu.NewMenu(&items),
			Pushed:   ginput.NewPushed(),
			Soundset: soundSet,
		}
	}
}

func (me *Title) Main(screen *ebiten.Image) (bool) {
	me.Pushed.Update()
	if me.Pushed.Check(ginput.Up) {
		me.Soundset.MenuMove()
		me.Menu.Prev()
	} else if me.Pushed.Check(ginput.Down) {
		me.Soundset.MenuMove()
		me.Menu.Next()
	} else if me.Pushed.Check(ginput.Key1) {
		me.Soundset.MenuSubmit()
		me.SelectedId = MenuId(me.Menu.Value())
		return false
	}

	me.MenuFont.Draw(screen, 10, 10, me.Menu.String())
	me.Counter++
	return true
}

func (me *Title) Dispose() {
	me.Counter = 0
}

func (me *Title) ReturnValue() (scene.Parameter) {
	switch me.SelectedId {
		case Start:
			return []string{"stage1"}

		case KeyConfig:
			return []string{"keyconfig"}
	}
	return []string{"exit"}
}

