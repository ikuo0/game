
package log

import (
    "fmt"
    "os"
    "time"
)

type LogSt struct {
    Text chan string
    Stop chan bool
    LogFile *os.File
}

func Logging(l *LogSt) {
    for {
loop:
        select {
            case text := <- l.Text:
                fmt.Println(text)
                l.LogFile.WriteString(text + "\n")

            case <- l.Stop:
                break loop
        }
    }
}

func New(fileName string) (res *LogSt) {
    if f, e1 := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666); e1 != nil {
        panic("cannnot open :" + fileName + ": " + e1.Error())
        return
    } else {
        res = &LogSt {
            Text:    make(chan string),
            Stop:    make(chan bool),
            LogFile: f,
        }
        go Logging(res)
        return
    }
}

func (me *LogSt) Dispose() {
    me.Stop <- true
    me.LogFile.Close()
}

func (me *LogSt) Log(f string, args ...interface{}) {
    me.Text <- fmt.Sprintf(time.Now().String() + " " + f, args...)
}

var logBody *LogSt

func Start(fileName string) {
    logBody = New(fileName)
}

func Log(f string, args ...interface{}) {
    if logBody == nil {
        return
    } else {
        logBody.Log(f, args...)
    }
}

func Exit(f string, args ...interface{}) {
    Log(f, args...)
    End()
    time.Sleep(time.Millisecond * 50)
    os.Exit(1)
}

func End() {
    if logBody == nil {
        return
    } else {
        logBody.Dispose()
    }
}

