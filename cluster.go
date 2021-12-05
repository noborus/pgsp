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

var ClusterTableName = "pg_stat_progress_cluster"
var ClusterQuery string
var ClusterColumns []string

func GetCluster(ctx context.Context, db *sqlx.DB) ([]PGSProgress, error) {
	if len(ClusterColumns) == 0 {
		ClusterColumns = getColumns(Cluster{})
	}
	if ClusterQuery == "" {
		ClusterQuery = buildQuery(ClusterTableName, ClusterColumns)
		log.Println(ClusterQuery)
	}
	return selectCluster(ctx, db, ClusterQuery)

}

func selectCluster(ctx context.Context, db *sqlx.DB, query string) ([]PGSProgress, error) {
	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var as []PGSProgress
	for rows.Next() {
		var row Cluster
		err = rows.StructScan(&row)
		if err != nil {
			return nil, err
		}
		as = append(as, row)
	}
	return as, rows.Err()
}

func (v Cluster) Name() string {
	return ClusterTableName
}

func (v Cluster) Pid() int {
	return v.PID
}

func (v Cluster) Color() (string, string) {
	return "#5A56E0", "#EE6FF8"
}

func (v Cluster) Table() string {
	value := []string{
		strconv.Itoa(v.PID),
		strconv.Itoa(v.DATID),
		v.DATNAME,
		strconv.Itoa(v.RELID),
		v.Command,
		v.PHASE,
		strconv.FormatInt(v.ClusterIndexRelid, 10),
		strconv.FormatInt(v.HeapTuplesScanned, 10),
		strconv.FormatInt(v.HeapTuplesWritten, 10),
		strconv.FormatInt(v.HeapBlksTotal, 10),
		strconv.FormatInt(v.HeapBlksScanned, 10),
		strconv.FormatInt(v.IndexRebuildCount, 10),
	}
	buff := new(bytes.Buffer)

	t := tablewriter.NewWriter(buff)
	t.SetHeader(ClusterColumns[0:7])
	t.Append(value[0:7])
	t.Render()

	t2 := tablewriter.NewWriter(buff)
	t2.SetHeader(ClusterColumns[7:])
	t2.Append(value[7:])
	t2.Render()

	return buff.String()
}

func (v Cluster) Vertical() string {
	buff := new(bytes.Buffer)
	vt := vertical.NewWriter(buff)
	vt.SetHeader(ClusterColumns)
	vt.AppendStruct(v)
	vt.Render()
	return buff.String()
}

func (v Cluster) Progress() float64 {
	return float64(v.HeapBlksScanned) / float64(v.HeapBlksTotal)
}
