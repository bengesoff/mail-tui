package main

import (
	"fmt"
	"os"

	"mail-tui/internal/backend/fake"
	"mail-tui/internal/ui/email_list"
	"mail-tui/internal/ui/email_viewer"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	ListViewName   = "email_list"
	ViewerViewName = "email_viewer"
)

type model struct {
	activeView  string
	emailViewer *email_viewer.EmailViewerModel
	emailList   *email_list.EmailListModel
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.emailList.Init(),
		m.emailViewer.Init(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		commands []tea.Cmd
		cmd      tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "l":
			commands = append(commands, func() tea.Msg {
				// simulate loading an email
				return email_viewer.DisplayEmailMessage{
					EmailId: "1",
				}
			})
		}
	case email_viewer.DisplayEmailMessage:
		m.activeView = ViewerViewName
	}

	m.emailViewer, cmd = m.emailViewer.Update(msg)
	commands = append(commands, cmd)

	m.emailList, cmd = m.emailList.Update(msg)
	commands = append(commands, cmd)

	return m, tea.Batch(commands...)
}

func (m model) View() string {
	switch m.activeView {
	case ListViewName:
		return m.emailList.View()
	case ViewerViewName:
		return m.emailViewer.View()
	default:
		return "Unknown view " + m.activeView
	}
}

func initialModel() model {
	backend := &fake.FakeBackend{}
	return model{
		activeView:  ListViewName,
		emailViewer: email_viewer.NewEmailViewerModel(backend),
		emailList:   email_list.NewEmailListModel(backend),
	}
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
