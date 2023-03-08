package main

import (
	"fmt"
	"os"
	"time"

	"github.com/SamoKopecky/pqcom/main/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// cmd.Execute()

	send := make(chan []byte)
	recv := make(chan []byte)
	model := tui.NewModel(send, recv)
	go func() {
		time.Sleep(time.Second)
		for i := 0; i < 10; i++ {
			send <- []byte("Foo: hi")
			model.Update(tea.KeyEnter)
		}
	}()
	if _, err := tea.NewProgram(model, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error while running program:", err)
		os.Exit(1)
	}
}
