
package global

import (
	"github.com/ikuo0/game/lib/ginput"
	"github.com/ikuo0/game/lib/log"
	"github.com/ikuo0/game/lib/sound"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

var rootPath string

func SetRootPath(root string) {
	rootPath = root
}

type PathSt struct {
	Root string
	Ary []string
}

func (me *PathSt) Join() (string) {
	return filepath.Join(me.Ary...)
}

func (me *PathSt) Resource() (*PathSt) {
	me.Ary = append(me.Ary, "resource")
	return me
}

func (me *PathSt) Log() (*PathSt) {
	me.Ary = append(me.Ary, "log")
	return me
}

func (me *PathSt) KeyConfig() (string) {
	return me.File("keyconfig.json")
}

func (me *PathSt) File(relative string) (string) {
	me.Ary = append(me.Ary, relative)
	return me.Join()
}

func NewPath(root string) (*PathSt) {
	return &PathSt {
		Root: root,
	}
}

func Path() (*PathSt) {
	return NewPath(rootPath)
}

//########################################
//# Key Config
//########################################
type KeyConfigSt struct {
	Maps ginput.Keymaps
}

func (me *KeyConfigSt) Load(fileName string) (error) {
	if b, e1 := ioutil.ReadFile(fileName); e1 != nil {
		log.Log("global#KeyConfigSt#Load#filelib.Get error: %s", e1.Error())
		return e1
	} else if e2 := json.Unmarshal(b, me); e2 != nil {
		log.Log("global#KeyConfigSt#Load#json.Unmarshal error: %s", e2.Error())
		return e2
	} else {
		return nil
	}
}

func (me *KeyConfigSt) Save(fileName string) (error) {
	if b, e1 := json.Marshal(me); e1 != nil {
		log.Log("global#KeyConfigSt#Save#json.Marshal error: %s", e1.Error())
		return e1
	} else if e2 := ioutil.WriteFile(fileName, b, os.ModePerm); e2 != nil {
		log.Log("global#KeyConfigSt#Save#filelib.Put error: %s", e2.Error())
		return e2
	} else {
		return nil
	}
}

func (me *KeyConfigSt) Set(maps ginput.Keymaps) {
	me.Maps = maps
}

func KeyConfig() (*KeyConfigSt) {
	return &KeyConfigSt{}
}

//########################################
//# Soundset
//########################################
const SampleRate = 44100

type SoundsetSt struct {
	SubmitWav *sound.Wav
	CancelWav *sound.Wav
	MoveWav   *sound.Wav
}
/*
	SubmitWav *sound.Wav
	CancelWav *sound.Wav
	MoveWav   *sound.Wav
*/

func (me *SoundsetSt) MenuSubmit() {
	me.SubmitWav.Play()
}

func (me *SoundsetSt) MenuCancel() {
	me.CancelWav.Play()
}

func (me *SoundsetSt) MenuMove() {
	me.MoveWav.Play()
}

func Soundset() (*SoundsetSt, error) {
	if submitWav, e1 := sound.NewWav(`./resource/sound/menu_submit.wav`); e1 != nil {
		log.Log("global#Soundset#sound.NewWav error: %s", e1.Error())
		return nil, e1
	} else if cancelWav, e2 := sound.NewWav(`./resource/sound/menu_cancel.wav`); e2 != nil {
		log.Log("global#Soundset#sound.NewWav error: %s", e2.Error())
		return nil, e2
	} else if moveWav, e3 := sound.NewWav(`./resource/sound/menu_move.wav`); e3 != nil {
		log.Log("global#Soundset#sound.NewWav error: %s", e3.Error())
		return nil, e3
	} else {
		return &SoundsetSt {
			SubmitWav: submitWav,
			CancelWav: cancelWav,
			MoveWav:   moveWav,
		}, nil
	}
}






