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
	SPAnalyze     SPTaget = "Analyze"
	SPCreateIndex SPTaget = "CreateIndex"
	SPVacuum      SPTaget = "Vacuum"
	SPCluster     SPTaget = "Cluster"
	SPBaseBackup  SPTaget = "BaseBackup"
	SPCopy        SPTaget = "Copy"
)

type SPTable struct {
	Enable bool
	Get    func(ctx context.Context, db *sqlx.DB) ([]Progress, error)
}

type StatProgress map[SPTaget]*SPTable

type Pgsp struct {
	DB           *sqlx.DB
	StatProgress StatProgress
}

type Progress interface {
	Name() string
	Pid() int
	Color() (string, string)
	Table() string
	Vertical() string
	Progress() float64
}

func New(dsn string) (*Pgsp, error) {
	db, err := Connect(dsn)
	if err != nil {
		return nil, err
	}
	monitor := NewMonitor()
	return &Pgsp{
		DB:           db,
		StatProgress: monitor,
	}, nil
}

func NewMonitor() StatProgress {
	return StatProgress{
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

func (p *Pgsp) DisConnect() error {
	return p.DB.Close()
}

func (p *Pgsp) Targets(target []string) {
	if len(target) != 0 {
		enableF := false
		for _, t := range target {
			if v, ok := p.StatProgress[SPTaget(t)]; ok {
				enableF = true
				v.Enable = true
			}
		}
		// Return if there is even one target.
		if enableF {
			return
		}
	}

	// All targets.
	for _, v := range p.StatProgress {
		v.Enable = true
	}
}

func (p *Pgsp) TargetString() string {
	var ms []string
	for n, v := range p.StatProgress {
		if v.Enable {
			ms = append(ms, string(n))
		}
	}
	sort.Strings(ms)
	return strings.Join(ms, " ")
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
