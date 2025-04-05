package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type app struct{}

func (a app) Init() tea.Cmd {
	return nil
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.Type == tea.KeyCtrlC {
			return a, tea.Quit
		}
	}
	return a, nil
}

func (a app) View() string {
	return "Hello world"
}

func main() {
	prog := tea.NewProgram(app{})
	_, err := prog.Run()
	if err != nil {
		panic(err)
	}
}
