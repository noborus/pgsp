package pgsp

import (
	"bytes"
	"reflect"
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

func getColumns(s interface{}) []string {
	t := reflect.TypeOf(s)
	var columns []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		j := field.Tag.Get("db")
		columns = append(columns, j)
	}
	return columns
}
