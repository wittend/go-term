# Go-Term

A terminal emulator written in Go using `tview` and `tvxterm`.

## Features

- [x] Multi-tab support (`Ctrl+T` to open, `Ctrl+N` to switch)
- [x] Bordered windows
- [x] Customizable shell and title via CLI flags
- [x] Session logging support
- [x] Responsive resizing
- [x] Cross-platform (Linux, macOS, BSD)

## Installation

```bash
go build ./cmd/go-term
```

## Usage

```bash
./go-term -shell /bin/zsh -title "My Custom Term" -log session.log
```

## Keyboard Shortcuts

- `Ctrl+T`: Open new terminal tab
- `Ctrl+N`: Switch to next tab
- `Ctrl+Q`: Quit application

## Architecture

The project is structured modularly:
- `cmd/go-term`: Entry point and UI orchestration.
- `pkg/terminal`: Terminal session and tab management.
- `pkg/config`: (Planned) Configuration and theme management.
