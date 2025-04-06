package snake

import (
	"errors"
	"fmt"
	"image"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fogleman/gg"
	ascii "github.com/qeesung/image2ascii/convert"
)

type Vector struct {
	X, Y int
}

func New(width, height int) Model {
	model := Model{
		width:  width,
		height: height,
	}
	ctx := gg.NewContext(width, height)
	model.image = ctx.Image()
	model.pos = image.Point{model.image.Bounds().Max.X / 2, model.image.Bounds().Max.Y / 2}
	model.velocity = Vector{1, 0}
	model.Render()
	model.state = RUNNING
	return model
}

func (m *Model) Render() {
	ctx := gg.NewContextForImage(m.image)
	ctx.SetRGBA255(255, 0, 0, 255)
	ctx.DrawCircle(float64(m.pos.X), float64(m.pos.Y), 1)
	ctx.Fill()
	m.image = ctx.Image()
	converter := ascii.NewImageConverter()
	asciiString := converter.Image2ASCIIString(m.image, &ascii.DefaultOptions)
	m.imageBuf = asciiString
}

type Model struct {
	width, height int
	state         State
	image         image.Image
	imageBuf      string
	pos           image.Point
	velocity      Vector
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
		if direction == UP {
			m.velocity = Vector{0, -1}
		} else if direction == DOWN {
			m.velocity = Vector{0, 1}
		} else if direction == LEFT {
			m.velocity = Vector{-1, 0}
		} else if direction == RIGHT {
			m.velocity = Vector{1, 0}
		} else {
			panic(errors.New(fmt.Sprintf("Unknown direction %d", direction)))
		}
		m.Render()
		return m, nil
	}
	if _, ok := msg.(tickMsg); ok {
		if m.state == RUNNING {
			m.pos.X += m.velocity.X
			m.pos.Y += m.velocity.Y
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
	gg.SavePNG("out.png", m.image)
	return nil
}
