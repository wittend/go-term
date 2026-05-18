package terminal

import (
	"os"
	"testing"

	"github.com/rivo/tview"
)

func TestNewManager(t *testing.T) {
	app := tview.NewApplication()
	m := NewManager(app)
	if m == nil {
		t.Fatal("Expected NewManager to return a non-nil value")
	}
	if m.App != app {
		t.Errorf("Expected App to be set to %v, got %v", app, m.App)
	}
	if m.Pages == nil {
		t.Error("Expected Pages to be initialized")
	}
	if m.MainView == nil {
		t.Error("Expected MainView to be initialized")
	}
	if m.TabBar == nil {
		t.Error("Expected TabBar to be initialized")
	}
	if m.StatusBar == nil {
		t.Error("Expected StatusBar to be initialized")
	}
	if len(m.Sessions) != 0 {
		t.Errorf("Expected 0 sessions initially, got %d", len(m.Sessions))
	}
	if m.Current != -1 {
		t.Errorf("Expected Current to be -1 initially, got %d", m.Current)
	}
}

func TestAddSession(t *testing.T) {
	app := tview.NewApplication()
	m := NewManager(app)

	// Use 'echo' as a lightweight "shell" for testing
	err := m.AddSession("echo", "")
	if err != nil {
		t.Fatalf("Failed to add session: %v", err)
	}

	if len(m.Sessions) != 1 {
		t.Errorf("Expected 1 session, got %d", len(m.Sessions))
	}
	if m.Current != 0 {
		t.Errorf("Expected Current to be 0, got %d", m.Current)
	}

	// Add another session
	err = m.AddSession("echo", "")
	if err != nil {
		t.Fatalf("Failed to add second session: %v", err)
	}

	if len(m.Sessions) != 2 {
		t.Errorf("Expected 2 sessions, got %d", len(m.Sessions))
	}
	if m.Current != 1 {
		t.Errorf("Expected Current to be 1, got %d", m.Current)
	}
}

func TestNextSession(t *testing.T) {
	app := tview.NewApplication()
	m := NewManager(app)

	// No sessions yet
	m.NextSession()
	if m.Current != -1 {
		t.Errorf("Expected Current to remain -1 with no sessions, got %d", m.Current)
	}

	// Add two sessions
	m.AddSession("echo", "")
	m.AddSession("echo", "")

	if m.Current != 1 {
		t.Errorf("Expected Current to be 1 after adding two sessions, got %d", m.Current)
	}

	m.NextSession()
	if m.Current != 0 {
		t.Errorf("Expected Current to wrap around to 0, got %d", m.Current)
	}

	m.NextSession()
	if m.Current != 1 {
		t.Errorf("Expected Current to be 1 after another NextSession, got %d", m.Current)
	}
}

func TestNewTerminalLogging(t *testing.T) {
	app := tview.NewApplication()
	logPath := "test_session.log"
	defer os.Remove(logPath)

	term, err := NewTerminal(app, "echo", logPath, "TestTerm")
	if err != nil {
		t.Fatalf("Failed to create terminal with logging: %v", err)
	}
	defer term.LogFile.Close()

	if term.LogFile == nil {
		t.Error("Expected LogFile to be non-nil")
	}

	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Errorf("Expected log file %s to exist", logPath)
	}
}
