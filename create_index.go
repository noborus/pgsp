package pgsp

import (
	"bytes"
	"database/sql"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/olekukonko/tablewriter"
)

// pg_stat_progress_create_index
type CreateIndex struct {
	PID             int    `db:"pid"`
	DATID           int    `db:"datid"`
	DATNAME         string `db:"datname"`
	RELID           int    `db:"relid"`
	IndexRelid      int    `db:"index_relid"`
	Command         string `db:"command"`
	PHASE           string `db:"phase"`
	LockersTotal    int64  `db:"lockers_total"`
	LockersDone     int64  `db:"lockers_done"`
	LockersPid      int    `db:"current_locker_pid"`
	BlocksTotal     int64  `db:"blocks_total"`
	BlocksDone      int64  `db:"blocks_done"`
	TuplesTotal     int64  `db:"tuples_total"`
	TuplesDone      int64  `db:"tuples_done"`
	PartitionsTotal int64  `db:"partitions_total"`
	PartitionsDone  int64  `db:"partitions__done"`
}

var CreateIndexColumns = []string{
	"pid",
	"datid",
	"datname",
	"relid",
	"index_relid",
	"command",
	"phase",
	"lockers_total",
	"lockers_done",
	"current_locker_pid",
	"blocks_total",
	"blocks_done",
	"tuples_total",
	"tuples_done",
	"partitions_total",
	"partitions_done",
}

func GetCreateIndex(db *sql.DB) ([]CreateIndex, error) {
	// tableName := "index_progress"
	tableName := "pg_stat_progress_create_index"
	query := buildQuery(tableName, CreateIndexColumns)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	var as []CreateIndex
	for rows.Next() {
		var row CreateIndex
		err = rows.Scan(
			&row.PID,
			&row.DATID,
			&row.DATNAME,
			&row.RELID,
			&row.IndexRelid,
			&row.Command,
			&row.PHASE,
			&row.LockersTotal,
			&row.LockersDone,
			&row.LockersPid,
			&row.BlocksTotal,
			&row.BlocksDone,
			&row.TuplesTotal,
			&row.TuplesDone,
			&row.PartitionsTotal,
			&row.PartitionsDone)
		if err != nil {
			return nil, err
		}
		as = append(as, row)
	}
	return as, nil
}

func (v CreateIndex) String() []string {
	return []string{
		strconv.Itoa(v.PID),
		strconv.Itoa(v.DATID),
		v.DATNAME,
		strconv.Itoa(v.RELID),
		strconv.Itoa(v.IndexRelid),
		v.Command,
		v.PHASE,
		strconv.FormatInt(v.LockersTotal, 10),
		strconv.FormatInt(v.LockersDone, 10),
		strconv.Itoa(v.LockersPid),
		strconv.FormatInt(v.BlocksTotal, 10),
		strconv.FormatInt(v.BlocksDone, 10),
		strconv.FormatInt(v.TuplesTotal, 10),
		strconv.FormatInt(v.TuplesDone, 10),
		strconv.FormatInt(v.PartitionsTotal, 10),
		strconv.FormatInt(v.PartitionsDone, 10),
	}
}

func (v CreateIndex) Table() string {
	value := v.String()
	buff := new(bytes.Buffer)

	t := tablewriter.NewWriter(buff)
	t.SetHeader(CreateIndexColumns[0:6])
	t.Append(value[0:6])
	t.Render()

	t2 := tablewriter.NewWriter(buff)
	t2.SetHeader(CreateIndexColumns[6:])
	t2.Append(value[6:])
	t2.Render()

	return buff.String()
}

func (v CreateIndex) Name() string {
	return "pg_stat_progress_create_index"
}

func (v CreateIndex) Progress() float64 {
	if v.BlocksTotal != 0 {
		return float64(v.BlocksDone) / float64(v.BlocksTotal)
	}
	return float64(v.TuplesDone) / float64(v.TuplesTotal)
}

func (v CreateIndex) Pid() int {
	return v.PID
}
