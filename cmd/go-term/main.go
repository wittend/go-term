package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dave/go-term/pkg/terminal"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	// Command line arguments
	shell := flag.String("shell", "bash", "Shell to use")
	title := flag.String("title", "Go-Term", "Terminal title")
	log := flag.String("log", "", "File to capture session to")
	flag.Parse()

	_ = title // title is currently not used for window title since we use tabs
	app := tview.NewApplication()
	manager := terminal.NewManager(app)

	err := manager.AddSession(*shell, *log)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlT, tcell.KeyF1:
			manager.AddSession(*shell, "")
			return nil
		case tcell.KeyCtrlN, tcell.KeyF2:
			manager.NextSession()
			return nil
		case tcell.KeyCtrlQ, tcell.KeyF10:
			app.Stop()
			return nil
		}
		return event
	})

	if err := app.SetRoot(manager.MainView, true).Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
