// 漢字
package fontmap

import (
	"github.com/ikuo0/game/lib/fig"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	//"os"
	"sort"
	"strings"
)

const RowCount = 60
const ColCount = 60

type Info struct {
	Bin  []byte
	SrcX int
	SrcY int
	Han  bool
}

type FontMap struct {
	SrcWidth int
	SrcHeight int
	Canvas *ebiten.Image
	Info   []Info
}

func LoadFont(fileName string) (*truetype.Font, error) {
	if f, e1 := ebitenutil.OpenFile(fileName); e1 != nil {
		return nil, fmt.Errorf("フォントファイルオープンエラー %s %s", fileName, e1.Error())
	} else {
		defer f.Close()
		if b, e2 := ioutil.ReadAll(f); e2 != nil {
			return nil, fmt.Errorf("フォントファイル読み込みエラー %s %s", fileName, e2.Error())
		} else if tt, e3 := truetype.Parse(b); e3 != nil {
			return nil, fmt.Errorf("フォントパースエラー %s %s", fileName, e3.Error())
		} else {
			return tt, nil
		}
	}
}

func GetTTFace(size int, tt *truetype.Font) (font.Face) {
	const dpi = 72
	return truetype.NewFace(tt, &truetype.Options {
		Size:	float64(size),
		DPI:	 dpi,
		Hinting: font.HintingFull,
	})
}

func GetDrawer(dst *image.RGBA, face font.Face) (*font.Drawer) {
	return &font.Drawer{
		Dst:  dst,
		Src:  image.White,
		Face: face,
	}
}

//3489
// 60 * 60
func CreateCanvas(fontWidth, fontTotalHeight int) (*ebiten.Image, error) {
	if canvas, e1 := ebiten.NewImage(ColCount * fontWidth, RowCount * fontTotalHeight, ebiten.FilterNearest); e1 != nil {
		return nil, fmt.Errorf("フォント描画領域作成エラー %s", e1.Error())
	} else {
		return canvas, nil
	}
}

func Draw(srcWidth, srcHeight, baseLine int, canvas *ebiten.Image, dst *image.RGBA, d *font.Drawer) (res []Info, e error) {
	for i, bary := range Utf8Ary {
		posX := i % ColCount
		posY := i / RowCount
		srcX := posX * srcWidth
		srcY := posY * srcHeight
		d.Dot = fixed.P(srcX, srcY + baseLine)
		d.DrawString(string(bary))
		//d.DrawBytes(bary)
		res = append(res, Info {
			Bin:  bary,
			SrcX: srcX,
			SrcY: srcY,
			Han: len(bary) == 1,
		})
	}
	if e1 := canvas.ReplacePixels(dst.Pix); e1 != nil {
		e = fmt.Errorf("フォント描画エラー %s", e1.Error())
		return
	} else {
		return
	}
}

func BinPrint(b []byte) {
	cary := []string{}
	for _, c := range b {
		cary = append(cary, fmt.Sprintf("%x", c))
	}
	fmt.Printf("%s [%s]\n", string(b), strings.Join(cary, ","))
}

/*
type Metrics struct {
	// Height is the recommended amount of vertical space between two lines of
	// text.
	Height fixed.Int26_6

	// Ascent is the distance from the top of a line to its baseline.
	Ascent fixed.Int26_6

	// Descent is the distance from the bottom of a line to its baseline. The
	// value is typically positive, even though a descender goes below the
	// baseline.
	Descent fixed.Int26_6
}
	//addBottom := (met.Ascent.Floor() - met.Descent.Floor()) + (met.Ascent.Floor() - met.Height.Floor())
	// {24:00 21:08 2:57} -> + 20 || srcy + 3
	// {24:00 25:52 7:44} -> + 19
	//addBottom := met.Height.Floor()
	//diff := met.Height.Floor() - met.Ascent.Floor()
	//diff := met.Descent.Floor()
*/

