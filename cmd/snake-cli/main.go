package main

import (
	"errors"
	"fmt"
	"image"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fogleman/gg"
	ascii "github.com/qeesung/image2ascii/convert"
)

type Vector struct {
	X, Y int
}

func NewGame() app {
	var a app
	ctx := gg.NewContext(width, height*2) // Not sure why I have to double the height...
	a.image = ctx.Image()
	a.pos = image.Point{a.image.Bounds().Max.X / 2, a.image.Bounds().Max.Y / 2}
	a.velocity = Vector{1, 0}
	a.Render()
	a.state = RUNNING
	return a
}

func (a *app) Render() {
	ctx := gg.NewContextForImage(a.image)
	ctx.SetRGBA255(255, 0, 0, 255)
	ctx.DrawCircle(float64(a.pos.X), float64(a.pos.Y), 1)
	ctx.Fill()
	a.image = ctx.Image()
	converter := ascii.NewImageConverter()
	asciiString := converter.Image2ASCIIString(a.image, &ascii.DefaultOptions)
	a.imageBuf = asciiString
}

type State int

var height int
var width int

const (
	INIT    State = 0
	RUNNING State = 1
	PAUSED  State = 2 // maybe useful later (shrug)
)

type app struct {
	state    State
	image    image.Image
	imageBuf string
	pos      image.Point
	velocity Vector
}

type tickMsg time.Time

func doTick() tea.Cmd {
	return tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (a app) Init() tea.Cmd {
	return doTick()
}

type Direction int

const (
	UP    Direction = 0
	DOWN  Direction = 1
	LEFT  Direction = 2
	RIGHT Direction = 3
)

func UpCmd() tea.Msg {
	return UP
}

func DownCmd() tea.Msg {
	return DOWN
}

func LeftCmd() tea.Msg {
	return LEFT
}

func RightCmd() tea.Msg {
	return RIGHT
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if windowSizeMsg, ok := msg.(tea.WindowSizeMsg); ok {
		height = windowSizeMsg.Height
		width = windowSizeMsg.Width
		app := NewGame()
		return app, nil
	}
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.Type == tea.KeyCtrlC {
			return a, tea.Quit
		}
		switch keyMsg.String() {
		case "k":
			return a, UpCmd
		case "j":
			return a, DownCmd
		case "l":
			return a, RightCmd
		case "h":
			return a, LeftCmd
		case "p":
			if a.state == PAUSED {
				a.state = RUNNING
			} else {
				a.state = PAUSED
			}
			return a, nil
		case "n":
			app := NewGame()
			return app, nil
		}
	}
	if direction, ok := msg.(Direction); ok {
		if direction == UP {
			a.velocity = Vector{0, -1}
		} else if direction == DOWN {
			a.velocity = Vector{0, 1}
		} else if direction == LEFT {
			a.velocity = Vector{-1, 0}
		} else if direction == RIGHT {
			a.velocity = Vector{1, 0}
		} else {
			panic(errors.New(fmt.Sprintf("Unknown direction %d", direction)))
		}
		a.Render()
		return a, nil
	}
	if _, ok := msg.(tickMsg); ok {
		if a.state == RUNNING {
			a.pos.X += a.velocity.X
			a.pos.Y += a.velocity.Y
			a.Render()
		}
		return a, doTick()
	}
	return a, nil
}

func (a app) View() string {
	if a.state == INIT {
		return "Starting"
	} else if a.state == RUNNING {
		return a.imageBuf
	} else if a.state == PAUSED {
		return "Paused"
	} else {
		panic(fmt.Sprintf("Unknown state %d", a.state))
	}
}

func main() {
	prog := tea.NewProgram(app{})
	_, err := prog.Run()
	if err != nil {
		panic(err)
	}
}
