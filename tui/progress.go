package tui

import (
	"database/sql"
	"log"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/noborus/pgsp"
)

type model struct {
	percent    float64
	progresses []*progress.Model
}

type tickMsg time.Time

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second/10, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

const (
	padding  = 2
	maxWidth = 80
)

type Model struct {
	DB       *sql.DB
	count    int
	v        []pgsp.PGSProgress
	progress []*progress.Model
}

func (m Model) Init() tea.Cmd {
	return tickCmd()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		default:
			return m, nil
		}

	case tea.WindowSizeMsg:
		return m, nil

	case tickMsg:
		m.statProgress(m.DB)
		return m, tickCmd()
	}
	return m, nil
}

func addProgress(m *Model, v pgsp.PGSProgress) {
	m.v = append(m.v, v)

	pg, err := progress.NewModel(
	//		progress.WithScaledGradient("#FF7CCB", "#FDFF8C"),
	)
	if err != nil {
		log.Println(err)
		return
	}
	m.progress = append(m.progress, pg)
}

func (m *Model) statProgress(db *sql.DB) {
	m.v = []pgsp.PGSProgress{}
	m.progress = []*progress.Model{}

	cindex, err := pgsp.GetCreateIndex(m.DB)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range cindex {
		addProgress(m, v)
	}

	vacuum, err := pgsp.GetVacuum(m.DB)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range vacuum {
		addProgress(m, v)
	}
	analyze, err := pgsp.GetAnalyze(m.DB)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range analyze {
		addProgress(m, v)
	}

	cluster, err := pgsp.GetCluster(m.DB)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range cluster {
		addProgress(m, v)
	}
	backup, err := pgsp.GetBaseBackup(m.DB)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range backup {
		addProgress(m, v)
	}
}

func (m *Model) statProgressHisotry(db *sql.DB) {
	cindex, err := pgsp.GetCreateIndex(m.DB)
	if err != nil {
		log.Fatal(err)
	}
	vacuum, err := pgsp.GetVacuum(m.DB)
	if err != nil {
		log.Fatal(err)
	}
	analyze, err := pgsp.GetAnalyze(m.DB)
	if err != nil {
		log.Fatal(err)
	}

	m.count++

	if len(m.v) == 0 {
		m.v = make([]pgsp.PGSProgress, 3)
	}
	if len(cindex) > m.count+1 {
		m.v[0] = cindex[m.count]
	}
	if len(vacuum) > m.count+1 {
		m.v[1] = vacuum[m.count]
	}
	m.v[2] = analyze[0]
}

func (m Model) View() string {
	s := "quit: q, ctrl+c, esc"
	for n, v := range m.v {
		if v != nil {
			s += v.Name() + "\n"
			s += v.Table()
			p := v.Progress()
			if p > 0 {
				s += "\n" + m.progress[n].View(v.Progress()) + "\n"
			}
		}
	}
	return s
}
