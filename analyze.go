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

// pg_stat_progress_analyze.
type Analyze struct {
	PID                    int     `db:"pid"`
	DATID                  int     `db:"datid"`
	DATNAME                string  `db:"datname"`
	RELID                  int     `db:"relid"`
	PHASE                  string  `db:"phase"`
	SampleBLKSTotal        int64   `db:"sample_blks_total"`
	SampleBLKSScanned      int64   `db:"sample_blks_scanned"`
	ExtStatsTotal          int64   `db:"ext_stats_total"`
	ExtStatsComputed       int64   `db:"ext_stats_computed"`
	ChildTablesTotal       int64   `db:"child_tables_total"`
	ChildTablesDone        int64   `db:"child_tables_done"`
	CurrentChildTableRelid int     `db:"current_child_table_relid"`
	DelayTime              float64 `db:"delay_time"`
}

var (
	AnalyzeTableName = "pg_stat_progress_analyze"
	AnalyzeQuery     string
	AnalyzeColumns   []string
)

func GetAnalyze(ctx context.Context, db *sqlx.DB) ([]Progress, error) {
	if len(AnalyzeColumns) == 0 {
		AnalyzeColumns = getColumns(Analyze{})
	}
	if AnalyzeQuery == "" {
		AnalyzeQuery = buildQuery(AnalyzeTableName, AnalyzeColumns)
	}
	return selectAnalyze(ctx, db, AnalyzeQuery)
}

func selectAnalyze(ctx context.Context, db *sqlx.DB, query string) ([]Progress, error) {
	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var as []Progress
	for rows.Next() {
		var row Analyze
		err = rows.StructScan(&row)
		if err != nil {
			return nil, err
		}
		as = append(as, row)
	}
	return as, rows.Err()
}

func (v Analyze) Name() string {
	return AnalyzeTableName
}

func (v Analyze) Pid() int {
	return v.PID
}

func (v Analyze) Color() (color.Color, color.Color) {
	return color.RGBA{R: 255, G: 124, B: 203}, color.RGBA{R: 253, G: 255, B: 140}
}

func (v Analyze) Table() string {
	value := str.ToStrStruct(v)
	buff := new(bytes.Buffer)

	t := tablewriter.NewWriter(buff)
	t.Header(AnalyzeColumns)
	t.Append(value)
	t.Render()

	return buff.String()
}

func (v Analyze) Progress() float64 {
	if v.ChildTablesTotal != 0 {
		return float64(v.ChildTablesDone) / float64(v.ChildTablesTotal)
	}
	if v.ExtStatsTotal != 0 {
		return float64(v.ExtStatsComputed) / float64(v.ExtStatsTotal)
	}
	return float64(v.SampleBLKSScanned) / float64(v.SampleBLKSTotal)
}
