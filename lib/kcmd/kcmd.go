
package kcmd

import (
	"github.com/ikuo0/game/lib/ginput"
	//"fmt"
)

const Length = 256

//type KeyBuffer [Length]ginput.InputBits
type Buffer struct {
	Data [Length]ginput.InputBits
	Index int
}

func (me *Buffer) Update(bits ginput.InputBits) {
	me.Index++
	if me.Index >= Length {
		me.Index = 0
	}
	me.Data[me.Index] = bits
}

func (me *Buffer) Find(bits ginput.InputBits, startIndex, findFrame int) (int) {
	for i := 0; i < findFrame; i++ {
		if bits.And(ginput.Neutral) {
			if !me.Data[startIndex].And( ^ginput.Neutral & bits) {
				return startIndex
			}
		} else {
			if me.Data[startIndex].And(bits) {
				return startIndex
			}
		}
		if startIndex--; startIndex < 0 {
			startIndex = Length - 1
		}
	}
	return -1
}

func Check(command []ginput.InputBits, buf *Buffer, findFrame int) (bool) {
	idx := buf.Index
	for i := len(command) - 1; i >= 0; i-- {
		if idx = buf.Find(command[i], idx, findFrame); idx >= 0 {
			if idx--; idx < 0 {
				idx = Length - 1
			}
			continue
		} else {
			return false
		}
	}
	return true
}

/*
type Kcmd struct {
	Buffer KeyBuffer
	Index int
}

func (me *Kcmd) Update(bits ginput.InputBits) {
	me.Index++
	if me.Index >= Length {
		me.Index = 0
	}
	me.Buffer[me.Index] = bits
}

func (me *Kcmd) Check(command []ginput.InputBits, findFrame int) (bool) {
	idx := me.Index
	for i := len(command) - 1; i >= 0; i-- {
		if idx = me.Buffer.Find(command[i], idx, findFrame); idx >= 0 {
			if idx--; idx < 0 {
				idx = Length - 1
			}
			continue
		} else {
			return false
		}
	}
	return true
}

func New() (*Kcmd) {
	return &Kcmd {
	}
}
*/
