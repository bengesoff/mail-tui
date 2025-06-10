package email_composer

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/bengesoff/mail-tui/internal/core"
	"github.com/bengesoff/mail-tui/internal/ui"
)

const (
	toField = iota
	subjectField
	bodyField
	submitButton
)

var (
	focusedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#FF06B7", Dark: "#FF06B7"})
	blurredStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#767676", Dark: "#767676"})
	labelStyle = lipgloss.NewStyle().
			Bold(true)
	buttonStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.AdaptiveColor{Light: "#767676", Dark: "#767676"}).
			Padding(0, 1)

	focusedButtonStyle = buttonStyle.
				Background(lipgloss.AdaptiveColor{Light: "#FF06B7", Dark: "#FF06B7"}).
				Foreground(lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#000000"})
)

type EmailComposerModel struct {
	backend    core.EmailBackend
	focusIndex int
	toInput    textinput.Model
	subInput   textinput.Model
	bodyInput  textarea.Model
	width      int
	height     int
}

func NewEmailComposerModel(backend core.EmailBackend) *EmailComposerModel {
	toInput := textinput.New()
	toInput.Placeholder = "recipient@example.com"
	toInput.CharLimit = 256
	toInput.Width = 50

	subInput := textinput.New()
	subInput.Placeholder = "Email subject"
	subInput.CharLimit = 256
	subInput.Width = 50

	bodyInput := textarea.New()
	bodyInput.Placeholder = "Compose your email here..."
	bodyInput.SetWidth(50)
	bodyInput.SetHeight(10)

	return &EmailComposerModel{
		backend:    backend,
		focusIndex: toField,
		toInput:    toInput,
		subInput:   subInput,
		bodyInput:  bodyInput,
		width:      80,
		height:     24,
	}
}

func (m *EmailComposerModel) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		m.toInput.Focus(),
	)
}

func (m *EmailComposerModel) Update(msg tea.Msg) (*EmailComposerModel, tea.Cmd) {
	var commands []tea.Cmd

	switch msg := msg.(type) {
	case ui.ShowEmailComposerMessage:
		m.toInput.SetValue("")
		m.subInput.SetValue("")
		m.bodyInput.SetValue("")
		m.focusIndex = toField
		cmd := m.updateFieldFocus()
		return m, cmd

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateSizes()

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg {
				return ui.ShowEmailListMessage{}
			}

		case "tab", "shift+tab":
			return m, m.handleNavigation(msg.String())

		case "enter":
			if m.focusIndex == submitButton {
				return m, func() tea.Msg {
					// TODO: send the email
					return ui.ShowEmailListMessage{}
				}
			}
		}
	}

	var cmd tea.Cmd
	switch m.focusIndex {
	case toField:
		m.toInput, cmd = m.toInput.Update(msg)
		commands = append(commands, cmd)
	case subjectField:
		m.subInput, cmd = m.subInput.Update(msg)
		commands = append(commands, cmd)
	case bodyField:
		m.bodyInput, cmd = m.bodyInput.Update(msg)
		commands = append(commands, cmd)
	}

	return m, tea.Batch(commands...)
}

func (m *EmailComposerModel) handleNavigation(key string) tea.Cmd {
	switch key {
	case "shift+tab":
		m.focusIndex--
		if m.focusIndex < 0 {
			m.focusIndex = submitButton
		}
	case "tab":
		m.focusIndex++
		if m.focusIndex > submitButton {
			m.focusIndex = toField
		}
	}

	return m.updateFieldFocus()
}

func (m *EmailComposerModel) updateFieldFocus() tea.Cmd {
	var cmds []tea.Cmd

	switch m.focusIndex {
	case toField:
		cmds = append(cmds, m.toInput.Focus())
		m.subInput.Blur()
		m.bodyInput.Blur()
		m.toInput.PromptStyle = focusedStyle
		m.toInput.TextStyle = focusedStyle
		m.subInput.PromptStyle = blurredStyle
		m.subInput.TextStyle = blurredStyle

	case subjectField:
		cmds = append(cmds, m.subInput.Focus())
		m.toInput.Blur()
		m.bodyInput.Blur()
		m.toInput.PromptStyle = blurredStyle
		m.toInput.TextStyle = blurredStyle
		m.subInput.PromptStyle = focusedStyle
		m.subInput.TextStyle = focusedStyle

	case bodyField:
		cmds = append(cmds, m.bodyInput.Focus())
		m.toInput.Blur()
		m.subInput.Blur()
		m.toInput.PromptStyle = blurredStyle
		m.toInput.TextStyle = blurredStyle
		m.subInput.PromptStyle = blurredStyle
		m.subInput.TextStyle = blurredStyle

	case submitButton:
		m.toInput.Blur()
		m.subInput.Blur()
		m.bodyInput.Blur()
		m.toInput.PromptStyle = blurredStyle
		m.toInput.TextStyle = blurredStyle
		m.subInput.PromptStyle = blurredStyle
		m.subInput.TextStyle = blurredStyle
	}

	return tea.Batch(cmds...)
}

func (m *EmailComposerModel) updateSizes() {
	m.toInput.Width = m.width
	m.subInput.Width = m.width

	usedHeight := 15 // height used by title and other fields
	bodyHeight := max(m.height-usedHeight, 5)

	m.bodyInput.SetWidth(m.width)
	m.bodyInput.SetHeight(bodyHeight)
}

func (m *EmailComposerModel) View() string {
	var b strings.Builder

	b.WriteString(labelStyle.Render("Compose Email"))
	b.WriteString("\n\n")

	b.WriteString(labelStyle.Render("To:"))
	b.WriteString("\n")
	b.WriteString(m.toInput.View())
	b.WriteString("\n\n")

	b.WriteString(labelStyle.Render("Subject:"))
	b.WriteString("\n")
	b.WriteString(m.subInput.View())
	b.WriteString("\n\n")

	b.WriteString(labelStyle.Render("Body:"))
	b.WriteString("\n")
	b.WriteString(m.bodyInput.View())
	b.WriteString("\n\n")

	if m.focusIndex == submitButton {
		b.WriteString(focusedButtonStyle.Render("Send"))
	} else {
		b.WriteString(buttonStyle.Render("Send"))
	}
	b.WriteString("\n\n")

	b.WriteString(blurredStyle.Render("Tab/Shift+Tab: Navigate • Enter: Send • Esc: Cancel"))

	return b.String()
}
