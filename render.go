package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type RenderModel struct {
	progress progress.Model

	counter *BandwidthCounter
}

var setColor = func(color string) func(strs ...string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render
}
var helpStyle = setColor("#626262")
var primaryStyle = setColor("#CCE0AC")
var infoStyle = setColor("#F4DEB3")
var dangerStyle = setColor("#FF8A8A")

func NewRenderModel(counter *BandwidthCounter) RenderModel {
	return RenderModel{
		progress: progress.New(progress.WithDefaultGradient()),
		counter:  counter,
	}
}

func (m RenderModel) Init() tea.Cmd {
	return nil
}

func (m RenderModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width / 2
		return m, tea.Quit
	default:
		return m, tea.Quit
	}
}

func (m RenderModel) View() string {
	pad := " "

	// 列出信息
	total := "Total: " + primaryStyle(humanSpace(m.counter.MonthlyBWLimitB))
	used := "Used: " + infoStyle(humanSpace(m.counter.BWCounterB))
	var limit string
	localtion, _ := time.LoadLocation("America/Los_Angeles")
	now := time.Now().In(localtion)
	if now.Day() > m.counter.BWResetDayOfMonth {
		limit = helpStyle("Network traffic will be reset on the ") + dangerStyle(fmt.Sprint(m.counter.BWResetDayOfMonth)+"th") + helpStyle(" of next month?!🫠")
	} else if now.Day() == m.counter.BWResetDayOfMonth {
		limit = helpStyle("Network traffic has been reset today~~🥰")
	} else {
		limit = helpStyle("Network traffic will be reset on the ") + dangerStyle(fmt.Sprint(m.counter.BWResetDayOfMonth, "th")) + helpStyle(" of this month!!🤩")
	}
	limit = helpStyle(limit)

	var percentage float64
	if m.counter.MonthlyBWLimitB != 0 {
		percentage = float64(m.counter.BWCounterB) / float64(m.counter.MonthlyBWLimitB)
	}
	return "\n" + pad + m.progress.ViewAs(percentage) + "\n\n" +
		pad + total + "\n" +
		pad + used + "\n\n" +
		pad + limit + "\n"
}

var spaceMap = map[int]string{
	0: "B",
	1: "KB",
	2: "MB",
	3: "GB",
	4: "TB",
}

func humanSpace(b int64) string {
	newb := float64(b)
	var i int
	for newb > 1000 && i < len(spaceMap)-1 {
		newb /= 1000
		i++
	}
	return fmt.Sprintf("%.2f%s", newb, spaceMap[i])
}
