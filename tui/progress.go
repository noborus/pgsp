package tui

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jmoiron/sqlx"
	"github.com/noborus/pgsp"
)

var (
	UpdateInterval    time.Duration
	AfterCompletion   time.Duration
	RightMargin       int = 10
	MinimumTableWidth int = 120
	MaxVerticalRows   int = 15
)

var Debug = false

type tickMsg time.Time

func DebugLogf(format string, v ...interface{}) {
	if Debug {
		log.Printf(format, v...)
	}
}

func DebugLog(v ...interface{}) {
	if Debug {
		log.Print(v...)
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(UpdateInterval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type pgrs struct {
	time time.Time
	v    pgsp.Progress
	p    *progress.Model
}

type Model struct {
	db      *sqlx.DB
	spinC   int
	pgrss   []pgrs
	width   int
	height  int
	monitor pgsp.StatProgress
	status  string
}

var spin []string = []string{"|", "/", "-", "\\"}

func (m Model) Init() tea.Cmd {
	return tickCmd()
}

type Option func(*Model) error

func NewModel(db *sqlx.DB, monitor pgsp.StatProgress) Model {
	model := Model{
		db:      db,
		monitor: monitor,
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
		err := m.updateProgress(ctx)
		if err != nil {
			fmt.Printf("update error:%v", err)
		}
		return m, tickCmd()
	}
	return m, nil
}

func (m Model) View() string {
	s := m.status
	s += "quit: q, ctrl+c, esc\n"
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
			if m.width >= MinimumTableWidth {
				s += pgrs.v.Table()
			} else {
				if num*MaxVerticalRows < m.height {
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

func (m *Model) updateProgress(ctx context.Context) error {
	m.status = pgsp.TargetString(m.monitor)

	for _, table := range m.monitor {
		result, err := table.Get(ctx, m.db)
		if err != nil {
			DebugLog(err)
			m.status += err.Error()
		}
		for _, v := range result {
			m.pgrss = m.addProgress(m.pgrss, v)
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

func (m Model) addProgress(pgrss []pgrs, v pgsp.Progress) []pgrs {
	for n, pgr := range pgrss {
		if pgr.v.Name() == v.Name() && pgr.v.Pid() == v.Pid() {
			pgrss[n].v = v
			pgrss[n].time = time.Now()
			return pgrss
		}
	}

	pg, err := progress.NewModel(
		progress.WithScaledGradient(v.Color()),
		progress.WithWidth(m.width-RightMargin),
	)
	if err != nil {
		DebugLog(err)
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
