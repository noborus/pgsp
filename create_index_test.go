package pgsp

import (
	"testing"

	_ "github.com/lib/pq"
)

func TestCreateIndex_Table(t *testing.T) {
	type fields struct {
		PID             int
		DATID           int
		DATNAME         string
		RELID           int
		IndexRelid      int
		Command         string
		PHASE           string
		LockersTotal    int64
		LockersDone     int64
		LockersPid      int
		BlocksTotal     int64
		BlocksDone      int64
		TuplesTotal     int64
		TuplesDone      int64
		PartitionsTotal int64
		PartitionsDone  int64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "test1",
			fields: fields{},
			want: `+-----+-------+---------+-------+-------------+---------+-------+---------------+--------------+
| PID | DATID | DATNAME | RELID | INDEX RELID | COMMAND | PHASE | LOCKERS TOTAL | LOCKERS DONE |
+-----+-------+---------+-------+-------------+---------+-------+---------------+--------------+
|   0 |     0 |         |     0 |           0 |         |       |             0 |            0 |
+-----+-------+---------+-------+-------------+---------+-------+---------------+--------------+
+--------------------+--------------+-------------+--------------+-------------+------------------+-----------------+
| CURRENT LOCKER PID | BLOCKS TOTAL | BLOCKS DONE | TUPLES TOTAL | TUPLES DONE | PARTITIONS TOTAL | PARTITIONS DONE |
+--------------------+--------------+-------------+--------------+-------------+------------------+-----------------+
|                  0 |            0 |           0 |            0 |           0 |                0 |               0 |
+--------------------+--------------+-------------+--------------+-------------+------------------+-----------------+
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := CreateIndex{
				PID:             tt.fields.PID,
				DATID:           tt.fields.DATID,
				DATNAME:         tt.fields.DATNAME,
				RELID:           tt.fields.RELID,
				IndexRelid:      tt.fields.IndexRelid,
				Command:         tt.fields.Command,
				PHASE:           tt.fields.PHASE,
				LockersTotal:    tt.fields.LockersTotal,
				LockersDone:     tt.fields.LockersDone,
				LockersPid:      tt.fields.LockersPid,
				BlocksTotal:     tt.fields.BlocksTotal,
				BlocksDone:      tt.fields.BlocksDone,
				TuplesTotal:     tt.fields.TuplesTotal,
				TuplesDone:      tt.fields.TuplesDone,
				PartitionsTotal: tt.fields.PartitionsTotal,
				PartitionsDone:  tt.fields.PartitionsDone,
			}
			if got := v.Table(); got != tt.want {
				t.Errorf("CreateIndex.Table() = %v, want %v", got, tt.want)
			}
		})
	}
}
