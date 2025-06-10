package app

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/bengesoff/mail-tui/internal/core"
	"github.com/bengesoff/mail-tui/internal/ui"
	"github.com/bengesoff/mail-tui/internal/ui/email_composer"
	"github.com/bengesoff/mail-tui/internal/ui/email_list"
	"github.com/bengesoff/mail-tui/internal/ui/email_viewer"
)

type ViewName string

const (
	ListViewName     ViewName = "email_list"
	ViewerViewName   ViewName = "email_viewer"
	ComposerViewName ViewName = "email_composer"
)

type AppModel struct {
	activeView    ViewName
	emailViewer   *email_viewer.EmailViewerModel
	emailList     *email_list.EmailListModel
	emailComposer *email_composer.EmailComposerModel
}

func NewAppModel(backend core.EmailBackend) *AppModel {
	return &AppModel{
		activeView:    ListViewName,
		emailViewer:   email_viewer.NewEmailViewerModel(backend),
		emailList:     email_list.NewEmailListModel(backend),
		emailComposer: email_composer.NewEmailComposerModel(backend),
	}
}

func (m AppModel) Init() tea.Cmd {
	return func() tea.Msg {
		return ui.ShowEmailListMessage{}
	}
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			case ComposerViewName:
				m.emailComposer, cmd = m.emailComposer.Update(msg)
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
	case ui.ShowEmailComposerMessage:
		m.activeView = ComposerViewName
		m.emailComposer, cmd = m.emailComposer.Update(msg)
		commands = append(commands, cmd)
	default:
		m.emailList, cmd = m.emailList.Update(msg)
		commands = append(commands, cmd)
		m.emailViewer, cmd = m.emailViewer.Update(msg)
		commands = append(commands, cmd)
		m.emailComposer, cmd = m.emailComposer.Update(msg)
		commands = append(commands, cmd)
	}

	return m, tea.Batch(commands...)
}

func (m AppModel) View() string {
	switch m.activeView {
	case ListViewName:
		return m.emailList.View()
	case ViewerViewName:
		return m.emailViewer.View()
	case ComposerViewName:
		return m.emailComposer.View()
	default:
		return "Unknown view " + string(m.activeView)
	}
}
