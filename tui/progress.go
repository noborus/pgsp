package tui

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/noborus/pgsp"
)

var (
	UpdateInterval  time.Duration
	AfterCompletion time.Duration
	RightMargin     int = 10
)

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(UpdateInterval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type pgrs struct {
	time time.Time
	v    pgsp.PGSProgress
	p    *progress.Model
}

type Model struct {
	DB     *sql.DB
	spinC  int
	pgrss  []pgrs
	width  int
	height int

	CreateIndex bool
	Vacuum      bool
	Analyze     bool
	Cluster     bool
	BaseBackup  bool
}

var colorTables map[string][]string = map[string][]string{
	"pg_stat_progress_analyze":      {"#FF7CCB", "#FDFF8C"},
	"pg_stat_progress_basebackup":   {"#FDFF8C", "#FF7CCB"},
	"pg_stat_progress_cluster":      {"#5A56E0", "#EE6FF8"},
	"pg_stat_progress_create_index": {"#EE6FF8", "#5A56E0"},
	"pg_stat_progress_vacuum":       {"#5A56E0", "#FF7CCB"},
}

var spin []string = []string{"|", "/", "-", "\\"}

func (m Model) Init() tea.Cmd {
	return tickCmd()
}

type Option func(*Model) error

func NewModel(db *sql.DB) Model {
	model := Model{
		DB:          db,
		CreateIndex: true,
		Vacuum:      true,
		Analyze:     true,
		Cluster:     true,
		BaseBackup:  true,
	}
	return model
}

func NewProgram(m Model, fullScreen bool) *tea.Program {
	p := tea.NewProgram(m)
	if fullScreen {
		p.EnterAltScreen()
	}
	return p
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	ctx := context.TODO()
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		default:
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		for _, pgrs := range m.pgrss {
			pgrs.p.Width = m.width - RightMargin
		}
		return m, nil

	case tickMsg:
		m.spinC++
		if m.spinC > len(spin)-1 {
			m.spinC = 0
		}
		err := m.updateProgress(ctx, m.DB)
		if err != nil {
			fmt.Printf("update error:%v", err)
		}
		return m, tickCmd()
	}
	return m, nil
}

func (m Model) View() string {
	s := "quit: q, ctrl+c, esc\n"
	num := len(m.pgrss)
	if num == 0 {
		s = spin[m.spinC] + " " + s
		return s
	}
	var style = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4"))
	for _, pgrs := range m.pgrss {
		if pgrs.p != nil {
			s += style.Render(pgrs.v.Name()) + "\n"
			if m.width >= 120 {
				s += pgrs.v.Table()
			} else {
				if num*15 < m.height {
					s += pgrs.v.Vertical()
				}
			}
			p := pgrs.v.Progress()
			if p > 0 && p <= 1 {
				if time.Since(pgrs.time) > time.Second*1 {
					// Deleted records are considered 100%.
					s += "\n" + pgrs.p.View(1)
					s += " " + time.Since(pgrs.time).Truncate(time.Second).String()
				} else {
					s += "\n" + pgrs.p.View(p)
				}
				s += "\n"
			}
		}
	}
	return s
}

func (m *Model) addCreateIndex(ctx context.Context) error {
	cindex, err := pgsp.GetCreateIndex(ctx, m.DB)
	if err != nil {
		return err
	}
	for _, v := range cindex {
		m.pgrss = m.addProgress(m.pgrss, v)
	}
	return nil
}

func (m *Model) addVacuum(ctx context.Context) error {
	vacuum, err := pgsp.GetVacuum(ctx, m.DB)
	if err != nil {
		return err
	}
	for _, v := range vacuum {
		m.pgrss = m.addProgress(m.pgrss, v)
	}
	return nil
}

func (m *Model) addAnalyze(ctx context.Context) error {
	analyze, err := pgsp.GetAnalyze(ctx, m.DB)
	if err != nil {
		return err
	}
	for _, v := range analyze {
		m.pgrss = m.addProgress(m.pgrss, v)
	}
	return nil
}

func (m *Model) addCluster(ctx context.Context) error {
	cluster, err := pgsp.GetCluster(ctx, m.DB)
	if err != nil {
		return err
	}
	for _, v := range cluster {
		m.pgrss = m.addProgress(m.pgrss, v)
	}
	return nil
}

func (m *Model) addBaseBackup(ctx context.Context) error {
	backup, err := pgsp.GetBaseBackup(ctx, m.DB)
	if err != nil {
		return err
	}
	for _, v := range backup {
		m.pgrss = m.addProgress(m.pgrss, v)
	}
	return nil
}

func (m *Model) updateProgress(ctx context.Context, db *sql.DB) error {
	if m.CreateIndex {
		if err := m.addCreateIndex(ctx); err != nil {
			log.Println(err)
			m.CreateIndex = false
		}
	}

	if m.Vacuum {
		if err := m.addVacuum(ctx); err != nil {
			log.Println(err)
			m.Vacuum = false
		}
	}

	if m.Analyze {
		if err := m.addAnalyze(ctx); err != nil {
			log.Println(err)
			m.Analyze = false
		}
	}

	if m.Cluster {
		if err := m.addCluster(ctx); err != nil {
			log.Println(err)
			m.Cluster = false
		}
	}

	if m.BaseBackup {
		if err := m.addBaseBackup(ctx); err != nil {
			log.Println(err)
			m.BaseBackup = false
		}
	}

	pgrss := make([]pgrs, 0, len(m.pgrss))
	for _, pgrs := range m.pgrss {
		if time.Since(pgrs.time) < time.Second*AfterCompletion {
			pgrss = append(pgrss, pgrs)
		}
	}
	m.pgrss = pgrss
	return nil
}

func (m Model) addProgress(pgrss []pgrs, v pgsp.PGSProgress) []pgrs {
	for n, pgr := range pgrss {
		if pgr.v.Name() == v.Name() && pgr.v.Pid() == v.Pid() {
			pgrss[n].v = v
			pgrss[n].time = time.Now()
			return pgrss
		}
	}
	c := colorTables[v.Name()]
	pg, err := progress.NewModel(
		progress.WithScaledGradient(c[0], c[1]),
		progress.WithWidth(m.width-RightMargin),
	)
	if err != nil {
		log.Println(err)
		return nil
	}
	pgrs := pgrs{
		time: time.Now(),
		v:    v,
		p:    pg,
	}
	pgrss = append(pgrss, pgrs)
	return pgrss
}
