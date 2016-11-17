
package radian

import (
	"testing"
	//"math"
	"fmt"
)

func CaseDeg360() {
	for i := 0; i < 360; i++ {
		deg := NormalizeDeg(i)
		rad := ToRad(deg)
		fmt.Printf("%d\t%g\n", deg, rad)
	}
}

func CaseAround360() {
	for i := 0; i < 1000; i++ {
		fmt.Printf("%d\t%d\n", i, NormalizeDeg(i))
	}
}

func CaseAround() {
	r := Radian(0)
	advance := 1
	d := 360 / advance
	m := 360 % advance
	max := d
	if m > 0 {
		max += 1
	}
	for i := 0; i < max; i++ {
		fmt.Printf("%f\t%f\n", r, r.Deg())
		r = r.TurnRight(advance)
	}
}

func CaseFormat() {
	x := 0.3490658503988659
	fmt.Printf("%%e: %e\n", x)
	fmt.Printf("%%E: %E\n", x)
	fmt.Printf("%%f: %f\n", x)
	fmt.Printf("%%g: %g\n", x)
	fmt.Printf("%%G: %G\n", x)
}

func CaseOther() {
	a := 585
	b := -5124
	fmt.Println(a % 360, b % 360)
}

func Test(t *testing.T) {
	//CaseFormat()
	//CaseDeg360()
	//CaseAround360()
	CaseOther()
}

