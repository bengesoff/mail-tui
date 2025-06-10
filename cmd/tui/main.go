package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/bengesoff/mail-tui/internal/backend/fake"
	"github.com/bengesoff/mail-tui/internal/ui"
	"github.com/bengesoff/mail-tui/internal/ui/email_list"
	"github.com/bengesoff/mail-tui/internal/ui/email_viewer"
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
	return func() tea.Msg {
		return ui.ShowEmailListMessage{}
	}
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
		default:
			// only send key messages to the active view
			switch m.activeView {
			case ListViewName:
				m.emailList, cmd = m.emailList.Update(msg)
				commands = append(commands, cmd)
			case ViewerViewName:
				m.emailViewer, cmd = m.emailViewer.Update(msg)
				commands = append(commands, cmd)
			}
		}
	case ui.ShowEmailListMessage:
		m.activeView = ListViewName
		m.emailList, cmd = m.emailList.Update(msg)
		commands = append(commands, cmd)
	case ui.ShowEmailViewerMessage:
		m.activeView = ViewerViewName
		m.emailViewer, cmd = m.emailViewer.Update(msg)
		commands = append(commands, cmd)
	default:
		m.emailList, cmd = m.emailList.Update(msg)
		commands = append(commands, cmd)
		m.emailViewer, cmd = m.emailViewer.Update(msg)
		commands = append(commands, cmd)
	}

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
