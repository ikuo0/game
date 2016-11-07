
package sound

import (
	"github.com/ikuo0/game/lib/log"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/vorbis"
	"github.com/hajimehoshi/ebiten/audio/wav"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"fmt"
	"io/ioutil"
	"time"
)

var Context *audio.Context = nil
var SmapleRate int = 44100
var BgmVolume  float64 = 24
var SEVolume   float64 = 24

func NewContext(sampleRate int) () {
	if a, e1 := audio.NewContext(sampleRate); e1 != nil {
		log.Exit("sound#NewContext#NewContext error: %s", e1.Error())
	} else {
		Context = a
		SmapleRate = sampleRate
	}
}

func Initialize(sampleRate int, seVolume, bgmVolume float64) {
	NewContext(sampleRate)
	SEVolume  = seVolume / float64(128)
	BgmVolume = bgmVolume / float64(128)
}

func Update() {
	Context.Update()
}

//########################################
//# Wave
//########################################
type PlayerQueue []*audio.Player
func (me *PlayerQueue) Pop() (*audio.Player, bool) {
	if len(*me) < 1 {
		return nil, false
	} else {
		res, body := (*me)[0], (*me)[1:]
		*me = body
		return res, true
	}
}
func (me *PlayerQueue) Push(v *audio.Player) {
	*me = append(*me, v)
}
func (me PlayerQueue) Split() (PlayerQueue, PlayerQueue) {
	actives := PlayerQueue{}
	unactives := PlayerQueue{}
	for _, v := range me {
		if v.IsPlaying() {
			actives = append(actives, v)
		} else {
			unactives = append(unactives, v)
		}
	}
	return actives, unactives
}

type Wav struct {
	Data []byte
	SEPlayer *audio.Player
	Actives   PlayerQueue
	Unactives PlayerQueue
}

func (me *Wav) Load(fileName string) (error) {
	if f, e1 := ebitenutil.OpenFile(fileName ); e1 != nil {
		return fmt.Errorf("Wav#Load#OpenFile error: %s", e1)
	} else if s, e3 := wav.Decode(Context, f); e3 != nil {
		return fmt.Errorf("Wav#Load#Decode error: %s", e3)
	} else if b, e4 := ioutil.ReadAll(s); e4 != nil {
		return fmt.Errorf("Wav#Load#ReadAll error: %s", e4)
	} else {
		me.Data = b
		return nil
	}
}

//var activesLen, unactivesLen int
func (me *Wav) Play(pos time.Duration) {
	me.Actives, me.Unactives = me.Actives.Split()

/*
	if len(me.Actives) != activesLen || len(me.Unactives) != unactivesLen {
		a, u := len(me.Actives), len(me.Unactives)
		fmt.Println(a, u)
		activesLen, unactivesLen = a, u
	}
	*/

	if player, exist := me.Unactives.Pop(); exist {
		player.SetVolume(SEVolume)
		player.Seek(pos)
		player.Play()
		return
	} else if sePlayer, e1 := audio.NewPlayerFromBytes(Context, me.Data); e1 != nil {
		log.Log("Wav#Play#NewPlayerFromBytes error: %s", e1.Error())
		return
	} else {
		sePlayer.SetVolume(SEVolume)
		sePlayer.Seek(pos)
		sePlayer.Play()
		me.Actives.Push(sePlayer)
		return
	}
}

func (me *Wav) Dispose() (error) {
	var e error
	for _, v := range me.Actives {
		if e1 := v.Close(); e1 != nil {
			e = e1
		}
	}
	for _, v := range me.Unactives {
		if e1 := v.Close(); e1 != nil {
			e = e1
		}
	}
	return e
}

func NewWav(fileName string) (*Wav, error) {
	x := Wav{}
	if err := x.Load(fileName); err != nil {
		return nil, err
	} else {
		return &x, nil
	}
}


//########################################
//# Ogg
//########################################
type Ogg struct {
	Player *audio.Player
}

func (me *Ogg) Load(fileName string) (error) {
	if f, e1 := ebitenutil.OpenFile(fileName ); e1 != nil {
		return fmt.Errorf("Ogg#Load#OpenFile error: %s", e1)
	} else if s, e3 := vorbis.Decode(Context, f); e3 != nil {
		return fmt.Errorf("Ogg#Load#Decode error: %s", e3)
	} else if p, e4 := audio.NewPlayer(Context, s); e3 != nil {
		return fmt.Errorf("Ogg#Load#NewPlayer error: %s", e4)
	} else {
		me.Player  = p
		return nil
	}
}

func (me *Ogg) Play(pos time.Duration) {
	me.Player.Seek(pos)
	me.Player.SetVolume(BgmVolume)
	me.Player.Play()
}

func (me *Ogg) IsPlaying() (bool) {
	return me.Player.IsPlaying()
}

func (me *Ogg) Dispose() (error) {
	return me.Player.Close()
}

func NewOgg(fileName string) (*Ogg, error) {
	x := Ogg{}
	if err := x.Load(fileName); err != nil {
		return nil, err
	} else {
		return &x, nil
	}
}
