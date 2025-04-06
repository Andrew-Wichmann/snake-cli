package snake

import (
	"fmt"
	"image"
	"image/color"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fogleman/gg"
	ascii "github.com/qeesung/image2ascii/convert"
)

type Vector struct {
	X, Y int
}

type Model struct {
	width, height  int
	state          State
	image          *image.Image
	imageBuf       string
	asciiConverter *ascii.ImageConverter
}

func New(width, height int) Model {
	converter := ascii.NewImageConverter()
	model := Model{
		width:          width,
		height:         height,
		asciiConverter: converter,
		state:          RUNNING,
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
	//if direction, ok := msg.(Direction); ok {
	//	if direction == UP {
	//		m.velocity = Vector{0, -1}
	//	} else if direction == DOWN {
	//		m.velocity = Vector{0, 1}
	//	} else if direction == LEFT {
	//		m.velocity = Vector{-1, 0}
	//	} else if direction == RIGHT {
	//		m.velocity = Vector{1, 0}
	//	} else {
	//		panic(errors.New(fmt.Sprintf("Unknown direction %d", direction)))
	//	}
	//	m.Render()
	//	return m, nil
	//}
	if _, ok := msg.(tickMsg); ok {
		if m.state == RUNNING {
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
	ctx.DrawCircle(float64(m.width/2), float64(m.height/2), 10)
	ctx.SetColor(color.RGBA{0, 255, 0, 255})
	ctx.Fill()
	return ctx.Image()
}

func (m *Model) Render() {
	i := m.drawImage()
	asciiString := m.asciiConverter.Image2ASCIIString(i, &ascii.DefaultOptions)
	m.imageBuf = asciiString
}
