
package sound

import (
	"fmt"
	"io/ioutil"
	"os"
    "github.com/hajimehoshi/ebiten/audio"
    "github.com/hajimehoshi/ebiten/audio/vorbis"
    "github.com/hajimehoshi/ebiten/audio/wav"
    "github.com/hajimehoshi/ebiten/ebitenutil"
)

var Context *audio.Context = nil
var SmapleRate int = 44100

func NewContext(sampleRate int) () {
	if a, e1 := audio.NewContext(sampleRate); e1 != nil {
		fmt.Println("sound#NewContext#NewContext", e1)
		os.Exit(1)
	} else {
		Context = a
		SmapleRate = sampleRate
	}
}

func Update() {
	Context.Update()
}

//########################################
//# Wave
//########################################
type Wav struct {
	Data []byte
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

func (me *Wav) Play() {
	if sePlayer, e1 := audio.NewPlayerFromBytes(Context, me.Data); e1 != nil {
		fmt.Println("NewPlayerFromBytes error", e1)
		os.Exit(1)
		return
	} else {
		sePlayer.Play()
		return
	}
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

func (me *Ogg) Play(v int) {
	/*
	if me.Player.IsPlaying() {
		me.Player.Seek(0)
	}
	*/

	me.Player.SetVolume(float64(v) / 128)
	me.Player.Play()
}

func (me *Ogg) Dispose() {
	me.Player.Close()
}

func NewOgg(fileName string) (*Ogg, error) {
	x := Ogg{}
	if err := x.Load(fileName); err != nil {
		return nil, err
	} else {
		return &x, nil
	}
}
