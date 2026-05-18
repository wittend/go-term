package terminal

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/blacknon/tvxterm"
	"github.com/rivo/tview"
)

// Terminal represents a terminal emulator session.
type Terminal struct {
	View    *tvxterm.View
	Backend *tvxterm.PTYBackend
	LogFile *os.File
	Name    string
}

// Manager handles multiple terminal sessions.
type Manager struct {
	App       *tview.Application
	MainView  *tview.Flex
	Pages     *tview.Pages
	TabBar    *tview.TextView
	Toolbar   *tview.TextView
	StatusBar *tview.TextView
	Sessions  []*Terminal
	Current   int
}

// NewManager creates a new Manager and initializes the layout.
func NewManager(app *tview.Application) *Manager {
	m := &Manager{
		App:       app,
		Pages:     tview.NewPages(),
		TabBar:    tview.NewTextView().SetDynamicColors(true).SetRegions(true).SetWrap(false),
		Toolbar:   tview.NewTextView().SetDynamicColors(true).SetText(" [black:white] F1:New [-] [black:white] F2:Next [-] [black:white] F10:Quit [-]"),
		StatusBar: tview.NewTextView().SetDynamicColors(true),
		Sessions:  []*Terminal{},
		Current:   -1,
	}

	m.MainView = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(m.TabBar, 1, 1, false).
		AddItem(m.Toolbar, 1, 1, false).
		AddItem(m.Pages, 0, 1, true).
		AddItem(m.StatusBar, 1, 1, false)

	return m
}

// NewTerminal creates a new Terminal session.
func NewTerminal(app *tview.Application, shell string, logPath string, name string) (*Terminal, error) {
	view := tvxterm.New(app)
	view.SetBorder(true)

	cmd := exec.Command(shell)
	cmd.Env = append(os.Environ(), "TERM=xterm-256color")

	backend, err := tvxterm.NewPTYBackend(cmd, 80, 24)
	if err != nil {
		return nil, err
	}

	var logFile *os.File
	if logPath != "" {
		logFile, err = os.Create(logPath)
		if err != nil {
			return nil, err
		}
		// We would need to wrap the backend reader to capture data
		// This is a bit complex with tvxterm's internal structure.
		// For now, let's just keep the reference.
	}

	view.Attach(backend)

	return &Terminal{
		View:    view,
		Backend: backend,
		LogFile: logFile,
		Name:    name,
	}, nil
}

// AddSession adds a new terminal session.
func (m *Manager) AddSession(shell string, logPath string) error {
	name := fmt.Sprintf("Term %d", len(m.Sessions)+1)
	term, err := NewTerminal(m.App, shell, logPath, name)
	if err != nil {
		return err
	}

	id := fmt.Sprintf("%d", len(m.Sessions))
	m.Sessions = append(m.Sessions, term)
	m.Pages.AddPage(id, term.View, true, true)
	m.Pages.SwitchToPage(id)
	m.Current = len(m.Sessions) - 1

	term.View.SetTitle(name)

	m.UpdateUI()

	return nil
}

// NextSession switches to the next terminal session.
func (m *Manager) NextSession() {
	if len(m.Sessions) <= 1 {
		return
	}
	m.Current = (m.Current + 1) % len(m.Sessions)
	id := fmt.Sprintf("%d", m.Current)
	m.Pages.SwitchToPage(id)
	m.UpdateUI()
}

// UpdateUI updates the TabBar and StatusBar.
func (m *Manager) UpdateUI() {
	if m.Current < 0 || m.Current >= len(m.Sessions) {
		return
	}

	// Update TabBar
	var tabs string
	for i, session := range m.Sessions {
		if i == m.Current {
			tabs += fmt.Sprintf(" [white:blue] %s [-] ", session.Name)
		} else {
			tabs += fmt.Sprintf(" [black:gray] %s [-] ", session.Name)
		}
	}
	m.TabBar.SetText(tabs)

	// Update StatusBar
	currentSession := m.Sessions[m.Current]
	// Since we can't easily get cursor position from tvxterm, we'll show the name and placeholders
	status := fmt.Sprintf(" [black:white] Terminal: %s [-] [white:blue] Line: - Col: - [-]", currentSession.Name)
	m.StatusBar.SetText(status)
}
