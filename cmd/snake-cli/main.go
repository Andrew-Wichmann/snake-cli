package main

import (
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
	image, err := png.Decode(file)
	if err != nil {
		panic(err)
	}
	app.img = image
	app.Render()
	return app
}

func (a *app) Render() {
	ctx := gg.NewContextForImage(a.img)
	rect := a.img.Bounds()
	ctx.SetRGBA255(255, 0, 0, 255)
	ctx.DrawCircle(float64(rect.Max.X)/2, float64(rect.Max.Y)/2, 50)
	ctx.Fill()
	a.img = ctx.Image()
	converter := ascii.NewImageConverter()
	asciiString := converter.Image2ASCIIString(a.img, &ascii.DefaultOptions)
	a.imageBuf = asciiString
}

type app struct {
	img      image.Image
	imageBuf string
}

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
