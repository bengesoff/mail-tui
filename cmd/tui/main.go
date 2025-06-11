package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/bengesoff/mail-tui/internal/backend/fake"
	"github.com/bengesoff/mail-tui/internal/ui/app"
)

func main() {
	backend := fake.NewFakeBackend()
	appModel := app.NewAppModel(backend)

	program := tea.NewProgram(
		appModel,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion())

	_, err := program.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
