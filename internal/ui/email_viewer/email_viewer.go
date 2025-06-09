package email_viewer

import (
	"mail-tui/internal/core"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type DisplayEmailMessage struct {
	Email *core.Email
}

type EmailViewerModel struct {
	email *core.Email

	ready bool
	error string

	viewport viewport.Model
}

func (m *EmailViewerModel) Init() tea.Cmd {
	return nil
}

func (m *EmailViewerModel) Update(msg tea.Msg) (*EmailViewerModel, tea.Cmd) {
	switch msg := msg.(type) {
	case DisplayEmailMessage:
		m.email = msg.Email
		err := m.updateViewportContent()
		if err != nil {
			m.error = err.Error()
			return m, nil
		}
	case tea.WindowSizeMsg:
		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height)
			m.viewport.SetContent("No email selected")
			m.ready = true
		} else {
			m.viewport.Height = msg.Height
			m.viewport.Width = msg.Width
		}

		err := m.updateViewportContent()
		if err != nil {
			m.error = err.Error()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m *EmailViewerModel) updateViewportContent() error {
	if m.ready && m.email != nil {
		content, err := RenderEmail(m.email, m.viewport.Width)
		if err != nil {
			return err
		}
		m.viewport.SetContent(content)
	}
	return nil
}

func (m *EmailViewerModel) View() string {
	if !m.ready {
		return "Loading viewer..."
	}
	return m.viewport.View()
}
