package pgsp

import (
	"bytes"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/olekukonko/tablewriter"
)

// pg_stat_progress_analyze
type Analyze struct {
	PID                    int    `db:"pid"`
	DATID                  int    `db:"datid"`
	DATNAME                string `db:"datname"`
	RELID                  int    `db:"relid"`
	PHASE                  string `db:"phase"`
	SampleBLKSTotal        int64  `db:"sample_blks_total"`
	SampleBLKSScanned      int64  `db:"sample_blks_scanned"`
	ExtStatsTotal          int64  `db:"ext_stats_total"`
	ExtStatsComputed       int64  `db:"ext_stats_computed"`
	ChildTablesTotal       int64  `db:"child_tables_total"`
	ChildTablesDone        int64  `db:"child_tables_done"`
	CurrentChildTableRelid int    `db:"current_child_table_relid"`
}

var AnalyzeColumns = []string{
	"pid",
	"datid",
	"datname",
	"relid",
	"phase",
	"sample_blks_total",
	"sample_blks_scanned",
	"ext_stats_total",
	"ext_stats_computed",
	"child_tables_total",
	"child_tables_done",
	"current_child_table_relid",
}

func GetAnalyze(db *sql.DB) ([]Analyze, error) {
	// tableName := "analyze_progress"
	tableName := "pg_stat_progress_analyze"
	query := buildQuery(tableName, AnalyzeColumns)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	var as []Analyze
	for rows.Next() {
		var row Analyze
		err = rows.Scan(&row.PID, &row.DATID, &row.DATNAME, &row.RELID, &row.PHASE, &row.SampleBLKSTotal, &row.SampleBLKSScanned, &row.ExtStatsTotal, &row.ExtStatsComputed, &row.ChildTablesTotal, &row.ChildTablesDone, &row.CurrentChildTableRelid)
		if err != nil {
			return nil, err
		}
		as = append(as, row)
	}
	return as, nil
}

func (v Analyze) String() []string {
	pid := fmt.Sprintf("%v", v.PID)
	datid := fmt.Sprintf("%v", v.DATID)
	relid := fmt.Sprintf("%v", v.RELID)
	total := fmt.Sprintf("%v", v.SampleBLKSTotal)
	scanned := fmt.Sprintf("%v", v.SampleBLKSScanned)
	extTotal := fmt.Sprintf("%v", v.ExtStatsTotal)
	extComputed := fmt.Sprintf("%v", v.ExtStatsComputed)
	childTotal := fmt.Sprintf("%v", v.ChildTablesTotal)
	childDone := fmt.Sprintf("%v", v.ChildTablesDone)
	childRelid := fmt.Sprintf("%v", v.CurrentChildTableRelid)
	return []string{pid, datid, v.DATNAME, relid, v.PHASE, total, scanned, extTotal, extComputed, childTotal, childDone, childRelid}
}

func (v Analyze) Table() string {
	buff := new(bytes.Buffer)
	t := tablewriter.NewWriter(buff)
	t.SetHeader(VacuumColumns)
	t.Append(v.String())
	t.Render()
	return buff.String()
}

func (v Analyze) Name() string {
	return "pg_stat_progress_analyze"
}

func (v Analyze) Progress() float64 {
	return float64(v.SampleBLKSScanned) / float64(v.SampleBLKSTotal)
}
