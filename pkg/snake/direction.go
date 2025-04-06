package snake

import (
	tea "github.com/charmbracelet/bubbletea"
)

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
