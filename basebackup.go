package pgsp

import (
	"bytes"
	"database/sql"
	"strconv"

	_ "github.com/lib/pq"
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

func GetBaseBackup(db *sql.DB) ([]BaseBackup, error) {
	tableName := "pg_stat_progress_basebackup"
	query := buildQuery(tableName, BaseBackupColumns)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var as []BaseBackup
	for rows.Next() {
		var row BaseBackup
		err = rows.Scan(
			&row.PID,
			&row.PHASE,
			&row.BackupTotal,
			&row.BackupStreamed,
			&row.TablespacesTotal,
			&row.TablespacesStreamed,
		)
		if err != nil {
			return nil, err
		}
		as = append(as, row)
	}
	return as, nil
}

func (v BaseBackup) String() []string {
	return []string{
		strconv.Itoa(v.PID),
		v.PHASE,
		strconv.FormatInt(v.BackupTotal, 10),
		strconv.FormatInt(v.BackupStreamed, 10),
		strconv.FormatInt(v.TablespacesTotal, 10),
		strconv.FormatInt(v.TablespacesStreamed, 10),
	}
}

func (v BaseBackup) Table() string {
	buff := new(bytes.Buffer)
	t := tablewriter.NewWriter(buff)
	t.SetHeader(BaseBackupColumns)
	t.Append(v.String())
	t.Render()
	return buff.String()
}

func (v BaseBackup) Name() string {
	return "pg_stat_progress_basebackup"
}

func (v BaseBackup) Progress() float64 {
	if v.BackupTotal != 0 {
		return float64(v.BackupStreamed) / float64(v.BackupTotal)
	}
	return float64(v.TablespacesStreamed) / float64(v.TablespacesTotal)
}

func (v BaseBackup) Pid() int {
	return v.PID
}
