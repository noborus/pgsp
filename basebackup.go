package pgsp

import (
	"bytes"
	"context"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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

var BaseBackupColumns = []string{
	"pid",
	"phase",
	"backup_total",
	"backup_streamed",
	"tablespaces_total",
	"tablespaces_streamed",
}

var BaseBackupTableName = "pg_stat_progress_basebackup"

func GetBaseBackup(ctx context.Context, db *sqlx.DB) ([]BaseBackup, error) {
	query := buildQuery(BaseBackupTableName, BaseBackupColumns)
	return selectBaseBackup(ctx, db, query)
}

func selectBaseBackup(ctx context.Context, db *sqlx.DB, query string) ([]BaseBackup, error) {
	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var as []BaseBackup
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

func (v BaseBackup) Table() string {
	buff := new(bytes.Buffer)
	t := tablewriter.NewWriter(buff)
	t.SetHeader(BaseBackupColumns)
	t.Append([]string{
		strconv.Itoa(v.PID),
		v.PHASE,
		strconv.FormatInt(v.BackupTotal, 10),
		strconv.FormatInt(v.BackupStreamed, 10),
		strconv.FormatInt(v.TablespacesTotal, 10),
		strconv.FormatInt(v.TablespacesStreamed, 10),
	})
	t.Render()
	return buff.String()
}

func (v BaseBackup) Vertical() string {
	buff := new(bytes.Buffer)
	vt := vertical.NewWriter(buff)
	vt.SetHeader(BaseBackupColumns)
	vt.Append([]interface{}{
		v.PID,
		v.PHASE,
		v.BackupTotal,
		v.BackupStreamed,
		v.TablespacesTotal,
		v.TablespacesStreamed,
	})
	vt.Render()
	return buff.String()
}

func (v BaseBackup) Progress() float64 {
	if v.BackupTotal != 0 {
		return float64(v.BackupStreamed) / float64(v.BackupTotal)
	}
	return float64(v.TablespacesStreamed) / float64(v.TablespacesTotal)
}
