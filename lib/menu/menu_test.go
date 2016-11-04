
package menu

import (
	"testing"
	"fmt"
)

type Items struct {
	Data []string
}

func (me *Items) Len() (int) {
	return len(me.Data)
}

func (me *Items) Value(i int) (int) {
	return i
}

func (me *Items) Text(i int) (string) {
	return me.Data[i]
}

func NewItems(ary []string) (*Items) {
	return &Items {
		Data: ary,
	}
}

func Test1(t *testing.T) {
	items := NewItems([]string{"hoge", "piyo", "fuga", "exit"})
	x := NewMenu(items)

	op := []byte{'n','n','n','p','p','p', 'x', 'q'}

	for _, c := range op {
		switch c {
			case 'p':
				x.Prev()

			case 'n':
				x.Next()

			case 'x':
				val := x.Value()
				text := x.Text()
				fmt.Println("submit", val, text)

			case 'q':
				fmt.Println("exit")
				return

			default:
		}
		fmt.Println("-----------\n")
		fmt.Println(x.String())
		fmt.Println("\n")
	}
}
