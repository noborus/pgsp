package pgsp

import (
	"bytes"
	"context"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/noborus/pgsp/vertical"
	"github.com/olekukonko/tablewriter"
)

// pg_stat_progress_vacuum
type Vacuum struct {
	PID              int    `db:"pid"`
	DATID            int    `db:"datid"`
	DATNAME          string `db:"datname"`
	RELID            int    `db:"relid"`
	PHASE            string `db:"phase"`
	HeapBLKSTotal    int64  `db:"heap_blks_total"`
	HeapBLKSScanned  int64  `db:"heap_blks_scanned"`
	HeapBLKSVacuumed int64  `db:"heap_blks_vacuumed"`
	IndexVacuumCount int64  `db:"index_vacuum_count"`
	MaxDeadTuples    int64  `db:"max_dead_tuples"`
	NumDeadTuples    int64  `db:"num_dead_tuples"`
}

var VacuumTableName = "pg_stat_progress_vacuum"
var VacuumQuery string
var VacuumColumns []string

func GetVacuum(ctx context.Context, db *sqlx.DB) ([]Vacuum, error) {
	if len(VacuumColumns) == 0 {
		VacuumColumns = getColumns(Vacuum{})
	}
	if VacuumQuery == "" {
		VacuumQuery = buildQuery(VacuumTableName, VacuumColumns)
	}
	query := buildQuery(VacuumTableName, VacuumColumns)
	return selectVacuum(ctx, db, query)
}

func selectVacuum(ctx context.Context, db *sqlx.DB, query string) ([]Vacuum, error) {
	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var as []Vacuum
	for rows.Next() {
		var row Vacuum
		err = rows.StructScan(&row)
		if err != nil {
			return nil, err
		}
		as = append(as, row)
	}
	return as, nil
}

func (v Vacuum) Name() string {
	return VacuumTableName
}

func (v Vacuum) Pid() int {
	return v.PID
}

func (v Vacuum) Table() string {
	value := []string{
		strconv.Itoa(v.PID),
		strconv.Itoa(v.DATID),
		v.DATNAME,
		strconv.Itoa(v.RELID),
		v.PHASE,
		strconv.FormatInt(v.HeapBLKSTotal, 10),
		strconv.FormatInt(v.HeapBLKSScanned, 10),
		strconv.FormatInt(v.HeapBLKSVacuumed, 10),
		strconv.FormatInt(v.IndexVacuumCount, 10),
		strconv.FormatInt(v.MaxDeadTuples, 10),
		strconv.FormatInt(v.NumDeadTuples, 10),
	}

	buff := new(bytes.Buffer)
	t := tablewriter.NewWriter(buff)
	t.SetHeader(VacuumColumns[0:7])
	t.Append(value[0:7])
	t.Render()

	t2 := tablewriter.NewWriter(buff)
	t2.SetHeader(VacuumColumns[7:])
	t2.Append(value[7:])
	t2.Render()
	return buff.String()
}

func (v Vacuum) Vertical() string {
	buff := new(bytes.Buffer)
	vt := vertical.NewWriter(buff)
	vt.SetHeader(VacuumColumns)
	vt.AppendStruct(v)
	vt.Render()
	return buff.String()
}

func (v Vacuum) Progress() float64 {
	return float64(v.HeapBLKSScanned) / float64(v.HeapBLKSTotal)
}
