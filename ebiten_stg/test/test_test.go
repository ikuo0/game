
package test

import (
	"fmt"
	"testing"
)

func Benchmark1(b *testing.B) {
	var x interface{}
	x = int(1)
	res := int(0)
	for i := 0; i < 1000000; i++ {
		res += x.(int) + i
	}
	fmt.Println(res)
}

func Benchmark2(b *testing.B) {
	var x float64
	x = 1
	res := int(0)
	for i := 0; i < 1000000; i++ {
		res += int(x) + i
	}
	fmt.Println(res)
}

func BenchmarkA(b *testing.B) {
	var x []float64
	res := float64(0)
	for i := 0; i < 1000000; i++ {
		x = []float64{1,2,3,4}
		a := x
		res += a[0]
	}
	fmt.Println(res)
}

func BenchmarkB(b *testing.B) {
	var x interface{}
	res := float64(0)
	for i := 0; i < 1000000; i++ {
		x = []float64{1,2,3,4}
		a := x.([]float64)
		res += a[0]
	}
	fmt.Println(res)
}

