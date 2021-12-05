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
}

var CopyTableName = "pg_stat_progress_copy"
var CopyQuery string
var CopyColumns []string

func GetCopy(ctx context.Context, db *sqlx.DB) ([]Copy, error) {
	if len(CopyColumns) == 0 {
		CopyColumns = getColumns(Copy{})
	}
	if CopyQuery == "" {
		CopyQuery = buildQuery(CopyTableName, CopyColumns)
		log.Println(CopyQuery)
	}
	return selectCopy(ctx, db, CopyQuery)
}

func selectCopy(ctx context.Context, db *sqlx.DB, query string) ([]Copy, error) {
	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var as []Copy
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

func (v Copy) Table() string {
	value := []string{
		strconv.Itoa(v.PID),
		strconv.Itoa(v.DATID),
		v.DATNAME,
		strconv.Itoa(v.RELID),
		v.COMMAND,
		v.CTYPE,
		strconv.FormatInt(v.BYTESProcessed, 10),
		strconv.FormatInt(v.BYTESTotal, 10),
		strconv.FormatInt(v.TUPLESProcessed, 10),
		strconv.FormatInt(v.TUPLESExcluded, 10),
	}

	buff := new(bytes.Buffer)
	t := tablewriter.NewWriter(buff)
	t.SetHeader(CopyColumns[0:7])
	t.Append(value[0:7])
	t.Render()

	t2 := tablewriter.NewWriter(buff)
	t2.SetHeader(CopyColumns[7:])
	t2.Append(value[7:])
	t2.Render()
	return buff.String()
}

func (v Copy) Vertical() string {
	buff := new(bytes.Buffer)
	vt := vertical.NewWriter(buff)
	vt.SetHeader(CopyColumns)
	vt.AppendStruct(v)
	vt.Render()
	return buff.String()
}

func (v Copy) Progress() float64 {
	if v.BYTESTotal == 0 {
		return float64(0.5)
	}
	return float64(v.BYTESProcessed) / float64(v.BYTESTotal)
}
