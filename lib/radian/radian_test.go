
package radian

import (
	"testing"
	"math"
	"fmt"
)

func CaseDeg360() {
	for i := 0; i < 360; i++ {
		fmt.Println(float64(i) * math.Pi / 180)
	}
}

func Test(t *testing.T) {
	CaseDeg360()
}
