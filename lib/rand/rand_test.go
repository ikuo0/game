// 漢字
package rand

import (
    "testing"
    "fmt"
    "math/rand"
    "time"
)

func CaseInt() {
    rand.Seed(int64(time.Now().Nanosecond()))
    for i := 0; i < 10; i++ {
        fmt.Println(RandInt(0,10))
    }
}

func CaseFloat() {
    rand.Seed(int64(time.Now().Nanosecond()))
    for i := 0; i < 10; i++ {
        fmt.Println(Rand())
    }
}

func Test(t *testing.T) {
    //CaseInt()
    CaseFloat()
}
