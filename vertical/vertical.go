package vertical

import (
	"bytes"
	"io"
	"log"
	"reflect"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/noborus/pgsp/str"
)

type Vertical struct {
	out     io.Writer
	bar     rune
	padding int
	header  []string
	rows    [][]interface{}
}

func NewWriter(writer io.Writer) *Vertical {
	v := &Vertical{
		out:     writer,
		bar:     '|',
		padding: 1,
		header:  []string{},
		rows:    [][]interface{}{},
	}
	return v
}

func (v *Vertical) SetHeader(header []string) {
	v.header = header
}

func (v *Vertical) SetPadding(p int) {
	v.padding = p
}

func (v *Vertical) SetBar(b rune) {
	v.bar = b
}

func (v *Vertical) AppendStruct(value interface{}) {
	rf := reflect.TypeOf(value)
	rv := reflect.ValueOf(value)

	num := rf.NumField()
	row := make([]interface{}, num)
	for i := 0; i < num; i++ {
		row[i] = rv.Field(i).Interface()
	}
	v.rows = append(v.rows, row)
}

func (v *Vertical) Append(values []interface{}) {
	v.rows = append(v.rows, values)
}

func (v *Vertical) Render() {
	maxH := 0
	for _, h := range v.header {
		hlen := runewidth.StringWidth(h)
		if hlen > maxH {
			maxH = hlen
		}
	}

	buff := new(bytes.Buffer)
	for _, row := range v.rows {
		for n, r := range row {
			h := ""
			if len(v.header) > n {
				h = v.header[n]
			}
			buff.WriteString(strings.Repeat(" ", v.padding))
			buff.WriteString(h)
			buff.WriteString(strings.Repeat(" ", maxH-len(h)))
			buff.WriteString(" ")
			buff.WriteString(string(v.bar))
			buff.WriteString(" ")
			buff.WriteString(str.ToStr(r))
			buff.WriteString("\n")
		}
	}
	if _, err := io.WriteString(v.out, buff.String()); err != nil {
		log.Panicln(err)
	}
}
