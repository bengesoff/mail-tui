package email_list

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/bengesoff/mail-tui/internal/core"
)

var (
	emailSubjectStyle = lipgloss.NewStyle().
				Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#555555", Dark: "#bbbbbb"}).
			Padding(0, 0, 0, 2)

	selectedStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
			Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"}).
			Padding(0, 0, 0, 1)

	unreadStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"}).
			Padding(0, 0, 0, 2)

	selectedUnreadStyle = lipgloss.NewStyle().
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

	unread := !email.IsRead
	selected := index == m.Index()

	unreadIndicator := " "
	if unread {
		unreadIndicator = "●"
	}

	selectedIndicator := " "
	if selected {
		selectedIndicator = ">"
	}

	style := normalStyle.Render
	if selected && unread {
		style = selectedUnreadStyle.Render
	} else if selected && !unread {
		style = selectedStyle.Render
	} else if !selected && unread {
		style = unreadStyle.Render
	}

	sent, err := email.SentAt.MarshalText()
	if err != nil {
		return
	}

	_, _ = fmt.Fprintf(w, style("%s%s %s\n   %s (%s)"),
		selectedIndicator,
		unreadIndicator,
		emailSubjectStyle.Render(email.Subject),
		email.From,
		sent,
	)
}

func (d *listItemDelegate) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "view email"),
		),
		key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "compose email"),
		),
	}
}

func (d *listItemDelegate) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "view email"),
			),
			key.NewBinding(
				key.WithKeys("c"),
				key.WithHelp("c", "compose email"),
			),
		},
	}
}
