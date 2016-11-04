// 漢字
package script

import (
	"github.com/ikuo0/game/lib/event"
	"github.com/ikuo0/game/lib/orig"
)


//########################################
//# Stack
//########################################
type Stack struct {
	Counter int
	Data    interface{}
}

func (me *Stack) PC() (int) {
	return me.Counter
}
func (me *Stack) Advance() {
	me.Counter++
}
func (me *Stack) Push(data interface{}) {
	me.Data = data
}
func (me *Stack) Pop() (interface{}) {
	return me.Data
}

//########################################
//# Proc
//########################################
type Proc interface {
	Exec(*Stack)
	Yield() (bool)
	EventId() (event.Id)
}

//########################################
//# WaitProc
//########################################
type WaitProc struct {
	Duration int
	CountUp bool
}

func (me *WaitProc) Exec(stack *Stack) {
	cnt, _ := stack.Pop().(int)
	if cnt >= me.Duration {
		me.CountUp = true
		cnt = 0
		stack.Advance()
	} else {
		cnt++
		me.CountUp = false
	}
	stack.Push(cnt)
}
func (me *WaitProc) Yield() (bool) {
	return !me.CountUp
}
func (me *WaitProc) EventId() (event.Id) {
	return event.Wait
}
func NewWaitProc(duration int) (*WaitProc) {
	return &WaitProc {
		Duration: duration,
	}
}

//########################################
//# EventProc
//########################################
type EventProc struct {
	ResultEventId event.Id
	Data    interface{}
}

func (me *EventProc) Exec(stack *Stack) {
	stack.Push(me.Data)
	stack.Advance()
}
func (me *EventProc) Yield() (bool) {
	return false
}

func (me *EventProc) EventId() (event.Id) {
	return me.ResultEventId
}

func NewEventProc(id event.Id, arguments interface{}) (*EventProc) {
	return &EventProc {
		ResultEventId: id,
		Data:    arguments,
	}
}

//########################################
//# JumpProc
//########################################
type JumpProc struct {
	JumpCounter int
}

func (me *JumpProc) Exec(stack *Stack) {
	stack.Counter = me.JumpCounter
}
func (me *JumpProc) Yield() (bool) {
	return true
}
func (me *JumpProc) EventId() (event.Id) {
	return event.Jump
}
func NewJumpProc(pc int) (*JumpProc) {
	return &JumpProc {
		JumpCounter: pc,
	}
}

//########################################
//# Input
//########################################
type Input interface {
	Read(int) (Proc)
}

func Exec(program Input, stack *Stack, origin orig.Interface, output event.Trigger) {
	for {
		proc := program.Read(stack.PC())
		proc.Exec(stack)
		output.EventTrigger(proc.EventId(), stack.Pop(), origin)
		if proc.Yield() {
			break
		}
	}
}




//########################################
//# Source
//########################################
type Source struct {
	Procs []Proc
}

func (me *Source) Read(idx int) (Proc) {
	if idx >= len(me.Procs) {
		idx = 0
	}
	res := me.Procs[idx]
	return res
}

func NewSource(procs []Proc) (*Source) {
	return &Source {
		Procs: procs,
	}
}
