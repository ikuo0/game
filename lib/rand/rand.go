// 漢字
package rand

import (
    "math/rand"
)

func Seed(n int64) {
    rand.Seed(n)
}

func Rand() (float64) {
    return rand.Float64()
}

func RandInt(min, max int) (int) {
    return rand.Intn(max - min + 1) + min
}
