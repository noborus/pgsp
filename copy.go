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

// pg_stat_progress_copy
type Copy struct {
	PID             int    `db:"pid"`
	DATID           int    `db:"datid"`
	DATNAME         string `db:"datname"`
	RELID           int    `db:"relid"`
	COMMAND         string `db:"command"`
	CTYPE           string `db:"type"`
	BYTESProcessed  int64  `db:"bytes_processed"`
	BYTESTotal      int64  `db:"bytes_total"`
	TUPLESProcessed int64  `db:"tuples_processed"`
	TUPLESExcluded  int64  `db:"tuples_excluded"`
	TUPLESSkipped   int64  `db:"tuples_skipped"`
}

var (
	CopyTableName = "pg_stat_progress_copy"
	CopyQuery     string
	CopyColumns   []string
)

func GetCopy(ctx context.Context, db *sqlx.DB) ([]Progress, error) {
	if len(CopyColumns) == 0 {
		CopyColumns = getColumns(Copy{})
	}
	if CopyQuery == "" {
		CopyQuery = buildQuery(CopyTableName, CopyColumns)
	}
	return selectCopy(ctx, db, CopyQuery)
}

func selectCopy(ctx context.Context, db *sqlx.DB, query string) ([]Progress, error) {
	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var as []Progress
	for rows.Next() {
		var row Copy
		err = rows.StructScan(&row)
		if err != nil {
			return nil, err
		}
		as = append(as, row)
	}
	return as, nil
}

func (v Copy) Name() string {
	return CopyTableName
}

func (v Copy) Pid() int {
	return v.PID
}

func (v Copy) Color() (color.Color, color.Color) {
	return color.RGBA{R: 90, G: 246, B: 255}, color.RGBA{R: 124, G: 255, B: 203}
}

func (v Copy) Table() string {
	value := str.ToStrStruct(v)
	buff := new(bytes.Buffer)
	t := tablewriter.NewWriter(buff)
	t.Header(CopyColumns)
	t.Append(value)
	t.Render()
	return buff.String()
}

func (v Copy) Progress() float64 {
	if v.BYTESTotal == 0 {
		return float64(0.5)
	}
	return float64(v.BYTESProcessed) / float64(v.BYTESTotal)
}
