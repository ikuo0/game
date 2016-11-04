
package script

import (
	"../event"
	"fmt"
	"os"
	"testing"
)

//########################################
//# Object
//########################################
type Object struct {
	Identity string
	MyStack    Stack
}

func (me *Object) Stack() (*Stack) {
	return &me.MyStack
}

func NewObject(identity string) (*Object) {
	return &Object {
		Identity: identity,
	}
}

//########################################
//# Source
//########################################
type Source struct {
	Procs []Proc
	Index int
}

func (me *Source) Read(idx int) (Proc) {
	if idx >= len(me.Procs) {
		os.Exit(9)
	}
	res := me.Procs[idx]
	me.Index++
	return res
}

func NewSource(procs []Proc) (*Source) {
	return &Source {
		Procs: procs,
	}
}

//########################################
//# Monitor
//########################################
type Monitor struct {
}

func (me *Monitor) ScriptTrigger(stack *Stack) {
	fmt.Println(stack)
}

//########################################
//# Mainloop
//########################################
type Mainloop struct {
	Monitor Monitor
	Source *Source
	Objs []Object
}

func (me *Mainloop) Main() {
	for i := 0; i < 10; i++ {
		fmt.Println("###", i)
		for j := 0; j < len(me.Objs); j++ {
			Exec(me.Source, me.Objs[j].Stack(), &me.Monitor)
		}
	}
}

func NewMainloop() (*Mainloop) {
	return &Mainloop {
		Monitor: Monitor{},
		Source:  NewSource([]Proc {
			NewWaitProc(2),
			NewEventProc(event.Player, 999),
			NewWaitProc(2),
			NewJumpProc(0),
		}),
		Objs: []Object {
			*NewObject("hoge"),
			*NewObject("piyo"),
		},
	}
}

func Test(t *testing.T) {
	x := NewMainloop()
	x.Main()
}

