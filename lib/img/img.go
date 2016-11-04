
package img

import (
    "bytes"
    "fmt"
    "image"
    //_ "image/jpeg"
    _ "image/png"
    "io/ioutil"
    "../fig"
    "../log"
    "github.com/hajimehoshi/ebiten"
)

type Img struct {
    Eimage *ebiten.Image
}

func Decode(b []byte) (*ebiten.Image, error) {
    if img, _, e1 := image.Decode(bytes.NewBuffer(b)); e1 != nil {
        return nil, e1
    } else if eimg, e2 := ebiten.NewImageFromImage(img, ebiten.FilterNearest); e2 != nil {
        return nil, e2
    } else {
        return eimg, nil
    }
}

func New(fileName string) (res *Img, e error) {
    if b, e1 := ioutil.ReadFile(fileName); e1 != nil {
        log.Log("画像読み込みエラー %s %s", fileName, e1.Error())
        e = e1
        return
    } else if eimg, e2 := Decode(b); e2 != nil {
        log.Log("画像デコードエラー %s %s", fileName, e2.Error())
        e = e2
        return
    } else {
        res = &Img {
            Eimage: eimg,
        }
        return
    }
}

func (me *Img) Image() (*ebiten.Image) {
    return me.Eimage
}

func (me *Img) AutoPositions(w, h int) (res []fig.Point, e error) {
    imageWidth, imageHeight := me.Eimage.Size()
    if (imageWidth % w) != 0 {
        e = fmt.Errorf("画像幅がパターン幅の倍数ではありません 画像幅: %d, パターン幅: %d", imageWidth, w)
        return
    } else if (imageHeight % h) != 0 {
        e = fmt.Errorf("画像高がパターン高の倍数ではありません 画像高: %d, パターン高: %d", imageHeight, h)
        return
    } else {
        rangeX := imageWidth / w
        rangeY := imageHeight / h
        for y := 0; y < rangeY; y++ {
            for x := 0; x < rangeX; x++ {
                res = append(res, fig.Point {
                    X: x,
                    Y: y,
                })
            }
        }
        return
    }
}
