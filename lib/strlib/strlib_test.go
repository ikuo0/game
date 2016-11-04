
package strlib

import (
	"testing"
	"fmt"
)

func Test1(t *testing.T) {
	ary1 := []string{"hoge", "piyo", "fuga"}
	a := NewQueue(ary1)
	b := NewQueue(nil)
	for !a.Empty() {
		b.Push(a.Shift())
		fmt.Println(b.Ary())
	}
}
