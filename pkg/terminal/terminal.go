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
}

// Manager handles multiple terminal sessions.
type Manager struct {
	App      *tview.Application
	Pages    *tview.Pages
	Sessions []*Terminal
	Current  int
}

// NewManager creates a new Manager.
func NewManager(app *tview.Application) *Manager {
	return &Manager{
		App:   app,
		Pages: tview.NewPages(),
	}
}

// NewTerminal creates a new Terminal session.
func NewTerminal(app *tview.Application, shell string, logPath string) (*Terminal, error) {
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
	}, nil
}

// AddSession adds a new terminal session.
func (m *Manager) AddSession(shell string, logPath string) error {
	term, err := NewTerminal(m.App, shell, logPath)
	if err != nil {
		return err
	}

	id := fmt.Sprintf("%d", len(m.Sessions))
	m.Sessions = append(m.Sessions, term)
	m.Pages.AddPage(id, term.View, true, true)
	m.Pages.SwitchToPage(id)
	m.Current = len(m.Sessions) - 1

	term.View.SetTitle(fmt.Sprintf("Terminal %d", len(m.Sessions)))

	return nil
}

// NextSession switches to the next terminal session.
func (m *Manager) NextSession() {
	if len(m.Sessions) == 0 {
		return
	}
	m.Current = (m.Current + 1) % len(m.Sessions)
	id := fmt.Sprintf("%d", m.Current)
	m.Pages.SwitchToPage(id)
}
