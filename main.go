package main

import (
	"context"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	conf := JustConf{
		UrlWidthCounter: os.Getenv("JUST_MY_SOCKS_WIDTH_COUNTER"),
	}
	just := NewJust(conf)
	counter, err := just.GetWidthCounter(context.Background())
	if err != nil {
		panic(err)
	}
	render := NewRenderModel(counter)
	if _, err := tea.NewProgram(render).Run(); err != nil {
		panic(err)
	}
}
