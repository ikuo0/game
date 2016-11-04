
package event

import (
	"github.com/ikuo0/game/lib/orig"
)

type Id int

const (
	Wait Id = iota + 1
	Jump

	User Id = 100
)

func (me Id) String() (string) {
	switch me {
		case Wait: return "Wait"
		case Jump: return "Jump"
	}
	return "Unknown"
}

//########################################
//# Trigger
//########################################
type Trigger interface {
	EventTrigger(Id, interface{}, orig.Interface)
}

