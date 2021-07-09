package tui

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/noborus/pgsp"
)

var (
	UpdateInterval  time.Duration = time.Second / 10
	AfterCompletion time.Duration = 10
	MaxWidth        int           = 80
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
	DB    *sql.DB
	count int
	pgrss []pgrs
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
		for _, pgrs := range m.pgrss {
			if pgrs.p.Width > MaxWidth {
				pgrs.p.Width = MaxWidth
			}
		}
		return m, nil

	case tickMsg:
		err := m.updateProgress(m.DB)
		if err != nil {
			fmt.Printf("update error:%v", err)
		}
		return m, tickCmd()
	}
	return m, nil
}

func (m *Model) updateProgress(db *sql.DB) error {
	cindex, err := pgsp.GetCreateIndex(m.DB)
	if err != nil {
		return err
	}
	for _, v := range cindex {
		m.pgrss = addProgress(m.pgrss, v)
	}

	vacuum, err := pgsp.GetVacuum(m.DB)
	if err != nil {
		return err
	}
	for _, v := range vacuum {
		m.pgrss = addProgress(m.pgrss, v)
	}
	analyze, err := pgsp.GetAnalyze(m.DB)
	if err != nil {
		return err
	}
	for _, v := range analyze {
		m.pgrss = addProgress(m.pgrss, v)
	}

	cluster, err := pgsp.GetCluster(m.DB)
	if err != nil {
		return err
	}
	for _, v := range cluster {
		m.pgrss = addProgress(m.pgrss, v)
	}

	backup, err := pgsp.GetBaseBackup(m.DB)
	if err != nil {
		return err
	}
	for _, v := range backup {
		m.pgrss = addProgress(m.pgrss, v)
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

func addProgress(pgrss []pgrs, v pgsp.PGSProgress) []pgrs {
	for n, pgr := range pgrss {
		if pgr.v.Pid() == v.Pid() {
			pgrss[n].v = v
			pgrss[n].time = time.Now()
			return pgrss
		}
	}
	pg, err := progress.NewModel(
		progress.WithScaledGradient("#FF7CCB", "#FDFF8C"),
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

func (m Model) View() string {
	s := "quit: q, ctrl+c, esc\n"
	for _, pgrs := range m.pgrss {
		if pgrs.p != nil {
			s += pgrs.v.Name() + "\n"
			s += pgrs.v.Table()
			p := pgrs.v.Progress()
			if p > 0 && p <= 1 {
				s += "\n" + pgrs.p.View(p)
				if time.Since(pgrs.time) > time.Second*1 {
					s += " " + time.Since(pgrs.time).Truncate(time.Second).String()
				}
				s += "\n"
			}
		}
	}
	return s
}
