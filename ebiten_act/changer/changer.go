
package changer

import (
	"github.com/ikuo0/game/ebiten_act/standby"
	"github.com/ikuo0/game/lib/log"
	"github.com/ikuo0/game/lib/scene"
	"github.com/hajimehoshi/ebiten"
	"fmt"
)

type SceneList map[string]scene.New

type Changer struct {
	Waiting scene.Interface
	Current scene.Interface
	List SceneList
}

func NewChanger(list SceneList) (*Changer) {
	waiting := standby.New([]string{})
	return &Changer {
		Waiting: waiting,
		Current: waiting,
		List: list,
	}
}

func (me *Changer) Apply(args []string) {
	me.Current = me.Waiting
	if len(args) < 1 {
		log.Exit("changer#Apply no argument")
	} else {
		fmt.Println(args)
		sceneName := args[0]
		opts := scene.Parameter(args[1:])

		if fn, ok := me.List[sceneName]; !ok {
			log.Exit("changer#Apply illegal scen name %s", sceneName)
		} else {
			go func() {
				me.Current = fn(opts)
			} ()
		}
	}
}

func (me *Changer) Play(screen *ebiten.Image) (error) {
	isContinue := me.Current.Main(screen)
	if !isContinue {
		me.Current.Dispose()
		me.Apply(me.Current.ReturnValue())
		return nil
	} else {
		return nil
	}
}

