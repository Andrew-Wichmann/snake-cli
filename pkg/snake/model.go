package snake

import (
	"fmt"
	"image"
	"image/color"
	"math"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fogleman/gg"
	ascii "github.com/qeesung/image2ascii/convert"
)

// Not sure where to put this.
type Vector struct {
	X, Y int
}

type Model struct {
	width, height  int
	state          State
	image          *image.Image
	imageBuf       string
	asciiConverter *ascii.ImageConverter
	body           body
}

func New(width, height int) Model {
	converter := ascii.NewImageConverter()
	body := newBody()
	model := Model{
		width:          width,
		height:         height,
		asciiConverter: converter,
		state:          RUNNING,
		body:           body,
	}
	model.Render()
	return model
}

func (m Model) Init() tea.Cmd {
	return doTick()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "k":
			return m, UpCmd
		case "j":
			return m, DownCmd
		case "l":
			return m, RightCmd
		case "h":
			return m, LeftCmd
		case "ctrl+p":
			return m, m.Save
		case "p":
			if m.state == PAUSED {
				m.state = RUNNING
			} else {
				m.state = PAUSED
			}
			return m, nil
		}
	}
	if direction, ok := msg.(Direction); ok {
		m.body.changeDirection(direction)
		m.Render()
		return m, nil
	}
	if _, ok := msg.(tickMsg); ok {
		if m.state == RUNNING {
			m.body.grow()
			//m.pos.X += m.velocity.X
			//m.pos.Y += m.velocity.Y
			m.Render()
		}
		return m, doTick()
	}
	return m, nil
}

func (m Model) View() string {
	if m.state == RUNNING {
		return m.imageBuf
	} else if m.state == PAUSED {
		return "Paused"
	} else {
		panic(fmt.Sprintf("Unknown state %d", m.state))
	}
}

func (m Model) Save() tea.Msg {
	i := m.drawImage()
	gg.SavePNG("out.png", i)
	return nil
}

func (m Model) drawImage() image.Image {
	i := image.NewRGBA(image.Rect(0, 0, m.width, m.height))
	ctx := gg.NewContextForImage(i)
	ctx.SetColor(color.RGBA{255, 0, 0, 255})
	for _, segment := range m.body.segments {
		minx := min(segment.start.x, segment.end.x)
		miny := min(segment.start.y, segment.end.y)
		width := math.Abs(segment.start.x-segment.end.x) + 1
		height := math.Abs(segment.start.y-segment.end.y) + 1
		ctx.DrawRectangle(minx, miny, width, height)
		ctx.Fill()
	}
	return ctx.Image()
}

func (m *Model) Render() {
	i := m.drawImage()
	asciiString := m.asciiConverter.Image2ASCIIString(i, &ascii.DefaultOptions)
	m.imageBuf = asciiString
}
