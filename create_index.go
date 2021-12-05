package pgsp

import (
	"bytes"
	"context"
	"log"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/noborus/pgsp/vertical"
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
	PartitionsDone  int64  `db:"partitions_done"`
}

var CreateIndexTableName = "pg_stat_progress_create_index"

var CreateIndexQuery string
var CreateIndexColumns []string

func GetCreateIndex(ctx context.Context, db *sqlx.DB) ([]PGSProgress, error) {
	if len(CreateIndexColumns) == 0 {
		CreateIndexColumns = getColumns(CreateIndex{})
	}
	if CreateIndexQuery == "" {
		CreateIndexQuery = buildQuery(CreateIndexTableName, CreateIndexColumns)
		log.Println(CreateIndexQuery)
	}
	query := buildQuery(CreateIndexTableName, CreateIndexColumns)
	return selectCreateIndex(ctx, db, query)
}

func selectCreateIndex(ctx context.Context, db *sqlx.DB, query string) ([]PGSProgress, error) {
	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var as []PGSProgress
	for rows.Next() {
		var row CreateIndex
		err = rows.StructScan(&row)
		if err != nil {
			return nil, err
		}
		as = append(as, row)
	}
	return as, rows.Err()
}

func (v CreateIndex) Name() string {
	return CreateIndexTableName
}

func (v CreateIndex) Pid() int {
	return v.PID
}

func (v CreateIndex) Color() (string, string) {
	return "#EE6FF8", "#5A56E0"
}

func (v CreateIndex) Table() string {
	value := []string{
		strconv.Itoa(v.PID),
		strconv.Itoa(v.DATID),
		v.DATNAME,
		strconv.Itoa(v.RELID),
		strconv.Itoa(v.IndexRelid),
		v.Command,
		v.PHASE,
		strconv.FormatInt(v.LockersTotal, 10),
		strconv.FormatInt(v.LockersDone, 10),
		strconv.FormatInt(v.LockersPid, 10),
		strconv.FormatInt(v.BlocksTotal, 10),
		strconv.FormatInt(v.BlocksDone, 10),
		strconv.FormatInt(v.TuplesTotal, 10),
		strconv.FormatInt(v.TuplesDone, 10),
		strconv.FormatInt(v.PartitionsTotal, 10),
		strconv.FormatInt(v.PartitionsDone, 10),
	}

	buff := new(bytes.Buffer)

	t := tablewriter.NewWriter(buff)
	t.SetHeader(CreateIndexColumns[0:9])
	t.Append(value[0:9])
	t.Render()

	t2 := tablewriter.NewWriter(buff)
	t2.SetHeader(CreateIndexColumns[9:])
	t2.Append(value[9:])
	t2.Render()

	return buff.String()
}

func (v CreateIndex) Vertical() string {
	buff := new(bytes.Buffer)
	vt := vertical.NewWriter(buff)
	vt.SetHeader(CreateIndexColumns)
	vt.AppendStruct(v)
	vt.Render()
	return buff.String()
}

func (v CreateIndex) Progress() float64 {
	if v.BlocksTotal != 0 {
		return float64(v.BlocksDone) / float64(v.BlocksTotal)
	}
	if v.PartitionsTotal != 0 {
		return float64(v.PartitionsDone) / float64(v.PartitionsTotal)
	}
	return float64(v.TuplesDone) / float64(v.TuplesTotal)
}
