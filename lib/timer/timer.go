// 漢字
package timer

type Frame struct {
	Counter int
}

func NewFrame(cnt int) (*Frame) {
	return &Frame {
		Counter: cnt,
	}
}

func (me *Frame) Start(cnt int) {
	me.Counter = cnt
}

func (me *Frame) Up() (bool) {
	if me.Counter <= 0 {
		return true
	}
	me.Counter--
	return me.Counter <= 0
}
