package email_list

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/bengesoff/mail-tui/internal/core"
	"github.com/bengesoff/mail-tui/internal/ui"
)

type emailsLoadedMessage struct {
	emails []core.EmailMetadata
	error  error
}

type EmailListModel struct {
	emails  []core.EmailMetadata
	backend core.EmailBackend

	loading bool
	error   string

	list list.Model
}

func NewEmailListModel(backend core.EmailBackend) *EmailListModel {
	return &EmailListModel{
		emails:  []core.EmailMetadata{},
		backend: backend,
		list:    list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
	}
}

func (m *EmailListModel) Init() tea.Cmd {
	return nil
}

func (m *EmailListModel) Update(msg tea.Msg) (*EmailListModel, tea.Cmd) {
	var commands []tea.Cmd

	switch msg := msg.(type) {
	case ui.ShowEmailListMessage:
		m.loading = true
		m.error = ""
		commands = append(commands, m.loadEmails())
	case emailsLoadedMessage:
		m.loading = false
		if msg.error != nil {
			m.error = msg.error.Error()
		} else {
			m.emails = msg.emails
			m.error = ""
			var items []list.Item
			for _, email := range m.emails {
				items = append(items, &emailListItem{email})
			}
			m.list = newList(items)
			commands = append(commands, tea.WindowSize())
		}
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case tea.KeyEnter.String():
			commands = append(commands, func() tea.Msg {
				i := m.list.GlobalIndex()
				selectedEmail := m.emails[i]
				return ui.ShowEmailViewerMessage{
					EmailId: selectedEmail.Id,
				}
			})
		case "c":
			commands = append(commands, func() tea.Msg {
				return ui.ShowEmailComposerMessage{}
			})
		}
	}

	var listCommand tea.Cmd
	m.list, listCommand = m.list.Update(msg)
	commands = append(commands, listCommand)

	return m, tea.Batch(commands...)
}

func (m *EmailListModel) View() string {
	if m.loading {
		return "Loading emails..."
	}

	if m.error != "" {
		return "Error loading emails: " + m.error
	}

	return m.list.View()
}

func (m *EmailListModel) loadEmails() tea.Cmd {
	return func() tea.Msg {
		emails, err := m.backend.ListEmails()
		if err != nil {
			return emailsLoadedMessage{
				emails: nil,
				error:  err,
			}
		}
		return emailsLoadedMessage{
			emails: emails,
			error:  nil,
		}
	}
}

func newList(items []list.Item) list.Model {
	list := list.New(items, &listItemDelegate{}, 0, 0)
	list.Title = "Inbox"
	list.SetStatusBarItemName("email", "emails")
	list.SetFilteringEnabled(false)
	return list
}
