package pgsp

import (
	"bytes"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PGSProgress interface {
	Name() string
	Pid() int
	Table() string
	Vertical() string
	Progress() float64
}

func Connect(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func DisConnect(db *sqlx.DB) error {
	return db.Close()
}

func buildQuery(tableName string, columns []string) string {
	buff := new(bytes.Buffer)
	buff.WriteString("SELECT ")
	buff.WriteString(strings.Join(columns, ", "))
	buff.WriteString(" FROM ")
	buff.WriteString(tableName)
	return buff.String()
}
