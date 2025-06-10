package email_viewer

import (
	"mail-tui/internal/core"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type DisplayEmailMessage struct {
	EmailId core.EmailId
}

type EmailLoadedMessage struct {
	Email *core.Email
	Error error
}

type EmailViewerModel struct {
	email   *core.Email
	backend core.EmailBackend

	ready   bool
	loading bool
	error   string

	viewport viewport.Model
}

func NewEmailViewerModel(backend core.EmailBackend) *EmailViewerModel {
	return &EmailViewerModel{
		backend: backend,
	}
}

func (m *EmailViewerModel) Init() tea.Cmd {
	return nil
}

func (m *EmailViewerModel) Update(msg tea.Msg) (*EmailViewerModel, tea.Cmd) {
	var commands []tea.Cmd

	switch msg := msg.(type) {
	case DisplayEmailMessage:
		m.loading = true
		m.error = ""
		m.email = nil
		commands = append(commands, m.loadEmail(msg.EmailId))
	case EmailLoadedMessage:
		m.loading = false
		if msg.Error != nil {
			m.error = msg.Error.Error()
		} else {
			m.email = msg.Email
			m.error = ""
		}
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

	var viewportCommand tea.Cmd
	m.viewport, viewportCommand = m.viewport.Update(msg)
	commands = append(commands, viewportCommand)

	return m, tea.Batch(commands...)
}

func (m *EmailViewerModel) View() string {
	if !m.ready {
		return "Loading viewer..."
	}

	if m.loading {
		return "Loading email..."
	}

	if m.error != "" {
		return "Error: " + m.error
	}

	return m.viewport.View()
}

func (m *EmailViewerModel) updateViewportContent() error {
	if !m.ready {
		return nil
	}

	if m.email != nil {
		content, err := RenderEmail(m.email, m.viewport.Width)
		if err != nil {
			return err
		}
		m.viewport.SetContent(content)
	} else {
		m.viewport.SetContent("No email selected")
	}
	return nil
}

func (m *EmailViewerModel) loadEmail(emailId core.EmailId) tea.Cmd {
	return func() tea.Msg {
		email, err := m.backend.GetEmail(emailId)
		if err != nil {
			return EmailLoadedMessage{
				Email: nil,
				Error: err,
			}
		}

		return EmailLoadedMessage{
			Email: email,
			Error: nil,
		}
	}
}
