package pgsp

import (
	"testing"
)

func TestCopy_Vertical(t *testing.T) {
	type fields struct {
		PID             int
		DATID           int
		DATNAME         string
		RELID           int
		COMMAND         string
		CTYPE           string
		BYTESProcessed  int64
		BYTESTotal      int64
		TUPLESProcessed int64
		TUPLESExcluded  int64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test1",
			fields: fields{
				PID:             1,
				DATID:           1,
				DATNAME:         "name",
				RELID:           1,
				COMMAND:         "command",
				CTYPE:           "ctype",
				BYTESProcessed:  1,
				BYTESTotal:      10,
				TUPLESProcessed: 1,
				TUPLESExcluded:  10,
			},
			want: ` pid              | 1
 datid            | 1
 datname          | name
 relid            | 1
 command          | command
 type             | ctype
 bytes_processed  | 1
 bytes_total      | 10
 tuples_processed | 1
 tuples_excluded  | 10
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Copy{
				PID:             tt.fields.PID,
				DATID:           tt.fields.DATID,
				DATNAME:         tt.fields.DATNAME,
				RELID:           tt.fields.RELID,
				COMMAND:         tt.fields.COMMAND,
				CTYPE:           tt.fields.CTYPE,
				BYTESProcessed:  tt.fields.BYTESProcessed,
				BYTESTotal:      tt.fields.BYTESTotal,
				TUPLESProcessed: tt.fields.TUPLESProcessed,
				TUPLESExcluded:  tt.fields.TUPLESExcluded,
			}
			CopyColumns = getColumns(Copy{})
			if got := v.Vertical(); got != tt.want {
				t.Errorf("Copy.Vertical() = %v, want %v", got, tt.want)
			}
		})
	}
}
