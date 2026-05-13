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

	app := tview.NewApplication()
	manager := terminal.NewManager(app)

	err := manager.AddSession(*shell, *log)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlT:
			manager.AddSession(*shell, "")
			return nil
		case tcell.KeyCtrlN:
			manager.NextSession()
			return nil
		case tcell.KeyCtrlQ:
			app.Stop()
			return nil
		case tcell.KeyCtrlL:
			// Example of how we could trigger a "Connect to SSH" dialog
			return event
		}
		return event
	})

	manager.Pages.SetTitle(*title)

	if err := app.SetRoot(manager.Pages, true).Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
