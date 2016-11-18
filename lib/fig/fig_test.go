// 漢字
package fig


import (
	"testing"
	"fmt"
)

func BenchmarkFloatToInt(b *testing.B) {
	fmt.Println("start")
	var x float64
	for i := int(0); i < b.N; i++ {
		x = float64(i)
	}
	fmt.Println("end", x)
}

func BenchmarkIntToInt(b *testing.B) {
	fmt.Println("start")
	var x int
	for i := int(0); i < b.N; i++ {
		x = i
	}
	fmt.Println("end", x)
}

func BenchmarkFloatComp(b *testing.B) {
	fmt.Println("start")
	l := float64(1)
	m := float64(2)
	n := float64(0)
	for i := int(0); i < b.N; i++ {
		if l < m {
			n = l - m
		}
	}
	fmt.Println("end", n)
}

func BenchmarkIntComp(b *testing.B) {
	fmt.Println("start")
	l := int(1)
	m := int(2)
	n := int(0)
	for i := int(0); i < b.N; i++ {
		if l < m {
			n = l - m
		}
	}
	fmt.Println("end", n)
}


func Test(t *testing.T) {
	a := Rect {0, 0, 10, 10}
	fmt.Println(a)
	a.Relative(Point{5, 5})
	fmt.Println(a)
}

func TestHit1(t *testing.T) {
/*
{252 314 316 378}] [{304 336 336 368}]
[292, 286, 356, 350] [304, 336, 336, 368]
*/
	a := Rect{292, 286, 356, 350}
	b := Rect{304, 336, 336, 368}
	fmt.Println(a.Hit(&b))
	fmt.Println(b.Hit(&a))
}

func TestReference(t *testing.T) {
	a := []Rect{{292, 286, 356, 350}, {292, 286, 356, 350}, {292, 286, 356, 350}}
	fn := func (x *[]Rect) () {
		x = nil
	}
	fmt.Println(a)
	fn(&a)
	fmt.Println(a)
}

