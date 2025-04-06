package main

import (
	"fmt"

	"github.com/Andrew-Wichmann/snake-cli/pkg/snake"
	tea "github.com/charmbracelet/bubbletea"
)

type state int

const (
	INIT    state = 0
	RUNNING state = 1
)

type app struct {
	width  int
	height int
	snake  snake.Model
	state  state
}

func (a app) Init() tea.Cmd {
	return a.snake.Init()
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if windowSizeMsg, ok := msg.(tea.WindowSizeMsg); ok {
		a.width = windowSizeMsg.Width
		a.height = windowSizeMsg.Height * 2 // Not sure why I have to double the height...
		prog := snake.New(a.width, a.height)
		a.snake = prog
		a.state = RUNNING
		return a, nil
	}
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.String() == "n" {
			prog := snake.New(a.width, a.height)
			a.snake = prog
			a.state = RUNNING
			return a, nil
		}
		if keyMsg.Type == tea.KeyCtrlC {
			return a, tea.Quit
		}
	}
	snakeModel, cmd := a.snake.Update(msg)
	a.snake = snakeModel
	return a, cmd
}

func (a app) View() string {
	if a.state == INIT {
		return "Starting"
	} else if a.state == RUNNING {
		return a.snake.View()
	} else {
		panic(fmt.Sprintf("Unknown state %d", a.state))
	}
}

func main() {
	f, err := tea.LogToFile("debug.log", "snake-cli")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	prog := tea.NewProgram(app{})
	_, err = prog.Run()
	if err != nil {
		panic(err)
	}
}
