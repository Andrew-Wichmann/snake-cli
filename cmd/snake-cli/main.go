package main

import (
	"errors"
	"fmt"
	"image"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fogleman/gg"
	ascii "github.com/qeesung/image2ascii/convert"
)

func NewGame() app {
	var a app
	ctx := gg.NewContext(width, height*2) // Not sure why I have to double the height...
	a.image = ctx.Image()
	a.pos = image.Point{a.image.Bounds().Max.X / 2, a.image.Bounds().Max.Y / 2}
	a.Render()
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
}

func (a app) Init() tea.Cmd {
	return nil
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
		case "n":
			app := NewGame()
			return app, nil
		}
	}
	if direction, ok := msg.(Direction); ok {
		if direction == UP {
			a.pos.Y += -1
		} else if direction == DOWN {
			a.pos.Y += 1
		} else if direction == LEFT {
			a.pos.X += -1
		} else if direction == RIGHT {
			a.pos.X += 1
		} else {
			panic(errors.New(fmt.Sprintf("Unknown direction %d", direction)))
		}
		a.Render()
	}
	return a, nil
}

func (a app) View() string {
	return a.imageBuf
}

func main() {
	prog := tea.NewProgram(app{})
	_, err := prog.Run()
	if err != nil {
		panic(err)
	}
}
