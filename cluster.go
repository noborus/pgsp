package pgsp

import (
	"bytes"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/olekukonko/tablewriter"
)

// pg_stat_progress_Cluster
type Cluster struct {
	PID               int    `db:"pid"`
	DATID             int    `db:"datid"`
	DATNAME           string `db:"datname"`
	RELID             int    `db:"relid"`
	Command           string `db:"command"`
	PHASE             string `db:"phase"`
	ClusterIndexRelid int64  `db:"cluster_index_relid"`
	HeapTuplesScanned int64  `db:"heap_tuples_scanned"`
	HeapTuplesWritten int64  `db:"heap_tuples_written"`
	HeapBlksTotal     int64  `db:"heap_blks_total"`
	HeapBlksScanned   int64  `db:"heap_blks_scanned"`
	IndexRebuildCount int64  `db:"index_rebuild_count"`
}

var ClusterColumns = []string{
	"pid",
	"datid",
	"datname",
	"relid",
	"command",
	"phase",
	"cluster_index_relid",
	"heap_tuples_scanned",
	"heap_tuples_written",
	"heap_blks_total",
	"heap_blks_scanned",
	"index_rebuild_count",
}

func GetCluster(db *sql.DB) ([]Cluster, error) {
	// tableName := "cluster_progress"
	tableName := "pg_stat_progress_cluster"
	query := buildQuery(tableName, ClusterColumns)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	var as []Cluster
	for rows.Next() {
		var row Cluster
		err = rows.Scan(&row.PID, &row.DATID, &row.DATNAME, &row.RELID, &row.Command, &row.PHASE, &row.ClusterIndexRelid, &row.HeapBlksScanned, &row.HeapTuplesWritten, &row.HeapBlksTotal, &row.HeapBlksScanned, &row.IndexRebuildCount)
		if err != nil {
			return nil, err
		}
		as = append(as, row)
	}
	return as, nil
}

func (v Cluster) String() []string {
	pid := fmt.Sprintf("%v", v.PID)
	datid := fmt.Sprintf("%v", v.DATID)
	total := fmt.Sprintf("%v", v.HeapBlksTotal)
	scanned := fmt.Sprintf("%v", v.HeapBlksScanned)
	return []string{pid, datid, v.DATNAME, v.PHASE, total, scanned}
}

func (v Cluster) Table() string {
	buff := new(bytes.Buffer)
	t := tablewriter.NewWriter(buff)
	t.SetHeader(ClusterColumns)
	t.Append(v.String())
	t.Render()
	return buff.String()
}

func (v Cluster) Name() string {
	return "pg_stat_progress_cluster"
}

func (v Cluster) Progress() float64 {
	return float64(v.HeapBlksScanned) / float64(v.HeapBlksTotal)
}
