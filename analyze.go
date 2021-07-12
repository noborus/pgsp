package pgsp

import (
	"bytes"
	"database/sql"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/noborus/pgsp/vertical"

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

var AnalyzeTableName = "pg_stat_progress_analyze"

func GetAnalyze(db *sql.DB) ([]Analyze, error) {
	query := buildQuery(AnalyzeTableName, AnalyzeColumns)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var as []Analyze
	for rows.Next() {
		var row Analyze
		err = rows.Scan(
			&row.PID,
			&row.DATID,
			&row.DATNAME,
			&row.RELID,
			&row.PHASE,
			&row.SampleBLKSTotal,
			&row.SampleBLKSScanned,
			&row.ExtStatsTotal,
			&row.ExtStatsComputed,
			&row.ChildTablesTotal,
			&row.ChildTablesDone,
			&row.CurrentChildTableRelid,
		)
		if err != nil {
			return nil, err
		}
		as = append(as, row)
	}
	return as, rows.Err()
}

func (v Analyze) String() []string {
	return []string{
		strconv.Itoa(v.PID),
		strconv.Itoa(v.DATID),
		v.DATNAME,
		strconv.Itoa(v.RELID),
		v.PHASE,
		strconv.FormatInt(v.SampleBLKSTotal, 10),
		strconv.FormatInt(v.SampleBLKSScanned, 10),
		strconv.FormatInt(v.ExtStatsTotal, 10),
		strconv.FormatInt(v.ExtStatsComputed, 10),
		strconv.FormatInt(v.ChildTablesTotal, 10),
		strconv.FormatInt(v.ChildTablesDone, 10),
		strconv.Itoa(v.CurrentChildTableRelid),
	}
}

func (v Analyze) Table() string {
	value := v.String()
	buff := new(bytes.Buffer)

	t := tablewriter.NewWriter(buff)
	t.SetHeader(AnalyzeColumns[0:7])
	t.Append(value[0:7])
	t.Render()

	t2 := tablewriter.NewWriter(buff)
	t2.SetHeader(AnalyzeColumns[7:])
	t2.Append(value[7:])
	t2.Render()

	return buff.String()
}

func (v Analyze) Vertical() string {
	buff := new(bytes.Buffer)
	vt := vertical.NewWriter(buff)
	vt.SetHeader(AnalyzeColumns)
	vt.Append([]interface{}{
		v.PID,
		v.DATID,
		v.DATNAME,
		v.RELID,
		v.PHASE,
		v.SampleBLKSScanned,
		v.ExtStatsTotal,
		v.ExtStatsComputed,
		v.ChildTablesTotal,
		v.ChildTablesDone,
		v.CurrentChildTableRelid,
	})
	vt.Render()
	return buff.String()
}

func (v Analyze) Name() string {
	return AnalyzeTableName
}

func (v Analyze) Progress() float64 {
	return float64(v.SampleBLKSScanned) / float64(v.SampleBLKSTotal)
}

func (v Analyze) Pid() int {
	return v.PID
}
