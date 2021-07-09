package pgsp

import (
	"bytes"
	"database/sql"
	"strconv"

	_ "github.com/lib/pq"
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

var VacuumColumns = []string{
	"pid",
	"datid",
	"datname",
	"relid",
	"phase",
	"heap_blks_total",
	"heap_blks_scanned",
	"heap_blks_vacuumed",
	"index_vacuum_count",
	"max_dead_tuples",
	"num_dead_tuples",
}

func GetVacuum(db *sql.DB) ([]Vacuum, error) {
	//tableName := "vacuum_progress"
	tableName := "pg_stat_progress_vacuum"
	query := buildQuery(tableName, VacuumColumns)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	var as []Vacuum
	for rows.Next() {
		var row Vacuum
		err = rows.Scan(
			&row.PID,
			&row.DATID,
			&row.DATNAME,
			&row.RELID,
			&row.PHASE,
			&row.HeapBLKSTotal,
			&row.HeapBLKSScanned,
			&row.HeapBLKSVacuumed,
			&row.IndexVacuumCount,
			&row.MaxDeadTuples,
			&row.NumDeadTuples,
		)
		if err != nil {
			return nil, err
		}
		as = append(as, row)
	}
	return as, nil
}

func (v Vacuum) String() []string {
	return []string{
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
}

func (v Vacuum) Table() string {
	buff := new(bytes.Buffer)
	t := tablewriter.NewWriter(buff)
	t.SetHeader(VacuumColumns)
	t.Append(v.String())
	t.Render()
	return buff.String()
}

func (v Vacuum) Name() string {
	return "pg_stat_progress_vacuum"
}

func (v Vacuum) Progress() float64 {
	return float64(v.HeapBLKSScanned) / float64(v.HeapBLKSTotal)
}

func (v Vacuum) Pid() int {
	return v.PID
}
