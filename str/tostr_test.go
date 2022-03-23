package str_test

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/noborus/pgsp"
	"github.com/noborus/pgsp/str"
)

func TestToStrStruct(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "test1",
			args: args{
				value: pgsp.BaseBackup{
					PID:                 1,
					PHASE:               "t",
					BackupTotal:         sql.NullInt64{Int64: 1, Valid: true},
					BackupStreamed:      1,
					TablespacesTotal:    1,
					TablespacesStreamed: 1,
				},
			},
			want: []string{"1", "t", "1", "1", "1", "1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := str.ToStrStruct(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToStrStruct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToStr(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				v: 1,
			},
			want: "1",
		},
		{
			name: "testSQLNullInt64",
			args: args{
				v: sql.NullInt64{Int64: 1, Valid: true},
			},
			want: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := str.ToStr(tt.args.v); got != tt.want {
				t.Errorf("ToStr() = %v, want %v", got, tt.want)
			}
		})
	}
}
