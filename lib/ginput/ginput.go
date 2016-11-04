
package ginput

import (
	"github.com/hajimehoshi/ebiten"
	"fmt"
)

type Kind int

const (
	Keyboard Kind = 1 + iota
	Joypad
)

type PressValue int
type PressValues []PressValue
func (me PressValues) Include(v PressValue) (bool) {
	for _, val := range me {
		if val == v {
			return true
		}
	}
	return false
}

func GetKeyInput() (res []PressValue) {
	for i := PressValue(ebiten.Key0); i <= PressValue(ebiten.KeyUp); i++ {
		if ebiten.IsKeyPressed(ebiten.Key(i)) {
			res = append(res, i)
		}
	}
	return
}

func PadCount() (res int) {
	for i := 0; i < 16; i++ {
		if cnt := ebiten.GamepadButtonNum(i); cnt > 0 {
			res++
		}
	}
	return
}

const (
	ButtonIndex PressValue = 0
	AxisIndex   PressValue = 500
	PadRadix    PressValue = 1000
)

func (me PressValue) String() (string) {
	if me < PadRadix {
		return KeyToString(ebiten.Key(me))
	} else {
		d := me / PadRadix
		m := me % PadRadix
		if m < AxisIndex {
			return fmt.Sprintf("Pad%d Button%d", d, ButtonIndex + m)
		} else {
			return fmt.Sprintf("Pad%d Axis%d", d, m - AxisIndex)
		}
	}
}

const AxisThreshold = 0.1

func GetPadInput(padNum int) (res []PressValue) {
	baseIndex := PressValue((padNum + 1)) * PadRadix
	buttonNum := ebiten.GamepadButtonNum(padNum)
	for i := 0; i < buttonNum; i++ {
		if ebiten.IsGamepadButtonPressed(padNum, ebiten.GamepadButton(i)) {
			res = append(res, baseIndex + PressValue(i) + ButtonIndex)
		}
	}

	axisNum := ebiten.GamepadAxisNum(padNum)
	for i := 0; i < axisNum && i < 4; i++ {// Axisは2個まで
		power := ebiten.GamepadAxis(padNum, i)
		if power != 0 {
			if power < -AxisThreshold {
				res = append(res, baseIndex + PressValue(i) * 2 + AxisIndex) // 偶数
			} else if power > AxisThreshold {
				res = append(res, baseIndex + PressValue(i) * 2 + AxisIndex + 1)// 奇数
			}
		}
	}

	return
}

type InputBits int64

func (me InputBits) And(bits InputBits) (bool) {
	return (me & bits) == bits
}

func (me InputBits) Or(bits InputBits) (bool) {
	return (me & bits) != 0
}

func (me InputBits) NotAnd(bits InputBits) (bool) {
	return !me.And(bits)
}

//########################################
//# Key map
//########################################
const (
	Left     InputBits = 0x00000001
	Up       InputBits = 0x00000002
	Right    InputBits = 0x00000004
	Down     InputBits = 0x00000008
	Key1     InputBits = 0x00000010
	Key2     InputBits = 0x00000020
	Key3     InputBits = 0x00000040
	Key4     InputBits = 0x00000080
	Key5     InputBits = 0x00000100
	Key6     InputBits = 0x00000200
	Key7     InputBits = 0x00000400
	Key8     InputBits = 0x00000800
	Key9     InputBits = 0x00001000
	Key10    InputBits = 0x00010000
	Key11    InputBits = 0x00020000
	Key12    InputBits = 0x00040000
	Key13    InputBits = 0x00080000
	Key14    InputBits = 0x00200000
	FullBits InputBits = 0x00ffffff
	KeyMask  InputBits = 0x0000fff0
	AxisMask InputBits = 0x0000000f

	Naxis    InputBits = 0x1000000f
	Nkey1    InputBits = 0x10000010
	Nkey2    InputBits = 0x10000020
	Nkey3    InputBits = 0x10000040
	Nkey4    InputBits = 0x10000080
	Nkey5    InputBits = 0x10000100
	Nkey6    InputBits = 0x10000200
	Nkey7    InputBits = 0x10000400
	Nkey8    InputBits = 0x10000800
	Nkey9    InputBits = 0x10001000
	Neutral  InputBits = 0x10000000
)

