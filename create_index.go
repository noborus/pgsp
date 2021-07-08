package pgsp

import (
	"bytes"
	"database/sql"
	"fmt"

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
	LockersPid      int64  `db:"current_locker_pid"`
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
	pid := fmt.Sprintf("%v", v.PID)
	datid := fmt.Sprintf("%v", v.DATID)
	relid := fmt.Sprintf("%v", v.RELID)
	indexRelid := fmt.Sprintf("%v", v.IndexRelid)
	lockerTotal := fmt.Sprintf("%v", v.LockersTotal)
	lockerDone := fmt.Sprintf("%v", v.LockersDone)
	lockerPid := fmt.Sprintf("%v", v.LockersPid)
	blocksTotal := fmt.Sprintf("%v", v.BlocksTotal)
	blocksDone := fmt.Sprintf("%v", v.BlocksDone)
	tuplesTotal := fmt.Sprintf("%v", v.TuplesTotal)
	tuplesDone := fmt.Sprintf("%v", v.TuplesDone)
	partitionsTotal := fmt.Sprintf("%v", v.PartitionsTotal)
	partitionsDone := fmt.Sprintf("%v", v.PartitionsDone)
	return []string{pid, datid, v.DATNAME, relid, indexRelid, v.Command, v.PHASE, lockerTotal, lockerDone, lockerPid, blocksTotal, blocksDone, tuplesTotal, tuplesDone, partitionsTotal, partitionsDone}
}

func (v CreateIndex) Table() string {
	buff := new(bytes.Buffer)
	t := tablewriter.NewWriter(buff)
	t.SetHeader(CreateIndexColumns)
	t.Append(v.String())
	t.Render()
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
