package pgsp

import (
	"bytes"
	"database/sql"
	"strings"

	_ "github.com/lib/pq"
)

type PGSProgress interface {
	Pid() int
	Name() string
	Table() string
	Progress() float64
}

func Connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func DisConnect(db *sql.DB) error {
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
