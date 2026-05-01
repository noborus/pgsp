package pgsp

import (
	"bytes"
	"context"
	"image/color"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/noborus/pgsp/str"
	"github.com/olekukonko/tablewriter"
)

// pg_stat_progress_vacuum
type Vacuum struct {
	PID               int     `db:"pid"`
	DATID             int     `db:"datid"`
	DATNAME           string  `db:"datname"`
	RELID             int     `db:"relid"`
	PHASE             string  `db:"phase"`
	HeapBLKSTotal     int64   `db:"heap_blks_total"`
	HeapBLKSScanned   int64   `db:"heap_blks_scanned"`
	HeapBLKSVacuumed  int64   `db:"heap_blks_vacuumed"`
	IndexVacuumCount  int64   `db:"index_vacuum_count"`
	MaxDeadTupleBytes int64   `db:"max_dead_tuple_bytes"`
	DeadTupleBytes    int64   `db:"dead_tuple_bytes"`
	NumDeadItemIDs    int64   `db:"num_dead_item_ids"`
	IndexesTotal      int64   `db:"indexes_total"`
	IndexesProcessed  int64   `db:"indexes_processed"`
	DelayTime         float64 `db:"delay_time"`
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

func (v Vacuum) Color() (color.Color, color.Color) {
	return color.RGBA{R: 90, G: 86, B: 224}, color.RGBA{R: 255, G: 124, B: 203}
}

func (v Vacuum) Table() string {
	value := str.ToStrStruct(v)
	buff := new(bytes.Buffer)
	t := tablewriter.NewWriter(buff)
	t.Header(VacuumColumns)
	t.Append(value)
	t.Render()
	return buff.String()
}

func (v Vacuum) Progress() float64 {
	return float64(v.HeapBLKSScanned) / float64(v.HeapBLKSTotal)
}
