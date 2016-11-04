
package ttpl

import (
	"strconv"
	"strings"
	"fmt"
)

type Ttpl struct {
	Template string
	Values   map[string]string
}

func (me *Ttpl) SetInt(key string, value int) {
	me.Values[key] = strconv.FormatInt(int64(value), 10)
}

func (me *Ttpl) SetFloat(key string, value float64) {
	//me.Values[key] = strconv.FormatFloat(value, 'f', -1, 64)
	me.Values[key] = fmt.Sprintf("%0.2f", value)
}

func (me *Ttpl) SetStr(key string, value string) {
	me.Values[key] = value
}

func (me *Ttpl) Text() (string) {
	s := me.Template
	for k, v := range me.Values {
		target := "#" + k + "#"
		s = strings.Replace(s, target, v, 1)
	}
	return s
}

func New(tpl string) (*Ttpl) {
	return &Ttpl {
		Template: tpl,
		Values:   make(map[string]string),
	}
}
