
package log

import (
    "testing"
    "fmt"
    "time"
)

func Case1() {
    fmt.Println("start")
    defer fmt.Println("end")

    x := New("./log.txt")
    defer x.Dispose()

    x.Log("hogehoge %d", 12)

    time.Sleep(time.Second * 1)
}

func Test(t *testing.T) {
    Case1()
    //time.Sleep(time.Second * 1)
}
