package email_list

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/bengesoff/mail-tui/internal/core"
)

var (
	emailSubjectStyle = lipgloss.NewStyle().
				Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}).
			Padding(0, 0, 0, 2)

	selectedStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
			Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"}).
			Padding(0, 0, 0, 1)
)

type emailListItem struct {
	core.EmailMetadata
}

func (i *emailListItem) FilterValue() string {
	return i.Subject
}

type listItemDelegate struct{}

func (d *listItemDelegate) Height() int {
	return 2
}

func (d *listItemDelegate) Spacing() int {
	return 0
}

func (d *listItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d *listItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	email, ok := item.(*emailListItem)
	if !ok {
		return
	}

	selected := " "
	style := normalStyle.Render
	if index == m.Index() {
		selected = ">"
		style = selectedStyle.Render
	}

	sent, err := email.SentAt.MarshalText()
	if err != nil {
		return
	}

	fmt.Fprintf(w, style("%s %s\n  %s (%s)"), selected, emailSubjectStyle.Render(email.Subject), email.From, sent)
}