func New(fileName string, size int) (*FontMap, error) {
	if tt, e1 := LoadFont(fileName); e1 != nil {
		return nil, e1
	} else {
		face := GetTTFace(size, tt)
		met := face.Metrics()
		srcWidth  := size

		srcHeight := met.Ascent.Ceil() + met.Descent.Ceil()
		baseLine  := met.Ascent.Ceil()
/*
		srcHeight := met.Ascent.Floor() + met.Descent.Floor()
		baseLine  := met.Ascent.Floor()

		srcHeight := met.Ascent.Round() + met.Descent.Round()
		baseLine  := met.Ascent.Round()*/

		if canvas, e2 := CreateCanvas(srcWidth, srcHeight); e2 != nil {
			return nil, e2
		} else {
			w, h := canvas.Size()
			dst := image.NewRGBA(image.Rect(0, 0, w, h))
			drawer := GetDrawer(dst, face)
			if info, e3 := Draw(srcWidth, srcHeight, baseLine, canvas, dst, drawer); e3 != nil {
				return nil, e3
			} else {
				return &FontMap {
					SrcWidth:  srcWidth,
					SrcHeight: srcHeight,
					Canvas: canvas,
					Info:   info,
				}, nil
			}
		}
	}
}

func (me *FontMap) Search(x []byte) (*Info) {
	length := len(me.Info)
	i := sort.Search(length, func(i int) (bool) {
		return bytes.Compare(me.Info[i].Bin, x) >= 0
	})
	if i < length && bytes.Equal(me.Info[i].Bin, x) {
		return &me.Info[i]
	} else {
		return nil
	}
}

type FontSprite struct {
	fig.IntPoint
	*Info
}

type FontSprites struct {
	SrcWidth int
	SrcHeight int
	X	int
	Y	int
	FontArray []FontSprite
	*FontMap
	PreText string
}

func (me *FontSprites) SetSub(s string, yIndex int) {
	xPos := 0
	for _, v := range s {
		if info := me.FontMap.Search([]byte(string(v))); info != nil {
			me.FontArray = append(me.FontArray, FontSprite {
				IntPoint: fig.IntPoint {
					X: xPos,
					Y: yIndex * me.SrcHeight,
				},
				Info: info,
			})
			if info.Han {
				xPos += me.SrcWidth / 2
			} else {
				xPos += me.SrcWidth
			}
		} else {
			me.FontArray = append(me.FontArray, FontSprite {
				IntPoint: fig.IntPoint {
					X: xPos,
					Y: yIndex * me.SrcHeight,
				},
			})
			xPos += me.SrcWidth
		}
	}
}

func (me *FontSprites) Set(s string) {
	if s == me.PreText {
		return
	} else {
		me.PreText = s
		me.FontArray = nil
		strAry := strings.Split(s, "\n")
		for i, str := range strAry {
			me.SetSub(str, i)
		}
	}
}

func (me *FontSprites) Src(i int) (x0, y0, x1, y1 int) {
	fsp := me.FontArray[i]
	if fsp.Info == nil {
		return
	} else {
		return fsp.Info.SrcX, fsp.Info.SrcY, fsp.Info.SrcX + me.SrcWidth, fsp.Info.SrcY + me.SrcHeight
	}
}

func (me *FontSprites) Dst(i int) (x0, y0, x1, y1 int) {
	fsp := me.FontArray[i]
	if fsp.Info == nil {
		return
	} else {
		x := me.X + fsp.IntPoint.X
		y := me.Y + fsp.IntPoint.Y
		return x, y, x + me.SrcWidth, y + me.SrcHeight
	}
}

func (me *FontSprites) Len() (int) {
	return len(me.FontArray)
}

func (me *FontMap) Draw(screen *ebiten.Image, x, y int, txt string) {
	fs := FontSprites {
		SrcWidth:  me.SrcWidth,
		SrcHeight: me.SrcHeight,
		X: x,
		Y: y,
		FontMap: me,
	}
	fs.Set(txt)
	screen.DrawImage(me.Canvas, &ebiten.DrawImageOptions{
		ImageParts: &fs,
	})
}

