package vertical

import (
	"bytes"
	"io"
	"testing"
)

func TestVertical_Render(t *testing.T) {
	type fields struct {
		out    io.Writer
		bar    rune
		header []string
		rows   [][]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test1",
			fields: fields{
				bar:    '|',
				header: []string{"a", "ab"},
				rows: [][]interface{}{
					{"a", "b"},
				},
			},
			want: `a  | a
ab | b
`,
		},
		{
			name: "testLongHeader",
			fields: fields{
				bar:    '|',
				header: []string{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "b"},
				rows: [][]interface{}{
					{"a", "b"},
				},
			},
			want: `aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa | a
b                                    | b
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := new(bytes.Buffer)
			v := &Vertical{
				out:    writer,
				bar:    tt.fields.bar,
				header: tt.fields.header,
				rows:   tt.fields.rows,
			}
			v.Render()
			if tt.want != writer.String() {
				t.Errorf("Vertical.Render() not match\n[%v]\n, want \n[%v]\n", writer.String(), tt.want)
			}
		})
	}
}
