
package anime

type Frame struct {
	Duration int
	Counter int
}

func (me *Frame) Tick() (bool) {
	if me.Counter >= me.Duration {
		return true
	} else {
		me.Counter++
		return false
	}
}
func (me *Frame) Reset() {
	me.Counter = 0
}

func NewFrame(d int) (*Frame) {
	return &Frame {
		Duration: d,
	}
}

type Frames struct {
	Frames []Frame
	Count  int
	Around bool
}
func (me *Frames) Update() {
	if me.Frames[me.Count].Tick() {
		me.Frames[me.Count].Reset()
		me.Count++
		if me.Count >= len(me.Frames) {
			me.Around = true
			me.Count = 0
		} else {
			me.Around = false
		}
	}
}

func (me *Frames) Index() (int) {
	return me.Count
}

func (me *Frames) Arounded() (bool) {
	return me.Around
}

func NewFrames(frames ...int) (*Frames) {
	ary := []Frame{}
	for _, v := range frames {
		ary = append(ary, *NewFrame(v))
	}
	return &Frames {
		Frames: ary,
	}
}
