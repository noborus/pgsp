package pgsp

import (
	"bytes"
	"context"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/noborus/pgsp/str"
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

var (
	VacuumTableName = "pg_stat_progress_vacuum"
	VacuumQuery     string
	VacuumColumns   []string
)

func GetVacuum(ctx context.Context, db *sqlx.DB) ([]Progress, error) {
	if len(VacuumColumns) == 0 {
		VacuumColumns = getColumns(Vacuum{})
	}
	if VacuumQuery == "" {
		VacuumQuery = buildQuery(VacuumTableName, VacuumColumns)
	}
	return selectVacuum(ctx, db, VacuumQuery)
}

func selectVacuum(ctx context.Context, db *sqlx.DB, query string) ([]Progress, error) {
	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var as []Progress
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

func (v Vacuum) Color() (string, string) {
	return "#5A56E0", "#FF7CCB"
}

func (v Vacuum) Table() string {
	value := str.ToStrStruct(v)
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
