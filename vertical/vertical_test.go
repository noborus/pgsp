package vertical

import (
	"bytes"
	"database/sql"
	"testing"
)

func TestVertical_Render(t *testing.T) {
	type fields struct {
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
			name: "testNullInt64",
			fields: fields{
				bar:    '|',
				header: []string{"a", "ab"},
				rows: [][]interface{}{
					{sql.NullInt64{Int64: 1, Valid: true}, "b"},
				},
			},
			want: `a  | 1
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

func TestVertical_AppendStruct(t *testing.T) {
	type test struct {
		PID   int    `db:"pid"`
		PHASE string `db:"phase"`
	}
	tests := []struct {
		name string
		args interface{}
		want string
	}{
		{
			name: "testStruct1",
			args: test{
				PID:   1,
				PHASE: "test",
			},
			want: ` a | 1
 b | test
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := new(bytes.Buffer)
			v := NewWriter(writer)
			v.SetHeader([]string{"a", "b"})
			v.AppendStruct(tt.args)
			v.Render()
			if tt.want != writer.String() {
				t.Errorf("Vertical.Render() not match\n[%v]\n, want \n[%v]\n", writer.String(), tt.want)
			}
		})
	}
}
