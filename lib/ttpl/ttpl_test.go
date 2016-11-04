
package ttpl

import (
	"testing"
	"fmt"
)



func Test(t *testing.T) {
	tpl := `
Frame #FrameCount#

Objects #ObjectCount#

Boss #BossEndurance#

Score #Score#
`
	x := New(tpl)
	x.SetFloat("FrameCount", 55.8)
	x.SetInt("ObjectCount", 120)
	x.SetInt("BossEndurance", 350)
	x.SetInt("Score", 12500)

	fmt.Println("template", tpl)
	fmt.Println("-----------------------")
	fmt.Println(x.Text())
}
