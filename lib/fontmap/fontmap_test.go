// 漢字
package fontmap

import (
    "testing"
    "bytes"
    "fmt"
    "sort"
    "golang.org/x/image/math/fixed"
)

func CaseBinPrint() {
    // E38182
    BinPrint([]byte{0xE3, 0x81, 0x82})
}

func CaseBinPrint2() {
    bary := [][]byte {
        []byte{0xE3, 0x81, 0x82},
        []byte{0xE3, 0x81, 0x82},
        []byte{0xE3, 0x81, 0x82},
    }
    for _, b := range bary {
        BinPrint(b)
    }
}

func CaseGetCode() {
    bary := [][]byte {
        []byte(string("あ")),
        []byte(string("い")),
        []byte(string("ん")),
    }

    targetWord := []byte(string("ｑ"))

    i := sort.Search(len(bary), func(i int) bool {
        return bytes.Compare(bary[i], targetWord) >= 0
    })

    fmt.Println(i)

    str := "ほげ"
    for _, v := range str {
        fmt.Println([]byte(string(v)))
    }
}

func Utf8Search(b []byte) {
    i := sort.Search(len(Utf8Ary), func(i int) bool {
        return bytes.Compare(Utf8Ary[i], b) >= 0
    })
    if i < len(Utf8Ary) && bytes.Compare(Utf8Ary[i], b) == 0 {
        fmt.Println("Detect", string(Utf8Ary[i]))
    } else {
        fmt.Println("NoMatch", string(b))
    }
}

func CaseSearch() {
    Utf8Search([]byte("ほ"))
    Utf8Search([]byte("げ"))
}

func CaseFixedP() {
    for i := 0; i < 10000; i+=1000 {
        fmt.Println(i, fixed.P(i, 0))
    }

    fmt.Println("-----------------------")

    for i := 0; i < 10; i++ {
        fmt.Println(i, fixed.P(0, i))
    }
}

func CaseStrRange() {
    str := "あいうえお"
    for _, v := range str {
        fmt.Printf("[%s]\n", string(v))
    }
}

func CaseHan1() {
    str := "ABCDE"
    for _, v := range str {
        fmt.Printf("[%s]\n", string(v))
    }

    Utf8Search([]byte("A"))
}

func Test(t *testing.T) {
    //CaseBinPrint()
    //CaseBinPrint2()
    //CaseGetCode()
    //CaseSearch()
    //CaseFixedP()
    //CaseStrRange()
    CaseHan1()
}
