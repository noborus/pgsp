package tui

import (
	"context"
	"fmt"
	"log"
	"time"

	"charm.land/bubbles/v2/progress"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/noborus/pgsp"
)

var (
	UpdateInterval  time.Duration
	AfterCompletion time.Duration
	RightMargin     int = 1
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
	spinC      int
	pgrss      []pgrs
	width      int
	height     int
	monitor    *pgsp.Pgsp
	status     string
	fullscreen bool
	ready      bool
	viewport   viewport.Model
	content    string
}

var spin []string = []string{"|", "/", "-", "\\"}

var (
	titleStyle = func() lipgloss.Style {
		return lipgloss.NewStyle().Padding(0, 1)
	}()

	headerStyle = func() lipgloss.Style {
		return lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4"))
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)

func (m Model) Init() tea.Cmd {
	return tickCmd()
}

type Option func(*Model) error

func NewModel(monitor *pgsp.Pgsp, fullscreen bool) Model {
	model := Model{
		monitor:    monitor,
		fullscreen: fullscreen,
	}
	DebugLogf("NewModel: %v, %p", model, &model)
	return model
}

func NewProgram(m Model) *tea.Program {
	p := tea.NewProgram(m)
	// if fullScreen {
	// 	p.EnterAltScreen()
	// }
	return p
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	ctx := context.TODO()
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		k := msg.Key().Code
		mod := msg.Key().Mod
		if (k == 'c' && mod == tea.ModCtrl) || k == tea.KeyEsc || k == 'q' {
			return m, tea.Quit
		} else if k == tea.KeyDown && m.ready {
			m.viewport.ScrollDown(4)
		} else if k == tea.KeyUp && m.ready {
			m.viewport.ScrollUp(4)
		} else if k == tea.KeyLeft && m.ready {
			m.viewport.ScrollLeft(4)
		} else if k == tea.KeyRight && m.ready {
			m.viewport.ScrollRight(4)
		}

	case tea.WindowSizeMsg:
		DebugLogf("WindowSizeMsg: width=%d, height=%d", msg.Width, msg.Height)
		m.height = msg.Height
		m.width = msg.Width
		headerHeight := lipgloss.Height(m.headerView())
		verticalMarginHeight := headerHeight

		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(viewport.WithWidth(msg.Width), viewport.WithHeight(msg.Height-verticalMarginHeight))
			m.viewport.YPosition = headerHeight
			m.viewport.MouseWheelEnabled = true
			m.ready = true
		} else {
			m.viewport.SetWidth(msg.Width)
			m.viewport.SetHeight(msg.Height - verticalMarginHeight)
		}

		for _, pgrs := range m.pgrss {
			pgrs.p.SetWidth(m.viewport.Width() - RightMargin)
		}
	case tickMsg:
		m.spinC++
		if m.spinC > len(spin)-1 {
			m.spinC = 0
		}
		err := m.updateProgress(ctx)

		if err != nil {
			fmt.Printf("update error:%v", err)
		}

		m.content = m.progressView()

		cmds = append(cmds, tickCmd())
	}

	if m.ready {
		// Handle keyboard and mouse events in the viewport
		DebugLogf("Update viewport with msg: %T and content length: %d, %p", msg, len(m.content), &m)
		m.viewport.SetContent(m.content)
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	for _, pgrs := range m.pgrss {
		p, cmd := pgrs.p.Update(msg)
		pgrs.p = &p
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) headerView() string {
	title := titleStyle.Render(m.status + "quit: q, ctrl+c, esc; scroll: arrow keys\n")
	return title
}

func (m Model) View() tea.View {
	var v tea.View
	v.AltScreen = m.fullscreen
	v.MouseMode = tea.MouseModeCellMotion
	if !m.ready {
		v.SetContent("\n  Initializing...")
	} else {
		v.SetContent(fmt.Sprintf("%s\n%s", m.headerView(), m.viewport.View()))
	}
	return v
}

func (m *Model) progressView() string {
	var s string
	num := len(m.pgrss)
	if num == 0 {
		s = spin[m.spinC] + " " + s
	}

	for _, pgrs := range m.pgrss {
		if pgrs.p == nil {
			continue
		}
		s += headerStyle.Render(pgrs.v.Name()) + "\n"
		s += pgrs.v.Table()
		p := pgrs.v.Progress()
		if p > 0 && p <= 1 {
			if time.Since(pgrs.time) > time.Second*1 {
				// Deleted records are considered 100%.
				s += "\n" + pgrs.p.ViewAs(float64(1))
				s += " " + time.Since(pgrs.time).Truncate(time.Second).String()
			} else {
				s += "\n" + pgrs.p.ViewAs(p)
			}
			s += "\n"
		}
	}
	DebugLogf("Set content with length: %d, %p", len(s), &m)
	return s
}

func (m *Model) updateProgress(ctx context.Context) error {
	m.status = fmt.Sprintf("Monitor: %s\n", m.monitor.TargetString())

	for _, table := range m.monitor.StatProgress {
		result, err := table.Get(ctx, m.monitor.DB)
		if err != nil {
			table.Enable = false
			DebugLog(err)
			m.status += err.Error() + "\n"
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

	pg := progress.New(
		progress.WithScaled(true),
		progress.WithColors(v.Color()),
		progress.WithWidth(m.width-RightMargin),
	)
	pgrs := pgrs{
		time: time.Now(),
		v:    v,
		p:    &pg,
	}
	pgrss = append(pgrss, pgrs)
	return pgrss
}
