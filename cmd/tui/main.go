package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	emails []email
	cursor int
}

type email struct {
	from    string
	subject string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.emails)-1 {
				m.cursor++
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	output := "Emails:\n\n"
	for i, email := range m.emails {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		output += fmt.Sprintf("%s [%s] – %s\n", cursor, email.from, email.subject)
	}
	output += "\n Press q to quit.\n"
	return output
}

func initialModel() model {
	return model{
		emails: []email{
			{from: "test@test.com", subject: "test1"},
			{from: "test@test.com", subject: "test2"},
			{from: "test@test.com", subject: "test3"},
			{from: "test@test.com", subject: "test4"},
		},
		cursor: 0,
	}
}

func main() {
	program := tea.NewProgram(initialModel())
	_, err := program.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
