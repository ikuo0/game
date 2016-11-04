
package keyconfig

import (
	"github.com/ikuo0/game/ebiten_stg/global"
	"github.com/ikuo0/game/lib/fontmap"
	"github.com/ikuo0/game/lib/ginput"
	"github.com/ikuo0/game/lib/log"
	"github.com/ikuo0/game/lib/menu"
	"github.com/ikuo0/game/lib/scene"
	"github.com/hajimehoshi/ebiten"
	//"fmt"
	"os"
)

type MenuId int

const (
	AnyKey MenuId = 1 + iota
	Up
	Down
	Left
	Right
	Shot
	Shield
	Submit
	Cancel
)

func (me MenuId)  Bits() (ginput.InputBits) {
	switch me {
		case Up:     return ginput.Up
		case Down:   return ginput.Down
		case Left:   return ginput.Left
		case Right:  return ginput.Right
		case Shot:   return ginput.Key1
		case Shield: return ginput.Key2
	}
	return 0
}

type Item struct {
	Id     MenuId
	Text   string
}

type Items struct {
	Data []Item
}

func (me *Items) Len() (int) {
	return len(me.Data)
}

func (me *Items) Value(i int) (int) {
	if i >= len(me.Data) {
		return 0
	} else {
		return int(me.Data[i].Id)
	}
}

func (me *Items) Text(i int) (string) {
	if i >= len(me.Data) {
		return "unknown"
	} else {
		return me.Data[i].Text
	}
}

func (me *Items) CreateItem(menuId MenuId, bits ginput.InputBits, prefix string, conf ginput.Keymaps) (Item) {
	val := conf.ValueFromBits(bits)
	return Item {
		Id:   menuId,
		Text: prefix + "(" + val.String() + ")",
	}
}

func (me *Items) Update(maps ginput.Keymaps) {
	items := []Item{}

	items = append(items, me.CreateItem(Up,     ginput.Up, "Up", maps))
	items = append(items, me.CreateItem(Down,   ginput.Down, "Down", maps))
	items = append(items, me.CreateItem(Left,   ginput.Left, "Left", maps))
	items = append(items, me.CreateItem(Right,  ginput.Right, "Right", maps))
	items = append(items, me.CreateItem(Shot,   ginput.Key1, "Shot", maps))
	items = append(items, me.CreateItem(Shield, ginput.Key2, "Shield", maps))

	items = append(items, Item {
		Id:   Submit,
		Text: "Submit",
	})
	items = append(items, Item {
		Id:   Cancel,
		Text: "Cancel",
	})

	me.Data = items
}

//########################################
//# Chose
//########################################
type Chose struct {
	Values []ginput.PressValue
	First  ginput.PressValue
	Prev   ginput.PressValue
	Rec    []ginput.PressValue
}

func (me *Chose) IsFill(value ginput.PressValue, values []ginput.PressValue) (bool) {
	for _, v := range values {
		if value != v {
			return false
		}
	}
	return true
}

func (me *Chose) Check(values []ginput.PressValue) (ginput.PressValue) {
	if len(values) == 1 {
		if len(me.Rec) > 4 {
			me.Rec = me.Rec[1:]
		}
		value := values[0]
		me.Rec = append(me.Rec, value)
		if len(me.Rec) >= 4 && me.Rec[0] != me.Rec[1] && me.IsFill(me.Rec[1], me.Rec[1:3]) {
			return value
		} else {
			return 0
		}
	} else {
		me.Rec = []ginput.PressValue{0}
		return 0
	}
}

//########################################
//# KeyConfig
//########################################
type KeyConfig struct {
	MenuFont   *fontmap.FontMap
	Items      *Items
	Menu       *menu.Menu
	Pushed     *ginput.Pushed

	Chose      Chose
	Chosing    bool

	SelectedId MenuId

	Config     *global.KeyConfigSt

	Soundset   *global.SoundsetSt
}

func New(args scene.Parameter) (scene.Interface) {
	ginput.Initialize()

	conf := global.KeyConfig()
	if e1 := conf.Load(global.Path().KeyConfig()); e1 != nil {
		log.Log("keyconfig#New#conf.Load error %s", e1.Error())
		conf.Set(ginput.DefaultKeyBoard())
	}

	items := Items{}
	items.Update(conf.Maps)

	if menuFont, e1 := fontmap.New(global.Path().Resource().File(`font/ipa/ipag.ttf`), 24); e1 != nil {
		log.Log(`keyconfig#New#fontmap.New: %s`, e1.Error())
		os.Exit(1)
		return nil
	} else if soundSet, e2 := global.Soundset(); e2 != nil {
		log.Log(`keyconfig#New#global.Soundset: %s`, e2.Error())
		os.Exit(1)
		return nil
	} else {
		return &KeyConfig {
			MenuFont: menuFont,
			Items:    &items,
			Menu:     menu.NewMenu(&items),
			Pushed:   ginput.NewPushed(),
			Chosing:  false,
			Config:   conf,
			Soundset: soundSet,
		}
	}
}

func (me *KeyConfig) Main(screen *ebiten.Image) (bool) {
	defer func() {
		me.MenuFont.Draw(screen, 10, 10, me.Menu.String())
	} ()

	if me.Chosing {
		values := ginput.Values()
		if value := me.Chose.Check(values); value != 0 {
			me.Soundset.MenuSubmit()
			bits := MenuId(me.Menu.Value()).Bits()
			me.Config.Maps.Config(value, bits)
			me.Items.Update(me.Config.Maps)
			me.Chosing = false
			me.Pushed.Update()
		}
		return true
	} else {
		me.Pushed.Update()
		if me.Pushed.Check(ginput.Up) {
			me.Soundset.MenuMove()
			me.Menu.Prev()
		} else if me.Pushed.Check(ginput.Down) {
			me.Soundset.MenuMove()
			me.Menu.Next()
		} else if me.Pushed.Check(ginput.Key1) {
			menuId := MenuId(me.Menu.Value())
			switch menuId {
				case Up, Down, Left, Right, Shot, Shield:
					me.Chosing = true
					me.Chose = Chose{}
					me.Soundset.MenuSubmit()
					return true

				case Submit, Cancel:
					me.SelectedId = menuId
					if menuId == Submit {
						me.Soundset.MenuSubmit()
					} else {
						me.Soundset.MenuCancel()
					}
					return false
			}
		}
	}

	return true
}

func (me *KeyConfig) Dispose() {
	if me.SelectedId == Submit {
		me.Config.Save(global.Path().KeyConfig())
	}
}

func (me *KeyConfig) ReturnValue() (scene.Parameter) {
	return []string{"title"}
}

