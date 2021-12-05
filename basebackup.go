package pgsp

import (
	"bytes"
	"context"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/noborus/pgsp/str"
	"github.com/noborus/pgsp/vertical"
	"github.com/olekukonko/tablewriter"
)

// pg_stat_progress_basebackup
type BaseBackup struct {
	PID                 int    `db:"pid"`
	PHASE               string `db:"phase"`
	BackupTotal         int64  `db:"backup_total"`
	BackupStreamed      int64  `db:"backup_streamed"`
	TablespacesTotal    int64  `db:"tablespaces_total"`
	TablespacesStreamed int64  `db:"tablespaces_streamed"`
}

var BaseBackupTableName = "pg_stat_progress_basebackup"
var BaseBackupQuery string
var BaseBackupColumns []string

func GetBaseBackup(ctx context.Context, db *sqlx.DB) ([]PGSProgress, error) {
	if len(BaseBackupColumns) == 0 {
		BaseBackupColumns = getColumns(BaseBackup{})
	}
	if BaseBackupQuery == "" {
		BaseBackupQuery = buildQuery(BaseBackupTableName, BaseBackupColumns)
	}
	return selectBaseBackup(ctx, db, BaseBackupQuery)
}

func selectBaseBackup(ctx context.Context, db *sqlx.DB, query string) ([]PGSProgress, error) {
	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var as []PGSProgress
	for rows.Next() {
		var row BaseBackup
		err = rows.StructScan(&row)
		if err != nil {
			return nil, err
		}
		as = append(as, row)
	}
	return as, rows.Err()
}

func (v BaseBackup) Name() string {
	return BaseBackupTableName
}

func (v BaseBackup) Pid() int {
	return v.PID
}

func (v BaseBackup) Color() (string, string) {
	return "#FDFF8C", "#FF7CCB"
}

func (v BaseBackup) Table() string {
	buff := new(bytes.Buffer)
	t := tablewriter.NewWriter(buff)
	t.SetHeader(BaseBackupColumns)
	t.Append(str.ToStrStruct(v))
	t.Render()
	return buff.String()
}

func (v BaseBackup) Vertical() string {
	buff := new(bytes.Buffer)
	vt := vertical.NewWriter(buff)
	vt.SetHeader(BaseBackupColumns)
	vt.AppendStruct(v)
	vt.Render()
	return buff.String()
}

func (v BaseBackup) Progress() float64 {
	if v.BackupTotal != 0 {
		return float64(v.BackupStreamed) / float64(v.BackupTotal)
	}
	return float64(v.TablespacesStreamed) / float64(v.TablespacesTotal)
}
