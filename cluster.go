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

// pg_stat_progress_Cluster.
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

var (
	ClusterTableName = "pg_stat_progress_cluster"
	ClusterQuery     string
	ClusterColumns   []string
)

func GetCluster(ctx context.Context, db *sqlx.DB) ([]Progress, error) {
	if len(ClusterColumns) == 0 {
		ClusterColumns = getColumns(Cluster{})
	}
	if ClusterQuery == "" {
		ClusterQuery = buildQuery(ClusterTableName, ClusterColumns)
	}
	return selectCluster(ctx, db, ClusterQuery)
}

func selectCluster(ctx context.Context, db *sqlx.DB, query string) ([]Progress, error) {
	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var as []Progress
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

func (v Cluster) Color() (color.Color, color.Color) {
	return color.RGBA{R: 90, G: 86, B: 224}, color.RGBA{R: 238, G: 111, B: 248}
}

func (v Cluster) Table() string {
	value := str.ToStrStruct(v)
	buff := new(bytes.Buffer)

	t := tablewriter.NewWriter(buff)
	t.Header(ClusterColumns)
	t.Append(value)
	t.Render()

	return buff.String()
}

func (v Cluster) Progress() float64 {
	return float64(v.HeapBlksScanned) / float64(v.HeapBlksTotal)
}
