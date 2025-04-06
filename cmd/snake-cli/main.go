package main

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fogleman/gg"
	ascii "github.com/qeesung/image2ascii/convert"
)

func newApp() app {
	var app app
	file, err := os.Open("silly-grid.png")
	if err != nil {
		panic(err)
	}
	grid, err := png.Decode(file)
	if err != nil {
		panic(err)
	}
	app.image = grid
	app.pos = image.Point{app.image.Bounds().Max.X / 2, app.image.Bounds().Max.Y / 2}
	app.Render()
	return app
}

func (a *app) Render() {
	ctx := gg.NewContextForImage(a.image)
	ctx.SetRGBA255(255, 0, 0, 255)
	ctx.DrawCircle(float64(a.pos.X), float64(a.pos.Y), 10)
	ctx.Fill()
	a.image = ctx.Image()
	converter := ascii.NewImageConverter()
	asciiString := converter.Image2ASCIIString(a.image, &ascii.DefaultOptions)
	a.imageBuf = asciiString
}

type app struct {
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
	app := newApp()
	prog := tea.NewProgram(app)
	_, err := prog.Run()
	if err != nil {
		panic(err)
	}
}
