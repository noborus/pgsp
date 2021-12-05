package pgsp

import (
	"bytes"
	"context"
	"reflect"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type SPTaget string

const (
	SPAnalyze     = "Analyze"
	SPCreateIndex = "CreateIndex"
	SPVacuum      = "Vacuum"
	SPCluster     = "Cluster"
	SPBaseBackup  = "BaseBackup"
	SPCopy        = "Copy"
)

type SPTable struct {
	Enable bool
	Get    func(ctx context.Context, db *sqlx.DB) ([]Progress, error)
}

type StatProgress map[SPTaget]*SPTable

type Progress interface {
	Name() string
	Pid() int
	Color() (string, string)
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

func NewMonitor() StatProgress {
	return map[SPTaget]*SPTable{
		SPAnalyze: {
			Get: GetAnalyze,
		},
		SPCreateIndex: {
			Get: GetCreateIndex,
		},
		SPVacuum: {
			Get: GetVacuum,
		},
		SPCluster: {
			Get: GetCluster,
		},
		SPBaseBackup: {
			Get: GetBaseBackup,
		},
		SPCopy: {
			Get: GetCopy,
		},
	}

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

func Targets(sp StatProgress, target []string) StatProgress {
	if len(target) != 0 {
		enableF := false
		for _, t := range target {
			if v, ok := sp[SPTaget(t)]; ok {
				enableF = true
				v.Enable = true
			}
		}
		// Return if there is even one target.
		if enableF {
			return sp
		}
	}

	// All targets.
	for _, v := range sp {
		v.Enable = true
	}
	return sp
}

func TargetString(sp StatProgress) string {
	var ms []string
	for n, v := range sp {
		if v.Enable {
			ms = append(ms, string(n))
		}
	}
	sort.Strings(ms)
	return strings.Join(ms, " ")
}
