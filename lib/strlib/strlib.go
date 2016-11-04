
package strlib

import (
	"strings"
)

type Queue struct {
	Data []string
}

func (me *Queue) Pop() (string) {
	if me.Empty() {
		return ""
	} else {
		res := me.Data[len(me.Data)-1] 
		me.Data = me.Data[:len(me.Data)-1]// 末尾削除
		return res
	}
}

func (me *Queue) Shift() (string) {
	if me.Empty() {
		return ""
	} else {
		res := me.Data[0]
		me.Data = me.Data[1:]
		return res
	}
}

func (me *Queue) Push(s string) {
	me.Data = append(me.Data, s)
}

func (me *Queue) Len() (int) {
	return len(me.Data)
}

func (me *Queue) Empty() (bool) {
	return len(me.Data) < 1
}

func (me *Queue) Ary() ([]string) {
	return me.Data
}

func (me *Queue) SetAry(ary []string) {
	me.Data = ary
}

func (me *Queue) SetSplit(s, sep string) {
	me.Data = strings.Split(s, sep)
}

func NewQueue(ary []string) (*Queue) {
	return &Queue{
		Data: ary,
	}
}
