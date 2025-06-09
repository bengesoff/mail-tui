package main

import (
	"fmt"
	"os"
	"time"

	"mail-tui/internal/core"
	"mail-tui/internal/ui/email_viewer"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	emailViewer *email_viewer.EmailViewerModel
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		commands []tea.Cmd
		cmd      tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "l":
			commands = append(commands, func() tea.Msg {
				// simulate loading an email
				return email_viewer.DisplayEmailMessage{
					Email: &core.Email{
						From:    "from@example.com",
						To:      "to@example.com",
						Subject: "Test email",
						SentAt:  time.Now(),
						Body: "To whom it may concern,\n\n" +
							"This is a test email.\n\n" +
							"Yours sincerely,\n\n" +
							"Tester",
					},
				}
			})
		}
	}

	m.emailViewer, cmd = m.emailViewer.Update(msg)
	commands = append(commands, cmd)

	return m, tea.Batch(commands...)
}

func (m model) View() string {
	return m.emailViewer.View()
}

func initialModel() model {
	return model{&email_viewer.EmailViewerModel{}}
}

func main() {
	program := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion())
	_, err := program.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
