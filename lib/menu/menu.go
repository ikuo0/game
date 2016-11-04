
package menu

import (
	"strings"
)

type ItemInterface interface {
	Len() (int)
	Value(int) (int)
	Text(int) (string)
}

type Menu struct {
	Index int
	Items ItemInterface
}

func NewMenu(i ItemInterface) (*Menu) {
	return &Menu {
		Items: i,
	}
}

func (me *Menu) Prev() {
	me.Index--
	if me.Index < 0 {
		me.Index = me.Items.Len() - 1
	}
}

func (me *Menu) Next() {
	me.Index++
	if me.Index >= me.Items.Len() {
		me.Index = 0
	}
}

func (me *Menu) Value() (int) {
	return me.Items.Value(me.Index)
}

func (me *Menu) Text() (string) {
	return me.Items.Text(me.Index)
}

func (me *Menu) String() (string) {
	ary := []string{}
	for i := 0; i < me.Items.Len(); i++ {
		if i == me.Index {
			ary = append(ary, ">> " + me.Items.Text(i))
		} else {
			ary = append(ary, me.Items.Text(i))
		}
	}

	return strings.Join(ary, "\r\n")
}
