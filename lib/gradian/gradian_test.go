
package gradian

import (
	"github.com/ikuo0/game/lib/radian"
	"fmt"
	"testing"
)

func CaseNormalize() {
	for i := 0; i < 360; i++ {
		fmt.Printf("%d\t%d\n", i, NormalizeDeg(i))
	}
}

func CaseDegree() {
	a := Degree{}
	for i := 0; i < 360; i++ {
		fmt.Println(a.Deg, a.Radian())
		a.RightTurn(1)
	}
}

func CaseOther() {
	fmt.Println(270, radian.DegArray[270])
	fmt.Println(90, radian.DegArray[ToIndex(90)])
	fmt.Println(45, radian.DegArray[45])
	fmt.Println(315, radian.DegArray[ToIndex(315)])
}

func Test(t *testing.T) {
	//CaseNormalize()
	CaseDegree()
	//CaseOther()
}