type Keymap struct {
	Value PressValue
	Bits  InputBits
}

type Keymaps []Keymap

func (me Keymaps) BitsFromValue(val PressValue) (InputBits) {
	for _, v := range me {
		if v.Value == val {
			return v.Bits
		}
	}
	return 0
}

func (me Keymaps) ValueFromBits(bits InputBits) (PressValue) {
	for _, v := range me {
		if v.Bits == bits {
			return v.Value
		}
	}
	return 0
}

func (me Keymaps) Config(value PressValue, bits InputBits) (Keymaps) {
	for i, v := range me {
		if v.Bits == bits {
			me[i].Value = value
			return me
		}
	}

	return append(me, Keymap {
		Value: value,
		Bits:  bits,
	})
}

func DefaultKeyBoard() (Keymaps) {
	return []Keymap {
		{60, Left},
		{64, Right},
		{68, Up},
		{42, Down},
		{44, Key1},
		{66, Key2},
		{35, Key3},
		{33, Key4},
		{12, Key5},
	}
}

func DefaultJoypad(padNum int) (Keymaps) {
	baseIndex := PressValue((padNum + 1)) * PadRadix
	return []Keymap {
		{baseIndex + AxisIndex + 0,   Left},
		{baseIndex + AxisIndex + 1,   Right},
		{baseIndex + AxisIndex + 2,   Up},
		{baseIndex + AxisIndex + 3,   Down},

		{baseIndex + 11,   Right}, // RAP
		{baseIndex + 13,   Left},
		{baseIndex + 10,   Up},
		{baseIndex + 12,   Down},

		{baseIndex + ButtonIndex + 0, Key1},
		{baseIndex + ButtonIndex + 1, Key2},
		{baseIndex + ButtonIndex + 2, Key3},
		{baseIndex + ButtonIndex + 3, Key4},
		{baseIndex + ButtonIndex + 4, Key5},
	}
}


//########################################
//# DirectPushed
//########################################
type DirectPushed struct {
	CurrentInput PressValues
	PreInput     PressValues
}

func (me *DirectPushed) Update() {
	values := Values()
	me.PreInput = me.CurrentInput
	me.CurrentInput = values
}

func (me *DirectPushed) Check(val PressValue) (bool) {
	return me.CurrentInput.Include(val) && !me.PreInput.Include(val)
}

func NewDirectPushed() (*DirectPushed) {
	return &DirectPushed {
	}
}

//########################################
//# Pushed
//########################################
type Pushed struct {
	CurrentInput InputBits
	PreInput     InputBits
	MergeInput   InputBits
}

func (me *Pushed) Update() {
	bits := Standard()
	me.PreInput = me.CurrentInput
	me.CurrentInput = bits
	me.MergeInput = me.CurrentInput & ^me.PreInput
}

func (me *Pushed) Check(bits InputBits) (bool) {
	return me.MergeInput.And(bits)
}

func NewPushed() (*Pushed) {
	return &Pushed {
	}
}

//########################################
//# global funcs
//########################################
var padCount int
func RefreshPad() {
	padCount = PadCount()
}

var DefaultConfig []Keymap
func SetupDefauiltConfig() {
	maps := DefaultKeyBoard()

	max := PadCount()
	for i := 0; i < max; i++ {
		maps = append(maps, DefaultJoypad(i)...)
	}

	DefaultConfig = maps
}

var PushedInput *DirectPushed

func Initialize() {
	RefreshPad()
	SetupDefauiltConfig()
	PushedInput = NewDirectPushed()
}

func Values() ([]PressValue) {
	values := []PressValue{}

	values = append(values, GetKeyInput()...)

	for i := 0; i < padCount; i++ {
		values = append(values, GetPadInput(i)...)
	}

	return values
}

func Bits(values []PressValue, conf Keymaps) (InputBits) {
	res := InputBits(0)
	for _, v := range values {
		res |= conf.BitsFromValue(v)
	}

	return res
}

func Standard() (InputBits) {
	values := Values()
	return Bits(values, DefaultConfig)
}


